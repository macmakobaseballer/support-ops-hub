package customer

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

type customerQuerier interface {
	ListCustomersAll(ctx context.Context) ([]db.Customer, error)
	ListCustomersActive(ctx context.Context) ([]db.Customer, error)
	GetCustomer(ctx context.Context, id int64) (db.Customer, error)
}

type customerDTO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Notes     *string   `json:"notes"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toDTO(c db.Customer) customerDTO {
	dto := customerDTO{
		ID:        c.ID,
		Name:      c.Name,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.Notes.Valid {
		dto.Notes = &c.Notes.String
	}
	return dto
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("customer: json encode error", slog.Any("error", err))
	}
}

// HandleListCustomers handles GET /customers?is_active=true|false.
func HandleListCustomers(queries customerQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isActiveParam := r.URL.Query().Get("is_active")

		var customers []db.Customer
		var err error
		if isActiveParam == "true" {
			customers, err = queries.ListCustomersActive(r.Context())
		} else {
			customers, err = queries.ListCustomersAll(r.Context())
		}
		if err != nil {
			slog.ErrorContext(
				r.Context(), "customer list: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		dtos := make([]customerDTO, 0, len(customers))
		for _, c := range customers {
			dtos = append(dtos, toDTO(c))
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

// HandleGetCustomer handles GET /customers/{customerId}.
func HandleGetCustomer(queries customerQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "customerId"), 10, 64)
		if err != nil || id <= 0 {
			httperr.BadRequest("INVALID_CUSTOMER_ID", "顧客IDが不正です").Write(w)
			return
		}

		c, err := queries.GetCustomer(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.NotFound("CUSTOMER_NOT_FOUND", "指定された顧客企業は存在しません").Write(w)
				return
			}
			slog.ErrorContext(
				r.Context(), "customer get: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}
		writeJSON(w, http.StatusOK, toDTO(c))
	}
}
