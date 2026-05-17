package ticket

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
)

// Router returns the ticket-service sub-router mounted at /tickets.
// sqlDB is needed for the dynamic list query (parameterised filters).
func Router(sqlDB *sql.DB, queries *db.Queries) http.Handler {
	r := chi.NewRouter()
	r.Get("/", HandleListTickets(sqlDB, queries))
	r.Post("/", HandleCreateTicket(queries))
	r.Get("/{ticketId}", HandleGetTicket(queries))
	r.Put("/{ticketId}", HandleUpdateTicket(queries))
	r.Patch("/{ticketId}/status", HandleUpdateStatus(queries))
	r.Patch("/{ticketId}/assignee", HandleUpdateAssignee(queries))
	return r
}
