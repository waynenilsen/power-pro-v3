package statemachine

// Enrollment states.
const (
	EnrollmentActive        State = "ACTIVE"
	EnrollmentBetweenCycles State = "BETWEEN_CYCLES"
	EnrollmentQuit          State = "QUIT"
)

// enrollmentTransitions defines valid transitions for enrollment state machine.
// Valid transitions:
// - ACTIVE -> BETWEEN_CYCLES (cycle completed)
// - ACTIVE -> QUIT (user quits)
// - BETWEEN_CYCLES -> ACTIVE (new cycle started)
// - BETWEEN_CYCLES -> QUIT (user quits)
var enrollmentTransitions = []Transition{
	{From: EnrollmentActive, To: EnrollmentBetweenCycles},
	{From: EnrollmentActive, To: EnrollmentQuit},
	{From: EnrollmentBetweenCycles, To: EnrollmentActive},
	{From: EnrollmentBetweenCycles, To: EnrollmentQuit},
}

// EnrollmentStateMachine manages enrollment state transitions.
type EnrollmentStateMachine struct {
	state State
}

// NewEnrollmentStateMachine creates a new enrollment state machine with the given initial state.
func NewEnrollmentStateMachine(initialState State) *EnrollmentStateMachine {
	return &EnrollmentStateMachine{state: initialState}
}

// CurrentState returns the current enrollment state.
func (sm *EnrollmentStateMachine) CurrentState() State {
	return sm.state
}

// ValidTransitions returns all valid enrollment transitions.
func (sm *EnrollmentStateMachine) ValidTransitions() []Transition {
	return enrollmentTransitions
}

// CanTransitionTo checks if a transition to the target state is valid from the current state.
func (sm *EnrollmentStateMachine) CanTransitionTo(target State) bool {
	for _, t := range enrollmentTransitions {
		if t.From == sm.state && t.To == target {
			return true
		}
	}
	return false
}

// TransitionTo attempts to transition to the target state.
// Returns an InvalidTransitionError if the transition is not valid.
func (sm *EnrollmentStateMachine) TransitionTo(target State) error {
	if !sm.CanTransitionTo(target) {
		return NewInvalidTransitionError(sm.state, target)
	}
	sm.state = target
	return nil
}

// ValidEnrollmentStates returns all valid enrollment states.
func ValidEnrollmentStates() []State {
	return []State{EnrollmentActive, EnrollmentBetweenCycles, EnrollmentQuit}
}

// IsValidEnrollmentState checks if a state is a valid enrollment state.
func IsValidEnrollmentState(s State) bool {
	for _, valid := range ValidEnrollmentStates() {
		if s == valid {
			return true
		}
	}
	return false
}
