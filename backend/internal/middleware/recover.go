package middleware

import (
	"log/slog"
	"net/http"

	"github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr"
)

// Recover catches panics, logs them, and returns a 500 response.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "panic recovered",
					slog.Any("panic", rec),
					slog.String("request_id", RequestIDFromContext(r.Context())),
				)
				httperr.InternalError("サーバーエラーが発生しました").Write(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
