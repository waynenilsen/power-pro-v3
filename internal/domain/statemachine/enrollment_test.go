package statemachine

import (
	"errors"
	"testing"
)

func TestEnrollmentStateMachine_CurrentState(t *testing.T) {
	tests := []struct {
		name         string
		initialState State
		expected     State
	}{
		{"active", EnrollmentActive, EnrollmentActive},
		{"between cycles", EnrollmentBetweenCycles, EnrollmentBetweenCycles},
		{"quit", EnrollmentQuit, EnrollmentQuit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewEnrollmentStateMachine(tt.initialState)
			if sm.CurrentState() != tt.expected {
				t.Errorf("expected state %s, got %s", tt.expected, sm.CurrentState())
			}
		})
	}
}

func TestEnrollmentStateMachine_ValidTransitions(t *testing.T) {
	sm := NewEnrollmentStateMachine(EnrollmentActive)
	transitions := sm.ValidTransitions()

	if len(transitions) != 4 {
		t.Errorf("expected 4 transitions, got %d", len(transitions))
	}
}

func TestEnrollmentStateMachine_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name     string
		from     State
		to       State
		expected bool
	}{
		// Valid transitions from ACTIVE
		{"active to between cycles", EnrollmentActive, EnrollmentBetweenCycles, true},
		{"active to quit", EnrollmentActive, EnrollmentQuit, true},
		{"active to active", EnrollmentActive, EnrollmentActive, false},

		// Valid transitions from BETWEEN_CYCLES
		{"between cycles to active", EnrollmentBetweenCycles, EnrollmentActive, true},
		{"between cycles to quit", EnrollmentBetweenCycles, EnrollmentQuit, true},
		{"between cycles to between cycles", EnrollmentBetweenCycles, EnrollmentBetweenCycles, false},

		// QUIT is terminal - no valid transitions
		{"quit to active", EnrollmentQuit, EnrollmentActive, false},
		{"quit to between cycles", EnrollmentQuit, EnrollmentBetweenCycles, false},
		{"quit to quit", EnrollmentQuit, EnrollmentQuit, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewEnrollmentStateMachine(tt.from)
			if sm.CanTransitionTo(tt.to) != tt.expected {
				t.Errorf("CanTransitionTo(%s) from %s = %v, expected %v",
					tt.to, tt.from, sm.CanTransitionTo(tt.to), tt.expected)
			}
		})
	}
}

func TestEnrollmentStateMachine_TransitionTo(t *testing.T) {
	tests := []struct {
		name        string
		from        State
		to          State
		expectError bool
	}{
		// Valid transitions
		{"active to between cycles", EnrollmentActive, EnrollmentBetweenCycles, false},
		{"active to quit", EnrollmentActive, EnrollmentQuit, false},
		{"between cycles to active", EnrollmentBetweenCycles, EnrollmentActive, false},
		{"between cycles to quit", EnrollmentBetweenCycles, EnrollmentQuit, false},

		// Invalid transitions
		{"active to active", EnrollmentActive, EnrollmentActive, true},
		{"quit to active", EnrollmentQuit, EnrollmentActive, true},
		{"quit to between cycles", EnrollmentQuit, EnrollmentBetweenCycles, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewEnrollmentStateMachine(tt.from)
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

func TestValidEnrollmentStates(t *testing.T) {
	states := ValidEnrollmentStates()
	if len(states) != 3 {
		t.Errorf("expected 3 states, got %d", len(states))
	}

	expected := map[State]bool{
		EnrollmentActive:        true,
		EnrollmentBetweenCycles: true,
		EnrollmentQuit:          true,
	}

	for _, s := range states {
		if !expected[s] {
			t.Errorf("unexpected state %s", s)
		}
	}
}

func TestIsValidEnrollmentState(t *testing.T) {
	tests := []struct {
		state    State
		expected bool
	}{
		{EnrollmentActive, true},
		{EnrollmentBetweenCycles, true},
		{EnrollmentQuit, true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if IsValidEnrollmentState(tt.state) != tt.expected {
				t.Errorf("IsValidEnrollmentState(%s) = %v, expected %v",
					tt.state, IsValidEnrollmentState(tt.state), tt.expected)
			}
		})
	}
}

func TestEnrollmentStateMachine_ImplementsInterface(t *testing.T) {
	var _ StateMachine = (*EnrollmentStateMachine)(nil)
}
