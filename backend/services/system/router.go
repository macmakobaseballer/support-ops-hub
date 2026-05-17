package system

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
)

// Router returns the system-service sub-router mounted at /systems.
// M3 implements read-only endpoints. Full CRUD is added in M5.
func Router(queries *db.Queries) http.Handler {
	r := chi.NewRouter()
	r.Get("/", HandleListSystems(queries))
	r.Get("/{systemId}", HandleGetSystem(queries))
	r.Get("/{systemId}/assignees", HandleListAssignees(queries))
	return r
}
