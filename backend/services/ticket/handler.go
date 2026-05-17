package ticket

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/domain"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
)

// ticketQuerier is the minimal sqlc interface used by ticket handlers.
type ticketQuerier interface {
	GetTicket(ctx context.Context, id int64) (db.Ticket, error)
	GetTicketDetail(ctx context.Context, id int64) (db.GetTicketDetailRow, error)
	CreateTicket(ctx context.Context, arg db.CreateTicketParams) (sql.Result, error)
	UpdateTicket(ctx context.Context, arg db.UpdateTicketParams) error
	UpdateTicketStatus(ctx context.Context, arg db.UpdateTicketStatusParams) error
	UpdateTicketAssignee(ctx context.Context, arg db.UpdateTicketAssigneeParams) error
	CreateTicketHistory(ctx context.Context, arg db.CreateTicketHistoryParams) error
	IsUserInSystemAssignees(ctx context.Context, arg db.IsUserInSystemAssigneesParams) (int64, error)
}

// rawQuerier is used for the dynamic list query.
type rawQuerier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func validTicketType(s string) bool {
	switch db.TicketsType(s) {
	case db.TicketsTypeQuestion, db.TicketsTypeBug, db.TicketsTypeConfig, db.TicketsTypeData:
		return true
	}
	return false
}

func validTicketPriority(s string) bool {
	switch db.TicketsPriority(s) {
	case db.TicketsPriorityHigh, db.TicketsPriorityMedium, db.TicketsPriorityLow:
		return true
	}
	return false
}

func validTicketStatus(s string) bool {
	switch db.TicketsStatus(s) {
	case db.TicketsStatusNew, db.TicketsStatusInProgress, db.TicketsStatusWaiting, db.TicketsStatusDone:
		return true
	}
	return false
}

// ticketDTO is the API response shape for a single ticket.
type ticketDTO struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description"`
	Type          string    `json:"type"`
	Priority      string    `json:"priority"`
	Status        string    `json:"status"`
	CustomerID    int64     `json:"customer_id"`
	CustomerName  string    `json:"customer_name"`
	SystemID      int64     `json:"system_id"`
	SystemName    string    `json:"system_name"`
	AssigneeID    *int64    `json:"assignee_id"`
	AssigneeName  *string   `json:"assignee_name"`
	CreatedBy     int64     `json:"created_by"`
	CreatedByName string    `json:"created_by_name"`
	ReceivedAt    time.Time `json:"received_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type paginationDTO struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

type listResponseDTO struct {
	Data       []ticketDTO   `json:"data"`
	Pagination paginationDTO `json:"pagination"`
}

func toDTO(row db.GetTicketDetailRow) ticketDTO {
	dto := ticketDTO{
		ID:            row.ID,
		Title:         row.Title,
		Type:          string(row.Type),
		Priority:      string(row.Priority),
		Status:        string(row.Status),
		CustomerID:    row.CustomerID,
		CustomerName:  row.CustomerName,
		SystemID:      row.SystemID,
		SystemName:    row.SystemName,
		CreatedBy:     row.CreatedBy,
		CreatedByName: row.CreatedByName,
		ReceivedAt:    row.ReceivedAt,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
	if row.Description.Valid {
		dto.Description = &row.Description.String
	}
	if row.AssigneeID.Valid {
		dto.AssigneeID = &row.AssigneeID.Int64
	}
	if row.AssigneeName.Valid {
		dto.AssigneeName = &row.AssigneeName.String
	}
	return dto
}

// getUserID parses X-User-ID header injected by dev-auth middleware.
func getUserID(r *http.Request) (int64, bool) {
	s := r.Header.Get("X-User-ID")
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func parseTicketID(r *http.Request) (int64, bool) {
	s := chi.URLParam(r, "ticketId")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("ticket: json encode error", slog.Any("error", err))
	}
}

// ─── list filters ────────────────────────────────────────────────────────────

type listFilters struct {
	Status     string
	Priority   string
	Type       string
	CustomerID int64
	SystemID   int64
	Keyword    string
	Page       int
	PerPage    int
}

func parseListFilters(r *http.Request) listFilters {
	q := r.URL.Query()
	f := listFilters{
		Status:   q.Get("status"),
		Priority: q.Get("priority"),
		Type:     q.Get("type"),
		Keyword:  q.Get("keyword"),
		Page:     1,
		PerPage:  50,
	}
	if v, err := strconv.ParseInt(q.Get("customer_id"), 10, 64); err == nil && v > 0 {
		f.CustomerID = v
	}
	if v, err := strconv.ParseInt(q.Get("system_id"), 10, 64); err == nil && v > 0 {
		f.SystemID = v
	}
	if v, err := strconv.Atoi(q.Get("page")); err == nil && v > 0 {
		f.Page = v
	}
	if v, err := strconv.Atoi(q.Get("per_page")); err == nil && v > 0 && v <= 200 {
		f.PerPage = v
	}
	return f
}

const listSelectCols = `SELECT t.id, t.title, t.description, t.type, t.priority, t.status,
       t.customer_id, c.name AS customer_name,
       t.system_id,   s.name AS system_name,
       t.assignee_id, a.name AS assignee_name,
       t.created_by,  cb.name AS created_by_name,
       t.received_at, t.created_at, t.updated_at `

const listFromClauses = `FROM tickets t
JOIN customers c  ON c.id = t.customer_id
JOIN systems   s  ON s.id = t.system_id
LEFT JOIN users a ON a.id = t.assignee_id
JOIN users     cb ON cb.id = t.created_by`

func buildListQuery(f listFilters) (dataQuery, countQuery string, args, countArgs []any) {
	var conds []string
	var sharedArgs []any

	if f.Status != "" {
		conds = append(conds, "t.status = ?")
		sharedArgs = append(sharedArgs, f.Status)
	}
	if f.Priority != "" {
		conds = append(conds, "t.priority = ?")
		sharedArgs = append(sharedArgs, f.Priority)
	}
	if f.Type != "" {
		conds = append(conds, "t.type = ?")
		sharedArgs = append(sharedArgs, f.Type)
	}
	if f.CustomerID > 0 {
		conds = append(conds, "t.customer_id = ?")
		sharedArgs = append(sharedArgs, f.CustomerID)
	}
	if f.SystemID > 0 {
		conds = append(conds, "t.system_id = ?")
		sharedArgs = append(sharedArgs, f.SystemID)
	}
	if f.Keyword != "" {
		conds = append(conds, "(t.title LIKE ? OR t.description LIKE ?)")
		kw := "%" + f.Keyword + "%"
		sharedArgs = append(sharedArgs, kw, kw)
	}

	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}

	offset := (f.Page - 1) * f.PerPage
	dataQuery = fmt.Sprintf("%s%s%s ORDER BY t.received_at DESC LIMIT ? OFFSET ?",
		listSelectCols, listFromClauses, where)
	args = append(sharedArgs, f.PerPage, offset)

	countQuery = "SELECT COUNT(*) " + listFromClauses + where
	countArgs = sharedArgs
	return
}

func scanListRow(rows *sql.Rows) (ticketDTO, error) {
	var (
		row     db.GetTicketDetailRow
		descVal sql.NullString
		ainID   sql.NullInt64
		ainName sql.NullString
	)
	if err := rows.Scan(
		&row.ID, &row.Title, &descVal, &row.Type, &row.Priority, &row.Status,
		&row.CustomerID, &row.CustomerName,
		&row.SystemID, &row.SystemName,
		&ainID, &ainName,
		&row.CreatedBy, &row.CreatedByName,
		&row.ReceivedAt, &row.CreatedAt, &row.UpdatedAt,
	); err != nil {
		return ticketDTO{}, err
	}
	row.Description = descVal
	row.AssigneeID = ainID
	row.AssigneeName = ainName
	return toDTO(row), nil
}

// ─── handlers ────────────────────────────────────────────────────────────────

// HandleListTickets handles GET /tickets with dynamic filters and pagination.
func HandleListTickets(rawDB rawQuerier, queries ticketQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := parseListFilters(r)
		dataQ, countQ, args, countArgs := buildListQuery(f)

		var total int64
		if err := rawDB.QueryRowContext(r.Context(), countQ, countArgs...).Scan(&total); err != nil {
			slog.ErrorContext(
				r.Context(), "ticket list: count query",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		rows, err := rawDB.QueryContext(r.Context(), dataQ, args...)
		if err != nil {
			slog.ErrorContext(
				r.Context(), "ticket list: data query",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		defer func() {
			if err := rows.Close(); err != nil {
				slog.ErrorContext(r.Context(), "ticket list: close rows", slog.Any("error", err))
			}
		}()

		items := make([]ticketDTO, 0)
		for rows.Next() {
			dto, err := scanListRow(rows)
			if err != nil {
				slog.ErrorContext(
					r.Context(), "ticket list: scan",
					slog.Any("error", err),
					slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
				)
				httperr.InternalError("サーバーエラーが発生しました").Write(w)
				return
			}
			items = append(items, dto)
		}
		if err := rows.Err(); err != nil {
			slog.ErrorContext(
				r.Context(), "ticket list: rows err",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		totalPages := int((total + int64(f.PerPage) - 1) / int64(f.PerPage))
		if totalPages == 0 {
			totalPages = 1
		}
		writeJSON(w, http.StatusOK, listResponseDTO{
			Data: items,
			Pagination: paginationDTO{
				Total:      total,
				Page:       f.Page,
				PerPage:    f.PerPage,
				TotalPages: totalPages,
			},
		})
	}
}

// HandleGetTicket handles GET /tickets/{ticketId}.
func HandleGetTicket(queries ticketQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseTicketID(r)
		if !ok {
			httperr.BadRequest("INVALID_TICKET_ID", "チケットIDが不正です").Write(w)
			return
		}

		row, err := queries.GetTicketDetail(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.NotFound("TICKET_NOT_FOUND", "指定されたチケットは存在しません").Write(w)
				return
			}
			slog.ErrorContext(
				r.Context(), "ticket get: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		writeJSON(w, http.StatusOK, toDTO(row))
	}
}

// HandleCreateTicket handles POST /tickets.
func HandleCreateTicket(queries ticketQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := getUserID(r)
		if !ok {
			httperr.Unauthorized("認証が必要です").Write(w)
			return
		}

		var req struct {
			CustomerID  int64   `json:"customer_id"`
			SystemID    int64   `json:"system_id"`
			Type        string  `json:"type"`
			Priority    string  `json:"priority"`
			Title       string  `json:"title"`
			Description *string `json:"description"`
			AssigneeID  *int64  `json:"assignee_id"`
			ReceivedAt  *string `json:"received_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperr.BadRequest("INVALID_REQUEST", "リクエストの形式が不正です").Write(w)
			return
		}

		details := map[string]string{}
		if req.CustomerID == 0 {
			details["customer_id"] = "必須項目です"
		}
		if req.SystemID == 0 {
			details["system_id"] = "必須項目です"
		}
		if req.Title == "" {
			details["title"] = "必須項目です"
		}
		if len(req.Title) > 255 {
			details["title"] = "255文字以内で入力してください"
		}
		if !validTicketType(req.Type) {
			details["type"] = "無効な種別です"
		}
		if !validTicketPriority(req.Priority) {
			details["priority"] = "無効な優先度です"
		}
		if len(details) > 0 {
			httperr.UnprocessableEntity("VALIDATION_FAILED", "入力内容に誤りがあります").
				WithDetails(details).Write(w)
			return
		}

		// Validate assignee is in the system's assignee pool when specified.
		if req.AssigneeID != nil {
			count, err := queries.IsUserInSystemAssignees(r.Context(), db.IsUserInSystemAssigneesParams{
				SystemID: req.SystemID,
				UserID:   *req.AssigneeID,
			})
			if err != nil {
				slog.ErrorContext(
					r.Context(), "ticket create: check assignee pool",
					slog.Any("error", err),
					slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
				)
				httperr.InternalError("サーバーエラーが発生しました").Write(w)
				return
			}
			if count == 0 {
				httperr.UnprocessableEntity("ASSIGNEE_NOT_IN_POOL",
					"指定された担当者はこのシステムの担当者プールに含まれていません").Write(w)
				return
			}
		}

		params := db.CreateTicketParams{
			Title:      req.Title,
			Type:       db.TicketsType(req.Type),
			Priority:   db.TicketsPriority(req.Priority),
			CustomerID: req.CustomerID,
			SystemID:   req.SystemID,
			CreatedBy:  userID,
			ReceivedAt: time.Now(),
		}
		if req.Description != nil {
			params.Description = sql.NullString{String: *req.Description, Valid: true}
		}
		if req.AssigneeID != nil {
			params.AssigneeID = sql.NullInt64{Int64: *req.AssigneeID, Valid: true}
		}
		if req.ReceivedAt != nil {
			if t, err := time.Parse(time.RFC3339, *req.ReceivedAt); err == nil {
				params.ReceivedAt = t
			}
		}

		result, err := queries.CreateTicket(r.Context(), params)
		if err != nil {
			slog.ErrorContext(
				r.Context(), "ticket create: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		newID, err := result.LastInsertId()
		if err != nil {
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		row, err := queries.GetTicketDetail(r.Context(), newID)
		if err != nil {
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		writeJSON(w, http.StatusCreated, toDTO(row))
	}
}

// HandleUpdateTicket handles PUT /tickets/{ticketId} and records history.
func HandleUpdateTicket(queries ticketQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := getUserID(r)
		if !ok {
			httperr.Unauthorized("認証が必要です").Write(w)
			return
		}
		id, ok := parseTicketID(r)
		if !ok {
			httperr.BadRequest("INVALID_TICKET_ID", "チケットIDが不正です").Write(w)
			return
		}

		// AssigneeID uses **int64 to distinguish three states:
		//   nil      = field absent from JSON → keep current value
		//   &nil     = JSON null              → clear assignee
		//   &(&n)    = JSON number n          → set to n
		var req struct {
			Type        string  `json:"type"`
			Priority    string  `json:"priority"`
			Title       string  `json:"title"`
			Description *string `json:"description"`
			AssigneeID  **int64 `json:"assignee_id"`
			ReceivedAt  *string `json:"received_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperr.BadRequest("INVALID_REQUEST", "リクエストの形式が不正です").Write(w)
			return
		}

		details := map[string]string{}
		if req.Title == "" {
			details["title"] = "必須項目です"
		}
		if len(req.Title) > 255 {
			details["title"] = "255文字以内で入力してください"
		}
		if !validTicketType(req.Type) {
			details["type"] = "無効な種別です"
		}
		if !validTicketPriority(req.Priority) {
			details["priority"] = "無効な優先度です"
		}
		if len(details) > 0 {
			httperr.UnprocessableEntity("VALIDATION_FAILED", "入力内容に誤りがあります").
				WithDetails(details).Write(w)
			return
		}

		current, err := queries.GetTicket(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.NotFound("TICKET_NOT_FOUND", "指定されたチケットは存在しません").Write(w)
				return
			}
			slog.ErrorContext(
				r.Context(), "ticket update: get current",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		params := db.UpdateTicketParams{
			ID:         id,
			Title:      req.Title,
			Type:       db.TicketsType(req.Type),
			Priority:   db.TicketsPriority(req.Priority),
			AssigneeID: current.AssigneeID,
			ReceivedAt: current.ReceivedAt,
		}
		if current.Description.Valid {
			params.Description = current.Description
		}
		if req.Description != nil {
			params.Description = sql.NullString{String: *req.Description, Valid: true}
		} else if req.Description == nil && params.Description.Valid {
			// explicit null clears description
			params.Description = sql.NullString{}
		}
		switch {
		case req.AssigneeID == nil:
			// Field absent: keep current value.
			params.AssigneeID = current.AssigneeID
		case *req.AssigneeID == nil:
			// Explicit null: clear assignee.
			params.AssigneeID = sql.NullInt64{}
		default:
			// Validate new assignee is in the system's pool.
			count, err := queries.IsUserInSystemAssignees(r.Context(), db.IsUserInSystemAssigneesParams{
				SystemID: current.SystemID,
				UserID:   **req.AssigneeID,
			})
			if err != nil {
				slog.ErrorContext(
					r.Context(), "ticket update: check assignee pool",
					slog.Any("error", err),
					slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
				)
				httperr.InternalError("サーバーエラーが発生しました").Write(w)
				return
			}
			if count == 0 {
				httperr.UnprocessableEntity("ASSIGNEE_NOT_IN_POOL",
					"指定された担当者はこのシステムの担当者プールに含まれていません").Write(w)
				return
			}
			params.AssigneeID = sql.NullInt64{Int64: **req.AssigneeID, Valid: true}
		}
		if req.ReceivedAt != nil {
			if t, err := time.Parse(time.RFC3339, *req.ReceivedAt); err == nil {
				params.ReceivedAt = t
			}
		}

		if err := queries.UpdateTicket(r.Context(), params); err != nil {
			slog.ErrorContext(
				r.Context(), "ticket update: update",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		// Record history for changed fields.
		histFields := []struct {
			name    string
			oldVal  string
			newVal  string
			changed bool
		}{
			{"title", current.Title, req.Title, current.Title != req.Title},
			{"type", string(current.Type), req.Type, string(current.Type) != req.Type},
			{"priority", string(current.Priority), req.Priority, string(current.Priority) != req.Priority},
		}
		for _, f := range histFields {
			if !f.changed {
				continue
			}
			_ = queries.CreateTicketHistory(r.Context(), db.CreateTicketHistoryParams{
				TicketID:  id,
				ChangedBy: userID,
				FieldName: f.name,
				OldValue:  sql.NullString{String: f.oldVal, Valid: true},
				NewValue:  sql.NullString{String: f.newVal, Valid: true},
			})
		}

		row, err := queries.GetTicketDetail(r.Context(), id)
		if err != nil {
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		writeJSON(w, http.StatusOK, toDTO(row))
	}
}

// HandleUpdateStatus handles PATCH /tickets/{ticketId}/status.
func HandleUpdateStatus(queries ticketQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := getUserID(r)
		if !ok {
			httperr.Unauthorized("認証が必要です").Write(w)
			return
		}
		id, ok := parseTicketID(r)
		if !ok {
			httperr.BadRequest("INVALID_TICKET_ID", "チケットIDが不正です").Write(w)
			return
		}

		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperr.BadRequest("INVALID_REQUEST", "リクエストの形式が不正です").Write(w)
			return
		}

		nextStatus := db.TicketsStatus(req.Status)
		if !validTicketStatus(req.Status) {
			httperr.BadRequest("INVALID_STATUS", "無効なステータスです").
				WithDetails(map[string]string{"status": "無効なステータスです"}).Write(w)
			return
		}

		current, err := queries.GetTicket(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.NotFound("TICKET_NOT_FOUND", "指定されたチケットは存在しません").Write(w)
				return
			}
			slog.ErrorContext(
				r.Context(), "ticket status: get current",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		currentDomain := domain.TicketStatus(current.Status)
		nextDomain := domain.TicketStatus(nextStatus)
		if !currentDomain.IsValidTransition(nextDomain) {
			httperr.New(
				"INVALID_STATUS_TRANSITION",
				fmt.Sprintf("%s から %s への遷移は許可されていません", current.Status, req.Status),
				http.StatusUnprocessableEntity,
			).Write(w)
			return
		}

		if err := queries.UpdateTicketStatus(r.Context(), db.UpdateTicketStatusParams{
			ID:     id,
			Status: nextStatus,
		}); err != nil {
			slog.ErrorContext(
				r.Context(), "ticket status: update",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		_ = queries.CreateTicketHistory(r.Context(), db.CreateTicketHistoryParams{
			TicketID:  id,
			ChangedBy: userID,
			FieldName: "status",
			OldValue:  sql.NullString{String: string(current.Status), Valid: true},
			NewValue:  sql.NullString{String: req.Status, Valid: true},
		})

		row, err := queries.GetTicketDetail(r.Context(), id)
		if err != nil {
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		writeJSON(w, http.StatusOK, toDTO(row))
	}
}

// HandleUpdateAssignee handles PATCH /tickets/{ticketId}/assignee.
func HandleUpdateAssignee(queries ticketQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := getUserID(r)
		if !ok {
			httperr.Unauthorized("認証が必要です").Write(w)
			return
		}
		id, ok := parseTicketID(r)
		if !ok {
			httperr.BadRequest("INVALID_TICKET_ID", "チケットIDが不正です").Write(w)
			return
		}

		var req struct {
			AssigneeID *int64 `json:"assignee_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperr.BadRequest("INVALID_REQUEST", "リクエストの形式が不正です").Write(w)
			return
		}

		current, err := queries.GetTicket(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.NotFound("TICKET_NOT_FOUND", "指定されたチケットは存在しません").Write(w)
				return
			}
			slog.ErrorContext(
				r.Context(), "ticket assignee: get current",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		newAssignee := sql.NullInt64{}
		if req.AssigneeID != nil {
			count, err := queries.IsUserInSystemAssignees(r.Context(), db.IsUserInSystemAssigneesParams{
				SystemID: current.SystemID,
				UserID:   *req.AssigneeID,
			})
			if err != nil {
				slog.ErrorContext(
					r.Context(), "ticket assignee: check assignee pool",
					slog.Any("error", err),
					slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
				)
				httperr.InternalError("サーバーエラーが発生しました").Write(w)
				return
			}
			if count == 0 {
				httperr.UnprocessableEntity("ASSIGNEE_NOT_IN_POOL",
					"指定された担当者はこのシステムの担当者プールに含まれていません").Write(w)
				return
			}
			newAssignee = sql.NullInt64{Int64: *req.AssigneeID, Valid: true}
		}

		if err := queries.UpdateTicketAssignee(r.Context(), db.UpdateTicketAssigneeParams{
			ID:         id,
			AssigneeID: newAssignee,
		}); err != nil {
			slog.ErrorContext(
				r.Context(), "ticket assignee: update",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		oldVal := ""
		if current.AssigneeID.Valid {
			oldVal = strconv.FormatInt(current.AssigneeID.Int64, 10)
		}
		newVal := ""
		if req.AssigneeID != nil {
			newVal = strconv.FormatInt(*req.AssigneeID, 10)
		}
		_ = queries.CreateTicketHistory(r.Context(), db.CreateTicketHistoryParams{
			TicketID:  id,
			ChangedBy: userID,
			FieldName: "assignee_id",
			OldValue:  sql.NullString{String: oldVal, Valid: oldVal != ""},
			NewValue:  sql.NullString{String: newVal, Valid: newVal != ""},
		})

		row, err := queries.GetTicketDetail(r.Context(), id)
		if err != nil {
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		writeJSON(w, http.StatusOK, toDTO(row))
	}
}
