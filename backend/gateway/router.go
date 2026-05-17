package gateway

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	appmid "github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
	analyticsvc "github.com/macmakobaseballer/support-ops-hub/backend/services/analytics"
	authsvc "github.com/macmakobaseballer/support-ops-hub/backend/services/auth"
	customersvc "github.com/macmakobaseballer/support-ops-hub/backend/services/customer"
	systemsvc "github.com/macmakobaseballer/support-ops-hub/backend/services/system"
	ticketsvc "github.com/macmakobaseballer/support-ops-hub/backend/services/ticket"
)

// NewRouter builds the chi router with all routes and middleware for the gateway.
// sqlDB is passed alongside queries for services that need dynamic SQL (e.g. ticket list).
// allowedOrigins is a comma-separated list of origins (e.g. "http://localhost:3000").
func NewRouter(sqlDB *sql.DB, queries *db.Queries, allowedOrigins string) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: strings.Split(allowedOrigins, ","),
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept", "Authorization", "Content-Type",
			HeaderDevUserEmail, "X-Request-Id",
		},
		MaxAge: 300,
	}))
	r.Use(appmid.RequestID)
	r.Use(appmid.Logger)
	r.Use(appmid.Recover)
	r.Use(chimid.Timeout(25 * time.Second))

	devAuth := NewDevAuthMiddleware(queries)

	r.Route("/api/v1", func(r chi.Router) {
		// Unauthenticated: login does not require a user header
		r.Post("/auth/login", authsvc.HandleLogin(queries))

		// Authenticated: all other routes require dev-auth header
		r.Group(func(r chi.Router) {
			r.Use(devAuth.Handler)
			r.Post("/auth/logout", authsvc.HandleLogout())
			r.Get("/auth/me", authsvc.HandleMe(queries))

			// M3 services
			r.Mount("/tickets",   ticketsvc.Router(sqlDB, queries))
			r.Mount("/analytics", analyticsvc.Router(queries))
			r.Mount("/customers", customersvc.Router(queries))
			r.Mount("/systems",   systemsvc.Router(queries))
		})
	})

	return r
}
