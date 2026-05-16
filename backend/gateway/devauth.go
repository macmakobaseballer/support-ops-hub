// Package gateway implements the API Gateway with dev-auth middleware.
// In M7, DevAuthMiddleware is replaced by JWT Bearer token verification.
// Services continue reading the same X-User-ID / X-User-Role headers unchanged.
package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
)

const (
	// HeaderDevUserEmail is the dev-auth header sent by the frontend.
	HeaderDevUserEmail = "X-Dev-User-Email"
	// HeaderUserID is the injected user ID header consumed by services.
	HeaderUserID = "X-User-ID"
	// HeaderUserRole is the injected user role header consumed by services.
	HeaderUserRole = "X-User-Role"
)

// userLookup is a minimal interface over *db.Queries for testability.
type userLookup interface {
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
}

// DevAuthMiddleware translates X-Dev-User-Email into X-User-ID / X-User-Role.
type DevAuthMiddleware struct {
	queries userLookup
}

// NewDevAuthMiddleware creates a DevAuthMiddleware backed by the given queries.
func NewDevAuthMiddleware(queries userLookup) *DevAuthMiddleware {
	return &DevAuthMiddleware{queries: queries}
}

// Handler is the http.Handler middleware function.
func (m *DevAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get(HeaderDevUserEmail)
		if email == "" {
			httperr.Unauthorized("認証が必要です").Write(w)
			return
		}

		user, err := m.queries.GetUserByEmail(r.Context(), email)
		if err != nil {
			slog.WarnContext(r.Context(), "dev-auth: user not found",
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.Unauthorized("ユーザーが見つかりません").Write(w)
			return
		}

		if !user.IsActive {
			httperr.Forbidden("このアカウントは無効です").Write(w)
			return
		}

		r = r.Clone(r.Context())
		r.Header.Set(HeaderUserID, fmt.Sprintf("%d", user.ID))
		r.Header.Set(HeaderUserRole, string(user.Role))

		next.ServeHTTP(w, r)
	})
}
