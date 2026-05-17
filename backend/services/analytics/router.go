package analytics

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
)

// Router returns the analytics-service sub-router mounted at /analytics.
func Router(queries *db.Queries) http.Handler {
	r := chi.NewRouter()
	r.Get("/dashboard", HandleDashboard(queries))
	r.Get("/customers/{customerId}/report", HandleCustomerReport(queries))
	return r
}
