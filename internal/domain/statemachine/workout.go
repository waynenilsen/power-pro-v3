package statemachine

// Workout states.
const (
	WorkoutInProgress State = "IN_PROGRESS"
	WorkoutCompleted  State = "COMPLETED"
	WorkoutAbandoned  State = "ABANDONED"
)

// workoutTransitions defines valid transitions for workout state machine.
// Valid transitions:
// - IN_PROGRESS -> COMPLETED (workout finished)
// - IN_PROGRESS -> ABANDONED (workout cancelled)
//
// Note: No transitions FROM COMPLETED or ABANDONED (terminal states)
var workoutTransitions = []Transition{
	{From: WorkoutInProgress, To: WorkoutCompleted},
	{From: WorkoutInProgress, To: WorkoutAbandoned},
}

// WorkoutStateMachine manages workout state transitions.
type WorkoutStateMachine struct {
	state State
}

// NewWorkoutStateMachine creates a new workout state machine with the given initial state.
func NewWorkoutStateMachine(initialState State) *WorkoutStateMachine {
	return &WorkoutStateMachine{state: initialState}
}

// CurrentState returns the current workout state.
func (sm *WorkoutStateMachine) CurrentState() State {
	return sm.state
}

// ValidTransitions returns all valid workout transitions.
func (sm *WorkoutStateMachine) ValidTransitions() []Transition {
	return workoutTransitions
}

// CanTransitionTo checks if a transition to the target state is valid from the current state.
func (sm *WorkoutStateMachine) CanTransitionTo(target State) bool {
	for _, t := range workoutTransitions {
		if t.From == sm.state && t.To == target {
			return true
		}
	}
	return false
}

// TransitionTo attempts to transition to the target state.
// Returns an InvalidTransitionError if the transition is not valid.
func (sm *WorkoutStateMachine) TransitionTo(target State) error {
	if !sm.CanTransitionTo(target) {
		return NewInvalidTransitionError(sm.state, target)
	}
	sm.state = target
	return nil
}

// ValidWorkoutStates returns all valid workout states.
func ValidWorkoutStates() []State {
	return []State{WorkoutInProgress, WorkoutCompleted, WorkoutAbandoned}
}

// IsValidWorkoutState checks if a state is a valid workout state.
func IsValidWorkoutState(s State) bool {
	for _, valid := range ValidWorkoutStates() {
		if s == valid {
			return true
		}
	}
	return false
}

// IsTerminalWorkoutState checks if a state is a terminal workout state
// (i.e., COMPLETED or ABANDONED - no further transitions possible).
func IsTerminalWorkoutState(s State) bool {
	return s == WorkoutCompleted || s == WorkoutAbandoned
}
