package statemachine

import (
	"errors"
	"testing"
)

func TestCycleStateMachine_CurrentState(t *testing.T) {
	tests := []struct {
		name         string
		initialState State
		expected     State
	}{
		{"pending", CyclePending, CyclePending},
		{"in progress", CycleInProgress, CycleInProgress},
		{"completed", CycleCompleted, CycleCompleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewCycleStateMachine(tt.initialState)
			if sm.CurrentState() != tt.expected {
				t.Errorf("expected state %s, got %s", tt.expected, sm.CurrentState())
			}
		})
	}
}

func TestCycleStateMachine_ValidTransitions(t *testing.T) {
	sm := NewCycleStateMachine(CyclePending)
	transitions := sm.ValidTransitions()

	if len(transitions) != 3 {
		t.Errorf("expected 3 transitions, got %d", len(transitions))
	}
}

func TestCycleStateMachine_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name     string
		from     State
		to       State
		expected bool
	}{
		// Valid transitions from PENDING
		{"pending to in progress", CyclePending, CycleInProgress, true},
		{"pending to completed", CyclePending, CycleCompleted, false},
		{"pending to pending", CyclePending, CyclePending, false},

		// Valid transitions from IN_PROGRESS
		{"in progress to completed", CycleInProgress, CycleCompleted, true},
		{"in progress to pending", CycleInProgress, CyclePending, false},
		{"in progress to in progress", CycleInProgress, CycleInProgress, false},

		// Valid transitions from COMPLETED
		{"completed to pending", CycleCompleted, CyclePending, true},
		{"completed to in progress", CycleCompleted, CycleInProgress, false},
		{"completed to completed", CycleCompleted, CycleCompleted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewCycleStateMachine(tt.from)
			if sm.CanTransitionTo(tt.to) != tt.expected {
				t.Errorf("CanTransitionTo(%s) from %s = %v, expected %v",
					tt.to, tt.from, sm.CanTransitionTo(tt.to), tt.expected)
			}
		})
	}
}

func TestCycleStateMachine_TransitionTo(t *testing.T) {
	tests := []struct {
		name        string
		from        State
		to          State
		expectError bool
	}{
		// Valid transitions
		{"pending to in progress", CyclePending, CycleInProgress, false},
		{"in progress to completed", CycleInProgress, CycleCompleted, false},
		{"completed to pending", CycleCompleted, CyclePending, false},

		// Invalid transitions
		{"pending to completed", CyclePending, CycleCompleted, true},
		{"in progress to pending", CycleInProgress, CyclePending, true},
		{"completed to in progress", CycleCompleted, CycleInProgress, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewCycleStateMachine(tt.from)
			err := sm.TransitionTo(tt.to)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				var invalidErr *InvalidTransitionError
				if !errors.As(err, &invalidErr) {
					t.Errorf("expected InvalidTransitionError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if sm.CurrentState() != tt.to {
					t.Errorf("expected state %s after transition, got %s", tt.to, sm.CurrentState())
				}
			}
		})
	}
}

func TestCycleStateMachine_FullCycle(t *testing.T) {
	sm := NewCycleStateMachine(CyclePending)

	// Start the cycle
	if err := sm.TransitionTo(CycleInProgress); err != nil {
		t.Fatalf("failed to start cycle: %v", err)
	}
	if sm.CurrentState() != CycleInProgress {
		t.Errorf("expected IN_PROGRESS, got %s", sm.CurrentState())
	}

	// Complete the cycle
	if err := sm.TransitionTo(CycleCompleted); err != nil {
		t.Fatalf("failed to complete cycle: %v", err)
	}
	if sm.CurrentState() != CycleCompleted {
		t.Errorf("expected COMPLETED, got %s", sm.CurrentState())
	}

	// Reset for new cycle
	if err := sm.TransitionTo(CyclePending); err != nil {
		t.Fatalf("failed to reset cycle: %v", err)
	}
	if sm.CurrentState() != CyclePending {
		t.Errorf("expected PENDING, got %s", sm.CurrentState())
	}
}

func TestValidCycleStates(t *testing.T) {
	states := ValidCycleStates()
	if len(states) != 3 {
		t.Errorf("expected 3 states, got %d", len(states))
	}

	expected := map[State]bool{
		CyclePending:    true,
		CycleInProgress: true,
		CycleCompleted:  true,
	}

	for _, s := range states {
		if !expected[s] {
			t.Errorf("unexpected state %s", s)
		}
	}
}

func TestIsValidCycleState(t *testing.T) {
	tests := []struct {
		state    State
		expected bool
	}{
		{CyclePending, true},
		{CycleInProgress, true},
		{CycleCompleted, true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if IsValidCycleState(tt.state) != tt.expected {
				t.Errorf("IsValidCycleState(%s) = %v, expected %v",
					tt.state, IsValidCycleState(tt.state), tt.expected)
			}
		})
	}
}

func TestCycleStateMachine_ImplementsInterface(t *testing.T) {
	var _ StateMachine = (*CycleStateMachine)(nil)
}
