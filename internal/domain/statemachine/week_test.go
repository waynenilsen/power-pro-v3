package statemachine

import (
	"errors"
	"testing"
)

func TestWeekStateMachine_CurrentState(t *testing.T) {
	tests := []struct {
		name         string
		initialState State
		expected     State
	}{
		{"pending", WeekPending, WeekPending},
		{"in progress", WeekInProgress, WeekInProgress},
		{"completed", WeekCompleted, WeekCompleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewWeekStateMachine(tt.initialState)
			if sm.CurrentState() != tt.expected {
				t.Errorf("expected state %s, got %s", tt.expected, sm.CurrentState())
			}
		})
	}
}

func TestWeekStateMachine_ValidTransitions(t *testing.T) {
	sm := NewWeekStateMachine(WeekPending)
	transitions := sm.ValidTransitions()

	if len(transitions) != 3 {
		t.Errorf("expected 3 transitions, got %d", len(transitions))
	}
}

func TestWeekStateMachine_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name     string
		from     State
		to       State
		expected bool
	}{
		// Valid transitions from PENDING
		{"pending to in progress", WeekPending, WeekInProgress, true},
		{"pending to completed", WeekPending, WeekCompleted, false},
		{"pending to pending", WeekPending, WeekPending, false},

		// Valid transitions from IN_PROGRESS
		{"in progress to completed", WeekInProgress, WeekCompleted, true},
		{"in progress to pending", WeekInProgress, WeekPending, false},
		{"in progress to in progress", WeekInProgress, WeekInProgress, false},

		// Valid transitions from COMPLETED
		{"completed to pending", WeekCompleted, WeekPending, true},
		{"completed to in progress", WeekCompleted, WeekInProgress, false},
		{"completed to completed", WeekCompleted, WeekCompleted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewWeekStateMachine(tt.from)
			if sm.CanTransitionTo(tt.to) != tt.expected {
				t.Errorf("CanTransitionTo(%s) from %s = %v, expected %v",
					tt.to, tt.from, sm.CanTransitionTo(tt.to), tt.expected)
			}
		})
	}
}

func TestWeekStateMachine_TransitionTo(t *testing.T) {
	tests := []struct {
		name        string
		from        State
		to          State
		expectError bool
	}{
		// Valid transitions
		{"pending to in progress", WeekPending, WeekInProgress, false},
		{"in progress to completed", WeekInProgress, WeekCompleted, false},
		{"completed to pending", WeekCompleted, WeekPending, false},

		// Invalid transitions
		{"pending to completed", WeekPending, WeekCompleted, true},
		{"in progress to pending", WeekInProgress, WeekPending, true},
		{"completed to in progress", WeekCompleted, WeekInProgress, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewWeekStateMachine(tt.from)
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

func TestWeekStateMachine_FullWeek(t *testing.T) {
	sm := NewWeekStateMachine(WeekPending)

	// Start the week
	if err := sm.TransitionTo(WeekInProgress); err != nil {
		t.Fatalf("failed to start week: %v", err)
	}
	if sm.CurrentState() != WeekInProgress {
		t.Errorf("expected IN_PROGRESS, got %s", sm.CurrentState())
	}

	// Complete the week
	if err := sm.TransitionTo(WeekCompleted); err != nil {
		t.Fatalf("failed to complete week: %v", err)
	}
	if sm.CurrentState() != WeekCompleted {
		t.Errorf("expected COMPLETED, got %s", sm.CurrentState())
	}

	// Reset for next week
	if err := sm.TransitionTo(WeekPending); err != nil {
		t.Fatalf("failed to reset week: %v", err)
	}
	if sm.CurrentState() != WeekPending {
		t.Errorf("expected PENDING, got %s", sm.CurrentState())
	}
}

func TestValidWeekStates(t *testing.T) {
	states := ValidWeekStates()
	if len(states) != 3 {
		t.Errorf("expected 3 states, got %d", len(states))
	}

	expected := map[State]bool{
		WeekPending:    true,
		WeekInProgress: true,
		WeekCompleted:  true,
	}

	for _, s := range states {
		if !expected[s] {
			t.Errorf("unexpected state %s", s)
		}
	}
}

func TestIsValidWeekState(t *testing.T) {
	tests := []struct {
		state    State
		expected bool
	}{
		{WeekPending, true},
		{WeekInProgress, true},
		{WeekCompleted, true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if IsValidWeekState(tt.state) != tt.expected {
				t.Errorf("IsValidWeekState(%s) = %v, expected %v",
					tt.state, IsValidWeekState(tt.state), tt.expected)
			}
		})
	}
}

func TestWeekStateMachine_ImplementsInterface(t *testing.T) {
	var _ StateMachine = (*WeekStateMachine)(nil)
}
