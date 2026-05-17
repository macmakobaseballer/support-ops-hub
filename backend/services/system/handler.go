package system

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
)

type systemQuerier interface {
	ListSystemsAll(ctx context.Context) ([]db.System, error)
	ListSystemsActive(ctx context.Context) ([]db.System, error)
	ListSystemsByCustomer(ctx context.Context, customerID int64) ([]db.System, error)
	ListSystemsActiveByCustomer(ctx context.Context, customerID int64) ([]db.System, error)
	GetSystem(ctx context.Context, id int64) (db.System, error)
	ListSystemAssigneesWithUsers(ctx context.Context, systemID int64) ([]db.ListSystemAssigneesWithUsersRow, error)
}

type systemDTO struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	CustomerID  int64            `json:"customer_id"`
	Description *string          `json:"description"`
	IsActive    bool             `json:"is_active"`
	Assignees   []userSummaryDTO `json:"assignees"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type userSummaryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func toDTO(s db.System) systemDTO {
	dto := systemDTO{
		ID:         s.ID,
		Name:       s.Name,
		CustomerID: s.CustomerID,
		IsActive:   s.IsActive,
		Assignees:  []userSummaryDTO{},
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
	if s.Description.Valid {
		dto.Description = &s.Description.String
	}
	return dto
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("system: json encode error", slog.Any("error", err))
	}
}

// HandleListSystems handles GET /systems?customer_id=X&is_active=true.
func HandleListSystems(queries systemQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		isActive := q.Get("is_active") == "true"
		customerIDStr := q.Get("customer_id")

		var systems []db.System
		var err error

		if customerIDStr != "" {
			cid, cerr := strconv.ParseInt(customerIDStr, 10, 64)
			if cerr != nil || cid <= 0 {
				httperr.BadRequest("INVALID_CUSTOMER_ID", "顧客IDが不正です").Write(w)
				return
			}
			if isActive {
				systems, err = queries.ListSystemsActiveByCustomer(r.Context(), cid)
			} else {
				systems, err = queries.ListSystemsByCustomer(r.Context(), cid)
			}
		} else {
			if isActive {
				systems, err = queries.ListSystemsActive(r.Context())
			} else {
				systems, err = queries.ListSystemsAll(r.Context())
			}
		}
		if err != nil {
			slog.ErrorContext(
				r.Context(), "system list: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		dtos := make([]systemDTO, 0, len(systems))
		for _, s := range systems {
			dtos = append(dtos, toDTO(s))
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"data": dtos,
			"pagination": map[string]int{
				"total":       len(dtos),
				"page":        1,
				"per_page":    len(dtos),
				"total_pages": 1,
			},
		})
	}
}

// HandleGetSystem handles GET /systems/{systemId}.
func HandleGetSystem(queries systemQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "systemId"), 10, 64)
		if err != nil || id <= 0 {
			httperr.BadRequest("INVALID_SYSTEM_ID", "システムIDが不正です").Write(w)
			return
		}

		s, err := queries.GetSystem(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.NotFound("SYSTEM_NOT_FOUND", "指定されたシステムは存在しません").Write(w)
				return
			}
			slog.ErrorContext(
				r.Context(), "system get: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		writeJSON(w, http.StatusOK, toDTO(s))
	}
}

// HandleListAssignees handles GET /systems/{systemId}/assignees.
func HandleListAssignees(queries systemQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "systemId"), 10, 64)
		if err != nil || id <= 0 {
			httperr.BadRequest("INVALID_SYSTEM_ID", "システムIDが不正です").Write(w)
			return
		}

		assignees, err := queries.ListSystemAssigneesWithUsers(r.Context(), id)
		if err != nil {
			slog.ErrorContext(
				r.Context(), "system assignees: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		dtos := make([]userSummaryDTO, 0, len(assignees))
		for _, a := range assignees {
			dtos = append(dtos, userSummaryDTO{ID: a.ID, Name: a.Name})
		}
		writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
	}
}
