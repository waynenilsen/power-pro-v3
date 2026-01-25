package statemachine

// Cycle states.
const (
	CyclePending    State = "PENDING"
	CycleInProgress State = "IN_PROGRESS"
	CycleCompleted  State = "COMPLETED"
)

// cycleTransitions defines valid transitions for cycle state machine.
// Valid transitions:
// - PENDING -> IN_PROGRESS (cycle started)
// - IN_PROGRESS -> COMPLETED (all weeks done)
// - COMPLETED -> PENDING (reset for new cycle)
var cycleTransitions = []Transition{
	{From: CyclePending, To: CycleInProgress},
	{From: CycleInProgress, To: CycleCompleted},
	{From: CycleCompleted, To: CyclePending},
}

// CycleStateMachine manages cycle state transitions.
type CycleStateMachine struct {
	state State
}

// NewCycleStateMachine creates a new cycle state machine with the given initial state.
func NewCycleStateMachine(initialState State) *CycleStateMachine {
	return &CycleStateMachine{state: initialState}
}

// CurrentState returns the current cycle state.
func (sm *CycleStateMachine) CurrentState() State {
	return sm.state
}

// ValidTransitions returns all valid cycle transitions.
func (sm *CycleStateMachine) ValidTransitions() []Transition {
	return cycleTransitions
}

// CanTransitionTo checks if a transition to the target state is valid from the current state.
func (sm *CycleStateMachine) CanTransitionTo(target State) bool {
	for _, t := range cycleTransitions {
		if t.From == sm.state && t.To == target {
			return true
		}
	}
	return false
}

// TransitionTo attempts to transition to the target state.
// Returns an InvalidTransitionError if the transition is not valid.
func (sm *CycleStateMachine) TransitionTo(target State) error {
	if !sm.CanTransitionTo(target) {
		return NewInvalidTransitionError(sm.state, target)
	}
	sm.state = target
	return nil
}

// ValidCycleStates returns all valid cycle states.
func ValidCycleStates() []State {
	return []State{CyclePending, CycleInProgress, CycleCompleted}
}

// IsValidCycleState checks if a state is a valid cycle state.
func IsValidCycleState(s State) bool {
	for _, valid := range ValidCycleStates() {
		if s == valid {
			return true
		}
	}
	return false
}
