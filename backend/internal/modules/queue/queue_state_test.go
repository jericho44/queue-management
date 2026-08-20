package queue_test

import (
	"testing"
)

type TicketStateMachine struct{}

func (s *TicketStateMachine) IsValidTransition(fromStatus, toStatus string) bool {
	validTransitions := map[string][]string{
		"WAITING":   {"CALLED", "CANCELLED", "TRANSFERRED"},
		"CALLED":    {"SERVING", "RECALLED", "SKIPPED", "NO_SHOW", "CANCELLED"},
		"SERVING":   {"COMPLETED", "SKIPPED", "NO_SHOW", "CANCELLED", "TRANSFERRED"},
		"COMPLETED": {},
		"SKIPPED":   {"WAITING"},
		"NO_SHOW":   {"WAITING"},
		"CANCELLED": {},
	}

	allowed, ok := validTransitions[fromStatus]
	if !ok {
		return false
	}

	for _, st := range allowed {
		if st == toStatus {
			return true
		}
	}
	return false
}

func TestTicketStateTransitions(t *testing.T) {
	sm := &TicketStateMachine{}

	if !sm.IsValidTransition("WAITING", "CALLED") {
		t.Errorf("Expected WAITING -> CALLED to be valid")
	}
	if !sm.IsValidTransition("CALLED", "SERVING") {
		t.Errorf("Expected CALLED -> SERVING to be valid")
	}
	if !sm.IsValidTransition("SERVING", "COMPLETED") {
		t.Errorf("Expected SERVING -> COMPLETED to be valid")
	}

	if sm.IsValidTransition("COMPLETED", "SERVING") {
		t.Errorf("Expected COMPLETED -> SERVING to be invalid")
	}
	if sm.IsValidTransition("CANCELLED", "CALLED") {
		t.Errorf("Expected CANCELLED -> CALLED to be invalid")
	}
}
