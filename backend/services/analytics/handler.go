package analytics

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
)

// analyticsQuerier is the minimal interface for analytics handlers.
type analyticsQuerier interface {
	GetAllStatusCounts(ctx context.Context) (db.GetAllStatusCountsRow, error)
	GetMonthlyTrend(ctx context.Context) ([]db.GetMonthlyTrendRow, error)
	GetOpenTypeBreakdown(ctx context.Context) ([]db.GetOpenTypeBreakdownRow, error)
	GetOpenPriorityBreakdown(ctx context.Context) ([]db.GetOpenPriorityBreakdownRow, error)
	GetTopSystemsByOpenCount(ctx context.Context) ([]db.GetTopSystemsByOpenCountRow, error)
	GetCustomer(ctx context.Context, id int64) (db.Customer, error)
	GetCustomerAllStatusCounts(ctx context.Context, customerID int64) (db.GetCustomerAllStatusCountsRow, error)
	GetCustomerMonthlyTrend(ctx context.Context, customerID int64) ([]db.GetCustomerMonthlyTrendRow, error)
	GetCustomerOpenTypeBreakdown(ctx context.Context, customerID int64) ([]db.GetCustomerOpenTypeBreakdownRow, error)
	GetCustomerSystemsBreakdown(ctx context.Context, customerID int64) ([]db.GetCustomerSystemsBreakdownRow, error)
}

// ─── response DTOs ────────────────────────────────────────────────────────────

type statusCountsDTO struct {
	New        int64 `json:"new"`
	InProgress int64 `json:"in_progress"`
	Waiting    int64 `json:"waiting"`
	Done       int64 `json:"done"`
}

type monthlyTrendItemDTO struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type typeCountsDTO struct {
	Question int64 `json:"question"`
	Bug      int64 `json:"bug"`
	Config   int64 `json:"config"`
	Data     int64 `json:"data"`
}

type priorityCountsDTO struct {
	High   int64 `json:"high"`
	Medium int64 `json:"medium"`
	Low    int64 `json:"low"`
}

type topSystemItemDTO struct {
	SystemID   int64  `json:"system_id"`
	SystemName string `json:"system_name"`
	Count      int64  `json:"count"`
}

type dashboardSummaryDTO struct {
	CountsByStatus       statusCountsDTO       `json:"counts_by_status"`
	MonthlyTrend         []monthlyTrendItemDTO `json:"monthly_trend"`
	OpenCountsByType     typeCountsDTO         `json:"open_counts_by_type"`
	OpenCountsByPriority priorityCountsDTO     `json:"open_counts_by_priority"`
	TopSystemsByOpenCount []topSystemItemDTO   `json:"top_systems_by_open_count"`
}

type systemBreakdownItemDTO struct {
	SystemID   int64           `json:"system_id"`
	SystemName string          `json:"system_name"`
	Counts     statusCountsDTO `json:"counts"`
	Total      int64           `json:"total"`
}

type customerReportDTO struct {
	CustomerID      int64                   `json:"customer_id"`
	CustomerName    string                  `json:"customer_name"`
	CountsByStatus  statusCountsDTO         `json:"counts_by_status"`
	MonthlyTrend    []monthlyTrendItemDTO   `json:"monthly_trend"`
	OpenCountsByType typeCountsDTO          `json:"open_counts_by_type"`
	SystemBreakdown []systemBreakdownItemDTO `json:"system_breakdown"`
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func buildMonthlyTrend(rows []db.GetMonthlyTrendRow) []monthlyTrendItemDTO {
	items := make([]monthlyTrendItemDTO, 0, len(rows))
	for _, r := range rows {
		items = append(items, monthlyTrendItemDTO{Month: r.Month, Count: r.Count})
	}
	return items
}

func buildCustomerMonthlyTrend(rows []db.GetCustomerMonthlyTrendRow) []monthlyTrendItemDTO {
	items := make([]monthlyTrendItemDTO, 0, len(rows))
	for _, r := range rows {
		items = append(items, monthlyTrendItemDTO{Month: r.Month, Count: r.Count})
	}
	return items
}

func buildTypeCountsFromRows(rows []db.GetOpenTypeBreakdownRow) typeCountsDTO {
	var tc typeCountsDTO
	for _, r := range rows {
		switch r.Type {
		case db.TicketsTypeQuestion:
			tc.Question = r.Count
		case db.TicketsTypeBug:
			tc.Bug = r.Count
		case db.TicketsTypeConfig:
			tc.Config = r.Count
		case db.TicketsTypeData:
			tc.Data = r.Count
		}
	}
	return tc
}

func buildTypeCountsFromCustomerRows(rows []db.GetCustomerOpenTypeBreakdownRow) typeCountsDTO {
	var tc typeCountsDTO
	for _, r := range rows {
		switch r.Type {
		case db.TicketsTypeQuestion:
			tc.Question = r.Count
		case db.TicketsTypeBug:
			tc.Bug = r.Count
		case db.TicketsTypeConfig:
			tc.Config = r.Count
		case db.TicketsTypeData:
			tc.Data = r.Count
		}
	}
	return tc
}

func buildPriorityCountsFromRows(rows []db.GetOpenPriorityBreakdownRow) priorityCountsDTO {
	var pc priorityCountsDTO
	for _, r := range rows {
		switch r.Priority {
		case db.TicketsPriorityHigh:
			pc.High = r.Count
		case db.TicketsPriorityMedium:
			pc.Medium = r.Count
		case db.TicketsPriorityLow:
			pc.Low = r.Count
		}
	}
	return pc
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("analytics: json encode error", slog.Any("error", err))
	}
}

// ─── handlers ────────────────────────────────────────────────────────────────

// HandleDashboard handles GET /analytics/dashboard.
func HandleDashboard(queries analyticsQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID := middleware.RequestIDFromContext(ctx)

		counts, err := queries.GetAllStatusCounts(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "analytics dashboard: status counts", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		trend, err := queries.GetMonthlyTrend(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "analytics dashboard: monthly trend", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		typeRows, err := queries.GetOpenTypeBreakdown(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "analytics dashboard: type breakdown", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		priorityRows, err := queries.GetOpenPriorityBreakdown(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "analytics dashboard: priority breakdown", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		topSystems, err := queries.GetTopSystemsByOpenCount(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "analytics dashboard: top systems", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		topItems := make([]topSystemItemDTO, 0, len(topSystems))
		for _, s := range topSystems {
			topItems = append(topItems, topSystemItemDTO{
				SystemID:   s.SystemID,
				SystemName: s.SystemName,
				Count:      s.Count,
			})
		}

		writeJSON(w, http.StatusOK, dashboardSummaryDTO{
			CountsByStatus: statusCountsDTO{
				New:        counts.NewCount,
				InProgress: counts.InProgressCount,
				Waiting:    counts.WaitingCount,
				Done:       counts.DoneCount,
			},
			MonthlyTrend:          buildMonthlyTrend(trend),
			OpenCountsByType:      buildTypeCountsFromRows(typeRows),
			OpenCountsByPriority:  buildPriorityCountsFromRows(priorityRows),
			TopSystemsByOpenCount: topItems,
		})
	}
}

// HandleCustomerReport handles GET /analytics/customers/{customerId}/report.
func HandleCustomerReport(queries analyticsQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID := middleware.RequestIDFromContext(ctx)

		customerID, err := strconv.ParseInt(chi.URLParam(r, "customerId"), 10, 64)
		if err != nil || customerID <= 0 {
			httperr.BadRequest("INVALID_CUSTOMER_ID", "顧客IDが不正です").Write(w)
			return
		}

		customer, err := queries.GetCustomer(ctx, customerID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.NotFound("CUSTOMER_NOT_FOUND", "指定された顧客企業は存在しません").Write(w)
				return
			}
			slog.ErrorContext(ctx, "analytics customer report: get customer", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		counts, err := queries.GetCustomerAllStatusCounts(ctx, customerID)
		if err != nil {
			slog.ErrorContext(ctx, "analytics customer report: status counts", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		trend, err := queries.GetCustomerMonthlyTrend(ctx, customerID)
		if err != nil {
			slog.ErrorContext(ctx, "analytics customer report: monthly trend", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		typeRows, err := queries.GetCustomerOpenTypeBreakdown(ctx, customerID)
		if err != nil {
			slog.ErrorContext(ctx, "analytics customer report: type breakdown", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		sysRows, err := queries.GetCustomerSystemsBreakdown(ctx, customerID)
		if err != nil {
			slog.ErrorContext(ctx, "analytics customer report: system breakdown", slog.Any("error", err), slog.String("request_id", reqID))
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		sysBreakdown := make([]systemBreakdownItemDTO, 0, len(sysRows))
		for _, s := range sysRows {
			sysBreakdown = append(sysBreakdown, systemBreakdownItemDTO{
				SystemID:   s.SystemID,
				SystemName: s.SystemName,
				Counts: statusCountsDTO{
					New:        s.NewCount,
					InProgress: s.InProgressCount,
					Waiting:    s.WaitingCount,
					Done:       s.DoneCount,
				},
				Total: s.Total,
			})
		}

		writeJSON(w, http.StatusOK, customerReportDTO{
			CustomerID:   customerID,
			CustomerName: customer.Name,
			CountsByStatus: statusCountsDTO{
				New:        counts.NewCount,
				InProgress: counts.InProgressCount,
				Waiting:    counts.WaitingCount,
				Done:       counts.DoneCount,
			},
			MonthlyTrend:    buildCustomerMonthlyTrend(trend),
			OpenCountsByType: buildTypeCountsFromCustomerRows(typeRows),
			SystemBreakdown: sysBreakdown,
		})
	}
}
