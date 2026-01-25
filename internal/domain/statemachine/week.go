package statemachine

// Week states.
const (
	WeekPending    State = "PENDING"
	WeekInProgress State = "IN_PROGRESS"
	WeekCompleted  State = "COMPLETED"
)

// weekTransitions defines valid transitions for week state machine.
// Valid transitions:
// - PENDING -> IN_PROGRESS (week started)
// - IN_PROGRESS -> COMPLETED (all workouts done)
// - COMPLETED -> PENDING (reset for next week)
var weekTransitions = []Transition{
	{From: WeekPending, To: WeekInProgress},
	{From: WeekInProgress, To: WeekCompleted},
	{From: WeekCompleted, To: WeekPending},
}

// WeekStateMachine manages week state transitions.
type WeekStateMachine struct {
	state State
}

// NewWeekStateMachine creates a new week state machine with the given initial state.
func NewWeekStateMachine(initialState State) *WeekStateMachine {
	return &WeekStateMachine{state: initialState}
}

// CurrentState returns the current week state.
func (sm *WeekStateMachine) CurrentState() State {
	return sm.state
}

// ValidTransitions returns all valid week transitions.
func (sm *WeekStateMachine) ValidTransitions() []Transition {
	return weekTransitions
}

// CanTransitionTo checks if a transition to the target state is valid from the current state.
func (sm *WeekStateMachine) CanTransitionTo(target State) bool {
	for _, t := range weekTransitions {
		if t.From == sm.state && t.To == target {
			return true
		}
	}
	return false
}

// TransitionTo attempts to transition to the target state.
// Returns an InvalidTransitionError if the transition is not valid.
func (sm *WeekStateMachine) TransitionTo(target State) error {
	if !sm.CanTransitionTo(target) {
		return NewInvalidTransitionError(sm.state, target)
	}
	sm.state = target
	return nil
}

// ValidWeekStates returns all valid week states.
func ValidWeekStates() []State {
	return []State{WeekPending, WeekInProgress, WeekCompleted}
}

// IsValidWeekState checks if a state is a valid week state.
func IsValidWeekState(s State) bool {
	for _, valid := range ValidWeekStates() {
		if s == valid {
			return true
		}
	}
	return false
}
