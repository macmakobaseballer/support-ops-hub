package httperr_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_SetsFields(t *testing.T) {
	e := httperr.New("TEST_CODE", "テストメッセージ", http.StatusBadRequest)
	assert.Equal(t, "TEST_CODE", e.Code)
	assert.Equal(t, "テストメッセージ", e.Message)
}

func TestWrite_StatusAndContentType(t *testing.T) {
	rr := httptest.NewRecorder()
	httperr.NotFound("RESOURCE_NOT_FOUND", "見つかりません").Write(rr)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestWrite_JSONBody(t *testing.T) {
	rr := httptest.NewRecorder()
	httperr.BadRequest("VALIDATION_FAILED", "入力エラー").
		WithDetails(map[string]string{"email": "必須項目です"}).Write(rr)

	var body map[string]any
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	assert.Equal(t, "VALIDATION_FAILED", body["code"])
	assert.Equal(t, "入力エラー", body["message"])
	details, ok := body["details"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "必須項目です", details["email"])
}

func TestUnauthorized_Returns401(t *testing.T) {
	rr := httptest.NewRecorder()
	httperr.Unauthorized("認証が必要です").Write(rr)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestForbidden_Returns403(t *testing.T) {
	rr := httptest.NewRecorder()
	httperr.Forbidden("権限がありません").Write(rr)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestInternalError_Returns500(t *testing.T) {
	rr := httptest.NewRecorder()
	httperr.InternalError("サーバーエラー").Write(rr)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
