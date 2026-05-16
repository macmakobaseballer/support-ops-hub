package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string  `json:"token"`
	User  userDTO `json:"user"`
}

type userDTO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// HandleLogin returns a dev stub token for M2.
// M7 replaces this with bcrypt password verification + real JWT signing.
func HandleLogin(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httperr.BadRequest("INVALID_REQUEST", "リクエストの形式が不正です").Write(w)
			return
		}
		if req.Email == "" {
			httperr.BadRequest("VALIDATION_FAILED", "メールアドレスは必須です").
				WithDetails(map[string]string{"email": "必須項目です"}).Write(w)
			return
		}

		user, err := queries.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.Unauthorized("メールアドレスまたはパスワードが正しくありません").Write(w)
				return
			}
			slog.ErrorContext(r.Context(), "login: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		if !user.IsActive {
			httperr.Unauthorized("このアカウントは無効です").Write(w)
			return
		}

		resp := loginResponse{
			Token: "dev-token:" + req.Email,
			User: userDTO{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Role:  string(user.Role),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.ErrorContext(r.Context(), "login: encode error", slog.Any("error", err))
		}
	}
}

// HandleLogout is a no-op. Clients discard the token locally (CLAUDE.md 実装上の注意事項5).
func HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// HandleMe reads X-User-ID (injected by dev-auth middleware) and returns the user.
func HandleMe(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.Header.Get("X-User-ID")
		if idStr == "" {
			httperr.Unauthorized("認証が必要です").Write(w)
			return
		}

		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httperr.Unauthorized("無効なユーザーIDです").Write(w)
			return
		}

		user, err := queries.GetUser(r.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httperr.Unauthorized("ユーザーが見つかりません").Write(w)
				return
			}
			slog.ErrorContext(r.Context(), "me: db error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
			httperr.InternalError("サーバーエラーが発生しました").Write(w)
			return
		}

		dto := userDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dto); err != nil {
			slog.ErrorContext(r.Context(), "me: encode error",
				slog.Any("error", err),
				slog.String("request_id", middleware.RequestIDFromContext(r.Context())),
			)
		}
	}
}
