// Package statemachine provides generic state machine interfaces and implementations
// for managing state transitions in the PowerPro domain.
package statemachine

import "fmt"

// State represents any state in a state machine.
type State string

// Transition represents a valid state transition.
type Transition struct {
	From State
	To   State
}

// StateMachine defines the interface for all state machines.
type StateMachine interface {
	// CurrentState returns the current state.
	CurrentState() State

	// ValidTransitions returns all valid transitions from all states.
	ValidTransitions() []Transition

	// CanTransitionTo checks if a transition to the target state is valid.
	CanTransitionTo(target State) bool

	// TransitionTo attempts to transition to the target state.
	// Returns an error if the transition is invalid.
	TransitionTo(target State) error
}

// InvalidTransitionError represents an invalid state transition attempt.
type InvalidTransitionError struct {
	From State
	To   State
}

// Error implements the error interface.
func (e *InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid transition from %s to %s", e.From, e.To)
}

// NewInvalidTransitionError creates a new InvalidTransitionError.
func NewInvalidTransitionError(from, to State) *InvalidTransitionError {
	return &InvalidTransitionError{From: from, To: to}
}
