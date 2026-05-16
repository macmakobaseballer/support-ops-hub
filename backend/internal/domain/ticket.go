package domain

// TicketStatus represents the lifecycle state of a ticket (CLAUDE.md ルール5).
type TicketStatus string

const (
	TicketStatusNew        TicketStatus = "new"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusWaiting    TicketStatus = "waiting"
	TicketStatusDone       TicketStatus = "done"
)

// validTransitions encodes the allowed status transitions.
// Violations must return HTTP 422 INVALID_STATUS_TRANSITION.
var validTransitions = map[TicketStatus][]TicketStatus{
	TicketStatusNew:        {TicketStatusInProgress},
	TicketStatusInProgress: {TicketStatusWaiting, TicketStatusDone},
	TicketStatusWaiting:    {TicketStatusInProgress, TicketStatusDone},
	TicketStatusDone:       {},
}

// IsValidTransition reports whether transitioning from s to next is allowed.
func (s TicketStatus) IsValidTransition(next TicketStatus) bool {
	targets, ok := validTransitions[s]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == next {
			return true
		}
	}
	return false
}
