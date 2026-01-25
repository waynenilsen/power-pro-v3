// Package errors provides state-specific error types for enrollment and workout operations.
package errors

// StateError represents a structured error response for state transition failures.
// These errors provide machine-readable codes and human-readable messages
// with optional details for debugging.
type StateError struct {
	// Category is the base error category for HTTP status code mapping.
	Category error
	// Code is a machine-readable error code (snake_case).
	Code string
	// Message is the human-readable error message.
	Message string
	// Details contains additional context for debugging.
	Details map[string]interface{}
}

// Error implements the error interface.
func (e *StateError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error for errors.Is and errors.As support.
func (e *StateError) Unwrap() error {
	return e.Category
}

// Is implements error comparison for errors.Is support.
func (e *StateError) Is(target error) bool {
	return target == e.Category
}

// GetCode returns the machine-readable error code.
func (e *StateError) GetCode() string {
	return e.Code
}

// GetDetails returns the error details.
func (e *StateError) GetDetails() map[string]interface{} {
	return e.Details
}

// State error codes
const (
	CodeWorkoutAlreadyInProgress = "workout_already_in_progress"
	CodeNoActiveWorkout          = "no_active_workout"
	CodeInvalidEnrollmentState   = "invalid_enrollment_state"
	CodeSessionNotActive         = "session_not_active"
	CodeNotEnrolled              = "not_enrolled"
)

// NewWorkoutAlreadyInProgress creates an error for when a workout is already in progress.
func NewWorkoutAlreadyInProgress(currentSessionID string) *StateError {
	return &StateError{
		Category: ErrConflict,
		Code:     CodeWorkoutAlreadyInProgress,
		Message:  "Complete or abandon current workout before starting a new one",
		Details: map[string]interface{}{
			"current_workout_session_id": currentSessionID,
		},
	}
}

// NewNoActiveWorkout creates an error for when no workout is in progress.
func NewNoActiveWorkout() *StateError {
	return &StateError{
		Category: ErrBadRequest,
		Code:     CodeNoActiveWorkout,
		Message:  "No workout is currently in progress",
	}
}

// NewInvalidEnrollmentState creates an error for invalid enrollment state transitions.
func NewInvalidEnrollmentState(requiredState, currentStatus string) *StateError {
	return &StateError{
		Category: ErrBadRequest,
		Code:     CodeInvalidEnrollmentState,
		Message:  "Cannot " + requiredState + " - enrollment is not in the required state",
		Details: map[string]interface{}{
			"current_status": currentStatus,
		},
	}
}

// NewSessionNotActive creates an error for operations on non-active sessions.
func NewSessionNotActive(sessionStatus string) *StateError {
	return &StateError{
		Category: ErrBadRequest,
		Code:     CodeSessionNotActive,
		Message:  "Cannot log sets to a session that is not in progress",
		Details: map[string]interface{}{
			"session_status": sessionStatus,
		},
	}
}

// NewNotEnrolled creates an error for when a user is not enrolled in a program.
func NewNotEnrolled() *StateError {
	return &StateError{
		Category: ErrBadRequest,
		Code:     CodeNotEnrolled,
		Message:  "User must be enrolled in a program to perform this action",
	}
}

// IsStateError checks if an error is a StateError.
func IsStateError(err error) bool {
	_, ok := err.(*StateError)
	return ok
}

// GetStateErrorCode extracts the code from a StateError.
// Returns empty string if the error is not a StateError.
func GetStateErrorCode(err error) string {
	if stateErr, ok := err.(*StateError); ok {
		return stateErr.Code
	}
	return ""
}

// GetStateErrorDetails extracts the details from a StateError.
// Returns nil if the error is not a StateError.
func GetStateErrorDetails(err error) map[string]interface{} {
	if stateErr, ok := err.(*StateError); ok {
		return stateErr.Details
	}
	return nil
}
