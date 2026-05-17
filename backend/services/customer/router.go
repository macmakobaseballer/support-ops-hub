package customer

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
)

// Router returns the customer-service sub-router mounted at /customers.
// M3 implements read-only endpoints. Full CRUD is added in M5.
func Router(queries *db.Queries) http.Handler {
	r := chi.NewRouter()
	r.Get("/", HandleListCustomers(queries))
	r.Get("/{customerId}", HandleGetCustomer(queries))
	return r
}
