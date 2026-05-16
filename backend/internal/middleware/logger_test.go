package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequestID_GeneratesWhenAbsent(t *testing.T) {
	var capturedID string
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = middleware.RequestIDFromContext(r.Context())
	}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.NotEmpty(t, capturedID)
	assert.Equal(t, capturedID, rr.Header().Get("X-Request-Id"))
}

func TestRequestID_PreservesExistingHeader(t *testing.T) {
	const existingID = "test-request-id-123"
	var capturedID string
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = middleware.RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", existingID)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, existingID, capturedID)
	assert.Equal(t, existingID, rr.Header().Get("X-Request-Id"))
}

func TestRequestIDFromContext_EmptyWhenNotSet(t *testing.T) {
	assert.Empty(t, middleware.RequestIDFromContext(context.Background()))
}
