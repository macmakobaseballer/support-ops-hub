package domain_test

import (
	"testing"

	"github.com/macmakobaseballer/support-ops-hub/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTicketStatus_IsValidTransition(t *testing.T) {
	tests := []struct {
		from  domain.TicketStatus
		to    domain.TicketStatus
		valid bool
	}{
		// new
		{domain.TicketStatusNew, domain.TicketStatusInProgress, true},
		{domain.TicketStatusNew, domain.TicketStatusWaiting, false},
		{domain.TicketStatusNew, domain.TicketStatusDone, false},
		{domain.TicketStatusNew, domain.TicketStatusNew, false},
		// in_progress
		{domain.TicketStatusInProgress, domain.TicketStatusWaiting, true},
		{domain.TicketStatusInProgress, domain.TicketStatusDone, true},
		{domain.TicketStatusInProgress, domain.TicketStatusNew, false},
		{domain.TicketStatusInProgress, domain.TicketStatusInProgress, false},
		// waiting
		{domain.TicketStatusWaiting, domain.TicketStatusInProgress, true},
		{domain.TicketStatusWaiting, domain.TicketStatusDone, true},
		{domain.TicketStatusWaiting, domain.TicketStatusNew, false},
		{domain.TicketStatusWaiting, domain.TicketStatusWaiting, false},
		// done (terminal)
		{domain.TicketStatusDone, domain.TicketStatusNew, false},
		{domain.TicketStatusDone, domain.TicketStatusInProgress, false},
		{domain.TicketStatusDone, domain.TicketStatusWaiting, false},
		{domain.TicketStatusDone, domain.TicketStatusDone, false},
		// unknown
		{"unknown", domain.TicketStatusInProgress, false},
	}

	for _, tt := range tests {
		name := string(tt.from) + "->" + string(tt.to)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.from.IsValidTransition(tt.to))
		})
	}
}
