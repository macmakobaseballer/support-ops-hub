package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubAuthQuerier is a test double for authQuerier.
type stubAuthQuerier struct {
	byEmail map[string]db.User
	byID    map[int64]db.User
}

func (s *stubAuthQuerier) GetUserByEmail(_ context.Context, email string) (db.User, error) {
	u, ok := s.byEmail[email]
	if !ok {
		return db.User{}, sql.ErrNoRows
	}
	return u, nil
}

func (s *stubAuthQuerier) GetUser(_ context.Context, id int64) (db.User, error) {
	u, ok := s.byID[id]
	if !ok {
		return db.User{}, sql.ErrNoRows
	}
	return u, nil
}

var (
	activeAdminUser = db.User{
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

func newStub() *stubAuthQuerier {
	return &stubAuthQuerier{
		byEmail: map[string]db.User{
			activeAdminUser.Email: activeAdminUser,
			inactiveUser.Email:    inactiveUser,
		},
		byID: map[int64]db.User{
			activeAdminUser.ID: activeAdminUser,
			inactiveUser.ID:    inactiveUser,
		},
	}
}

// --- HandleLogin tests ---

func TestHandleLogin_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	body := `{"email":"admin@example.com","password":"any"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))

	HandleLogin(newStub())(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp loginResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, "dev-token:admin@example.com", resp.Token)
	assert.Equal(t, "admin@example.com", resp.User.Email)
	assert.Equal(t, "admin", resp.User.Role)
}

func TestHandleLogin_UnknownEmail_Returns401(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login",
		strings.NewReader(`{"email":"ghost@example.com","password":"x"}`))

	HandleLogin(newStub())(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleLogin_InactiveUser_Returns403(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login",
		strings.NewReader(`{"email":"inactive@example.com","password":"x"}`))

	HandleLogin(newStub())(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestHandleLogin_EmptyEmail_Returns400(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login",
		strings.NewReader(`{"email":"","password":"x"}`))

	HandleLogin(newStub())(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleLogin_InvalidJSON_Returns400(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login",
		strings.NewReader(`not-json`))

	HandleLogin(newStub())(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// --- HandleMe tests ---

func TestHandleMe_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("X-User-ID", "1")

	HandleMe(newStub())(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var dto userDTO
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&dto))
	assert.Equal(t, "admin@example.com", dto.Email)
}

func TestHandleMe_NoHeader_Returns401(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)

	HandleMe(newStub())(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleMe_InvalidID_Returns401(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("X-User-ID", "not-a-number")

	HandleMe(newStub())(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleMe_UnknownID_Returns401(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("X-User-ID", "999")

	HandleMe(newStub())(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleMe_InactiveUser_Returns403(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("X-User-ID", "2")

	HandleMe(newStub())(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}
