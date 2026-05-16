package gateway

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/stretchr/testify/assert"
)

// stubUserLookup is a test double for the userLookup interface.
type stubUserLookup struct {
	user db.User
	err  error
}

func (s *stubUserLookup) GetUserByEmail(_ context.Context, _ string) (db.User, error) {
	return s.user, s.err
}

var (
	activeAdmin = db.User{
		ID: 1, Name: "管理者", Email: "admin@example.com",
		Role: db.UsersRoleAdmin, IsActive: true,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	inactiveUser = db.User{
		ID: 2, Name: "無効ユーザー", Email: "inactive@example.com",
		Role: db.UsersRoleMember, IsActive: false,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
)

func newDevAuthHandler(stub *stubUserLookup) http.Handler {
	m := NewDevAuthMiddleware(stub)
	return m.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Got-User-ID", r.Header.Get(HeaderUserID))
		w.Header().Set("X-Got-User-Role", r.Header.Get(HeaderUserRole))
		w.WriteHeader(http.StatusOK)
	}))
}

func TestDevAuth_MissingHeader_Returns401(t *testing.T) {
	h := newDevAuthHandler(&stubUserLookup{user: activeAdmin})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDevAuth_UserNotFound_Returns401(t *testing.T) {
	h := newDevAuthHandler(&stubUserLookup{err: sql.ErrNoRows})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderDevUserEmail, "ghost@example.com")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDevAuth_InactiveUser_Returns403(t *testing.T) {
	h := newDevAuthHandler(&stubUserLookup{user: inactiveUser})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderDevUserEmail, inactiveUser.Email)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDevAuth_ValidUser_InjectsHeaders(t *testing.T) {
	h := newDevAuthHandler(&stubUserLookup{user: activeAdmin})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderDevUserEmail, activeAdmin.Email)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "1", rr.Header().Get("X-Got-User-ID"))
	assert.Equal(t, "admin", rr.Header().Get("X-Got-User-Role"))
}
