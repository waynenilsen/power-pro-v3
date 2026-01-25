package statemachine

import (
	"errors"
	"testing"
)

func TestWorkoutStateMachine_CurrentState(t *testing.T) {
	tests := []struct {
		name         string
		initialState State
		expected     State
	}{
		{"in progress", WorkoutInProgress, WorkoutInProgress},
		{"completed", WorkoutCompleted, WorkoutCompleted},
		{"abandoned", WorkoutAbandoned, WorkoutAbandoned},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewWorkoutStateMachine(tt.initialState)
			if sm.CurrentState() != tt.expected {
				t.Errorf("expected state %s, got %s", tt.expected, sm.CurrentState())
			}
		})
	}
}

func TestWorkoutStateMachine_ValidTransitions(t *testing.T) {
	sm := NewWorkoutStateMachine(WorkoutInProgress)
	transitions := sm.ValidTransitions()

	if len(transitions) != 2 {
		t.Errorf("expected 2 transitions, got %d", len(transitions))
	}
}

func TestWorkoutStateMachine_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name     string
		from     State
		to       State
		expected bool
	}{
		// Valid transitions from IN_PROGRESS
		{"in progress to completed", WorkoutInProgress, WorkoutCompleted, true},
		{"in progress to abandoned", WorkoutInProgress, WorkoutAbandoned, true},
		{"in progress to in progress", WorkoutInProgress, WorkoutInProgress, false},

		// COMPLETED is terminal - no valid transitions
		{"completed to in progress", WorkoutCompleted, WorkoutInProgress, false},
		{"completed to abandoned", WorkoutCompleted, WorkoutAbandoned, false},
		{"completed to completed", WorkoutCompleted, WorkoutCompleted, false},

		// ABANDONED is terminal - no valid transitions
		{"abandoned to in progress", WorkoutAbandoned, WorkoutInProgress, false},
		{"abandoned to completed", WorkoutAbandoned, WorkoutCompleted, false},
		{"abandoned to abandoned", WorkoutAbandoned, WorkoutAbandoned, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewWorkoutStateMachine(tt.from)
			if sm.CanTransitionTo(tt.to) != tt.expected {
				t.Errorf("CanTransitionTo(%s) from %s = %v, expected %v",
					tt.to, tt.from, sm.CanTransitionTo(tt.to), tt.expected)
			}
		})
	}
}

func TestWorkoutStateMachine_TransitionTo(t *testing.T) {
	tests := []struct {
		name        string
		from        State
		to          State
		expectError bool
	}{
		// Valid transitions
		{"in progress to completed", WorkoutInProgress, WorkoutCompleted, false},
		{"in progress to abandoned", WorkoutInProgress, WorkoutAbandoned, false},

		// Invalid transitions - terminal states
		{"completed to in progress", WorkoutCompleted, WorkoutInProgress, true},
		{"completed to abandoned", WorkoutCompleted, WorkoutAbandoned, true},
		{"abandoned to in progress", WorkoutAbandoned, WorkoutInProgress, true},
		{"abandoned to completed", WorkoutAbandoned, WorkoutCompleted, true},

		// Invalid transitions - same state
		{"in progress to in progress", WorkoutInProgress, WorkoutInProgress, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewWorkoutStateMachine(tt.from)
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

func TestWorkoutStateMachine_CompleteWorkout(t *testing.T) {
	sm := NewWorkoutStateMachine(WorkoutInProgress)

	// Complete the workout
	if err := sm.TransitionTo(WorkoutCompleted); err != nil {
		t.Fatalf("failed to complete workout: %v", err)
	}
	if sm.CurrentState() != WorkoutCompleted {
		t.Errorf("expected COMPLETED, got %s", sm.CurrentState())
	}

	// Cannot transition from terminal state
	if err := sm.TransitionTo(WorkoutInProgress); err == nil {
		t.Error("expected error when transitioning from terminal state, got nil")
	}
}

func TestWorkoutStateMachine_AbandonWorkout(t *testing.T) {
	sm := NewWorkoutStateMachine(WorkoutInProgress)

	// Abandon the workout
	if err := sm.TransitionTo(WorkoutAbandoned); err != nil {
		t.Fatalf("failed to abandon workout: %v", err)
	}
	if sm.CurrentState() != WorkoutAbandoned {
		t.Errorf("expected ABANDONED, got %s", sm.CurrentState())
	}

	// Cannot transition from terminal state
	if err := sm.TransitionTo(WorkoutInProgress); err == nil {
		t.Error("expected error when transitioning from terminal state, got nil")
	}
}

func TestValidWorkoutStates(t *testing.T) {
	states := ValidWorkoutStates()
	if len(states) != 3 {
		t.Errorf("expected 3 states, got %d", len(states))
	}

	expected := map[State]bool{
		WorkoutInProgress: true,
		WorkoutCompleted:  true,
		WorkoutAbandoned:  true,
	}

	for _, s := range states {
		if !expected[s] {
			t.Errorf("unexpected state %s", s)
		}
	}
}

func TestIsValidWorkoutState(t *testing.T) {
	tests := []struct {
		state    State
		expected bool
	}{
		{WorkoutInProgress, true},
		{WorkoutCompleted, true},
		{WorkoutAbandoned, true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if IsValidWorkoutState(tt.state) != tt.expected {
				t.Errorf("IsValidWorkoutState(%s) = %v, expected %v",
					tt.state, IsValidWorkoutState(tt.state), tt.expected)
			}
		})
	}
}

func TestIsTerminalWorkoutState(t *testing.T) {
	tests := []struct {
		state    State
		expected bool
	}{
		{WorkoutInProgress, false},
		{WorkoutCompleted, true},
		{WorkoutAbandoned, true},
		{"INVALID", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if IsTerminalWorkoutState(tt.state) != tt.expected {
				t.Errorf("IsTerminalWorkoutState(%s) = %v, expected %v",
					tt.state, IsTerminalWorkoutState(tt.state), tt.expected)
			}
		})
	}
}

func TestWorkoutStateMachine_ImplementsInterface(t *testing.T) {
	var _ StateMachine = (*WorkoutStateMachine)(nil)
}
