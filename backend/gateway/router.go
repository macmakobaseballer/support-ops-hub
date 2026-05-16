package gateway

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	appmid "github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
	authsvc "github.com/macmakobaseballer/support-ops-hub/backend/services/auth"
)

// NewRouter builds the chi router with all routes and middleware for the gateway.
func NewRouter(queries *db.Queries) http.Handler {
	r := chi.NewRouter()

	r.Use(appmid.RequestID)
	r.Use(appmid.Logger)
	r.Use(appmid.Recover)
	r.Use(chimid.Timeout(30 * time.Second))

	devAuth := NewDevAuthMiddleware(queries)

	r.Route("/api/v1", func(r chi.Router) {
		// Unauthenticated: login does not require a user header
		r.Post("/auth/login", authsvc.HandleLogin(queries))

		// Authenticated: all other routes require dev-auth header
		r.Group(func(r chi.Router) {
			r.Use(devAuth.Handler)
			r.Post("/auth/logout", authsvc.HandleLogout())
			r.Get("/auth/me", authsvc.HandleMe(queries))
			// M3+ service routes will be mounted here:
			// r.Mount("/tickets",   ticketsvc.Router(queries))
			// r.Mount("/customers", customersvc.Router(queries))
			// r.Mount("/systems",   systemsvc.Router(queries))
		})
	})

	return r
}
