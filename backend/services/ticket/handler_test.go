package ticket

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── test doubles ────────────────────────────────────────────────────────────

// mockTicketQuerier implements ticketQuerier for unit tests.
type mockTicketQuerier struct {
	ticket       db.Ticket
	ticketDetail db.GetTicketDetailRow
	getErr       error
	createResult sql.Result
	createErr    error
	updateErr    error
}

func (m *mockTicketQuerier) GetTicket(_ context.Context, _ int64) (db.Ticket, error) {
	return m.ticket, m.getErr
}
func (m *mockTicketQuerier) GetTicketDetail(_ context.Context, _ int64) (db.GetTicketDetailRow, error) {
	return m.ticketDetail, m.getErr
}
func (m *mockTicketQuerier) CreateTicket(_ context.Context, _ db.CreateTicketParams) (sql.Result, error) {
	return m.createResult, m.createErr
}
func (m *mockTicketQuerier) UpdateTicket(_ context.Context, _ db.UpdateTicketParams) error {
	return m.updateErr
}
func (m *mockTicketQuerier) UpdateTicketStatus(_ context.Context, _ db.UpdateTicketStatusParams) error {
	return m.updateErr
}
func (m *mockTicketQuerier) UpdateTicketAssignee(_ context.Context, _ db.UpdateTicketAssigneeParams) error {
	return m.updateErr
}
func (m *mockTicketQuerier) CreateTicketHistory(_ context.Context, _ db.CreateTicketHistoryParams) error {
	return nil
}

// mockResult implements sql.Result for CreateTicket tests.
type mockResult struct{ id int64 }

func (m mockResult) LastInsertId() (int64, error) { return m.id, nil }
func (m mockResult) RowsAffected() (int64, error) { return 1, nil }

// chiRequest builds a request with chi URL params.
func chiRequest(method, path string, body []byte, params map[string]string) *http.Request {
	r := httptest.NewRequest(method, path, bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-User-ID", "1")

	if len(params) > 0 {
		rctx := chi.NewRouteContext()
		for k, v := range params {
			rctx.URLParams.Add(k, v)
		}
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	}
	return r
}

func sampleDetail() db.GetTicketDetailRow {
	return db.GetTicketDetailRow{
		ID:            1,
		Title:         "テストチケット",
		Type:          db.TicketsTypeQuestion,
		Priority:      db.TicketsPriorityMedium,
		Status:        db.TicketsStatusNew,
		CustomerID:    1,
		CustomerName:  "テスト企業",
		SystemID:      1,
		SystemName:    "基幹システム",
		CreatedBy:     1,
		CreatedByName: "山田太郎",
		ReceivedAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func sampleTicket() db.Ticket {
	return db.Ticket{
		ID:       1,
		Title:    "テストチケット",
		Type:     db.TicketsTypeQuestion,
		Priority: db.TicketsPriorityMedium,
		Status:   db.TicketsStatusNew,
	}
}

// ─── GetTicket tests ─────────────────────────────────────────────────────────

func TestGetTicket_Found_Returns200(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{ticketDetail: sampleDetail()}
	h := HandleGetTicket(q)

	r := chiRequest(http.MethodGet, "/tickets/1", nil, map[string]string{"ticketId": "1"})
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var body ticketDTO
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Equal(t, int64(1), body.ID)
	assert.Equal(t, "テストチケット", body.Title)
}

func TestGetTicket_NotFound_HasErrorCode(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{getErr: sql.ErrNoRows}
	h := HandleGetTicket(q)

	r := chiRequest(http.MethodGet, "/tickets/999", nil, map[string]string{"ticketId": "999"})
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var body map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Equal(t, "TICKET_NOT_FOUND", body["code"])
}

// ─── CreateTicket tests ───────────────────────────────────────────────────────

func TestCreateTicket_Valid_Returns201(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{
		createResult: mockResult{id: 1},
		ticketDetail: sampleDetail(),
	}
	h := HandleCreateTicket(q)

	body, _ := json.Marshal(map[string]any{
		"customer_id": 1, "system_id": 1,
		"type": "question", "priority": "medium",
		"title": "テストチケット",
	})
	r := chiRequest(http.MethodPost, "/tickets", body, nil)
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp ticketDTO
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, int64(1), resp.ID)
}

func TestCreateTicket_MissingTitle_Returns400(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{}
	h := HandleCreateTicket(q)

	body, _ := json.Marshal(map[string]any{
		"customer_id": 1, "system_id": 1,
		"type": "question", "priority": "medium",
		// title missing
	})
	r := chiRequest(http.MethodPost, "/tickets", body, nil)
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "VALIDATION_FAILED", resp["code"])
	details, _ := resp["details"].(map[string]any)
	assert.NotEmpty(t, details["title"])
}

func TestCreateTicket_InvalidType_Returns400(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{}
	h := HandleCreateTicket(q)

	body, _ := json.Marshal(map[string]any{
		"customer_id": 1, "system_id": 1,
		"type": "invalid_type", "priority": "medium",
		"title": "テスト",
	})
	r := chiRequest(http.MethodPost, "/tickets", body, nil)
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreateTicket_NoUserID_Returns401(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{}
	h := HandleCreateTicket(q)

	body, _ := json.Marshal(map[string]any{
		"customer_id": 1, "system_id": 1,
		"type": "question", "priority": "medium", "title": "テスト",
	})
	r := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	// No X-User-ID header
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── UpdateStatus tests — status transitions (ルール5 全パターン) ────────────

func TestUpdateStatus_ValidTransitions(t *testing.T) {
	t.Parallel()
	cases := []struct {
		from db.TicketsStatus
		to   string
	}{
		{db.TicketsStatusNew, "in_progress"},
		{db.TicketsStatusInProgress, "waiting"},
		{db.TicketsStatusInProgress, "done"},
		{db.TicketsStatusWaiting, "in_progress"},
		{db.TicketsStatusWaiting, "done"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.from)+"->"+tc.to, func(t *testing.T) {
			t.Parallel()
			detail := sampleDetail()
			detail.Status = tc.from
			q := &mockTicketQuerier{
				ticket:       db.Ticket{ID: 1, Status: tc.from},
				ticketDetail: detail,
			}
			h := HandleUpdateStatus(q)

			body, _ := json.Marshal(map[string]string{"status": tc.to})
			r := chiRequest(http.MethodPatch, "/tickets/1/status", body, map[string]string{"ticketId": "1"})
			w := httptest.NewRecorder()
			h(w, r)

			assert.Equal(t, http.StatusOK, w.Code, "from=%s to=%s", tc.from, tc.to)
		})
	}
}

func TestUpdateStatus_InvalidTransitions_Returns422(t *testing.T) {
	t.Parallel()
	cases := []struct {
		from db.TicketsStatus
		to   string
	}{
		{db.TicketsStatusNew, "waiting"},
		{db.TicketsStatusNew, "done"},
		{db.TicketsStatusDone, "new"},
		{db.TicketsStatusDone, "in_progress"},
		{db.TicketsStatusDone, "waiting"},
		{db.TicketsStatusInProgress, "new"},
		{db.TicketsStatusWaiting, "new"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.from)+"->"+tc.to, func(t *testing.T) {
			t.Parallel()
			q := &mockTicketQuerier{
				ticket:       db.Ticket{ID: 1, Status: tc.from},
				ticketDetail: sampleDetail(),
			}
			h := HandleUpdateStatus(q)

			body, _ := json.Marshal(map[string]string{"status": tc.to})
			r := chiRequest(http.MethodPatch, "/tickets/1/status", body, map[string]string{"ticketId": "1"})
			w := httptest.NewRecorder()
			h(w, r)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "from=%s to=%s", tc.from, tc.to)
			var resp map[string]string
			require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			assert.Equal(t, "INVALID_STATUS_TRANSITION", resp["code"])
		})
	}
}

func TestUpdateStatus_InvalidStatusValue_Returns400(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{ticket: sampleTicket()}
	h := HandleUpdateStatus(q)

	body, _ := json.Marshal(map[string]string{"status": "unknown_status"})
	r := chiRequest(http.MethodPatch, "/tickets/1/status", body, map[string]string{"ticketId": "1"})
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "INVALID_STATUS", resp["code"])
}

func TestUpdateStatus_NoUserID_Returns401(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{}
	h := HandleUpdateStatus(q)

	body, _ := json.Marshal(map[string]string{"status": "in_progress"})
	r := httptest.NewRequest(http.MethodPatch, "/tickets/1/status", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ticketId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateStatus_NotFound_Returns404(t *testing.T) {
	t.Parallel()
	q := &mockTicketQuerier{getErr: sql.ErrNoRows}
	h := HandleUpdateStatus(q)

	body, _ := json.Marshal(map[string]string{"status": "in_progress"})
	r := chiRequest(http.MethodPatch, "/tickets/999/status", body, map[string]string{"ticketId": "999"})
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── UpdateAssignee tests ─────────────────────────────────────────────────────

func TestUpdateAssignee_NullAssignee_Returns200(t *testing.T) {
	t.Parallel()
	detail := sampleDetail()
	q := &mockTicketQuerier{
		ticket:       sampleTicket(),
		ticketDetail: detail,
	}
	h := HandleUpdateAssignee(q)

	body, _ := json.Marshal(map[string]any{"assignee_id": nil})
	r := chiRequest(http.MethodPatch, "/tickets/1/assignee", body, map[string]string{"ticketId": "1"})
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
