// Package progression provides domain logic for progression strategies.
// This file defines the FailureCounter entity for tracking consecutive failures.
package progression

import (
	"errors"
	"time"
)

// FailureCounter errors.
var (
	ErrFailureCounterUserIDRequired       = errors.New("failure counter user_id is required")
	ErrFailureCounterLiftIDRequired       = errors.New("failure counter lift_id is required")
	ErrFailureCounterProgressionIDRequired = errors.New("failure counter progression_id is required")
	ErrConsecutiveFailuresNegative        = errors.New("consecutive failures cannot be negative")
)

// FailureCounter tracks consecutive failures for a user's lift within a specific progression.
// This enables failure-based progression strategies like deloads after N consecutive failures.
type FailureCounter struct {
	// ID is the unique identifier for this failure counter.
	ID string
	// UserID is the UUID of the user.
	UserID string
	// LiftID is the UUID of the lift being tracked.
	LiftID string
	// ProgressionID is the UUID of the progression this counter is associated with.
	// This allows different progressions to have different failure thresholds.
	ProgressionID string
	// ConsecutiveFailures is the number of consecutive failures (resets on success).
	ConsecutiveFailures int
	// LastFailureAt is when the last failure occurred. Nil if no failures recorded.
	LastFailureAt *time.Time
	// LastSuccessAt is when the last success occurred. Nil if no successes recorded.
	LastSuccessAt *time.Time
	// CreatedAt is when this counter was created.
	CreatedAt time.Time
	// UpdatedAt is when this counter was last updated.
	UpdatedAt time.Time
}

// CreateFailureCounterInput contains the input data for creating a new failure counter.
type CreateFailureCounterInput struct {
	UserID        string
	LiftID        string
	ProgressionID string
}

// FailureCounterValidationResult holds validation errors for failure counter operations.
type FailureCounterValidationResult struct {
	Valid  bool
	Errors []error
}

// NewFailureCounterValidationResult creates a valid result.
func NewFailureCounterValidationResult() *FailureCounterValidationResult {
	return &FailureCounterValidationResult{Valid: true}
}

// AddError adds an error to the result and marks it invalid.
func (r *FailureCounterValidationResult) AddError(err error) {
	r.Valid = false
	r.Errors = append(r.Errors, err)
}

// ValidateFailureCounterUserID validates the user ID is not empty.
func ValidateFailureCounterUserID(userID string) error {
	if userID == "" {
		return ErrFailureCounterUserIDRequired
	}
	return nil
}

// ValidateFailureCounterLiftID validates the lift ID is not empty.
func ValidateFailureCounterLiftID(liftID string) error {
	if liftID == "" {
		return ErrFailureCounterLiftIDRequired
	}
	return nil
}

// ValidateFailureCounterProgressionID validates the progression ID is not empty.
func ValidateFailureCounterProgressionID(progressionID string) error {
	if progressionID == "" {
		return ErrFailureCounterProgressionIDRequired
	}
	return nil
}

// ValidateConsecutiveFailures validates consecutive failures is non-negative.
func ValidateConsecutiveFailures(consecutiveFailures int) error {
	if consecutiveFailures < 0 {
		return ErrConsecutiveFailuresNegative
	}
	return nil
}

// NewFailureCounter validates input and creates a new FailureCounter entity.
// Returns a validation result with errors if validation fails.
func NewFailureCounter(input CreateFailureCounterInput, id string) (*FailureCounter, *FailureCounterValidationResult) {
	result := NewFailureCounterValidationResult()

	if err := ValidateFailureCounterUserID(input.UserID); err != nil {
		result.AddError(err)
	}

	if err := ValidateFailureCounterLiftID(input.LiftID); err != nil {
		result.AddError(err)
	}

	if err := ValidateFailureCounterProgressionID(input.ProgressionID); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &FailureCounter{
		ID:                  id,
		UserID:              input.UserID,
		LiftID:              input.LiftID,
		ProgressionID:       input.ProgressionID,
		ConsecutiveFailures: 0,
		LastFailureAt:       nil,
		LastSuccessAt:       nil,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, result
}

// Validate performs full validation on an existing failure counter.
func (f *FailureCounter) Validate() *FailureCounterValidationResult {
	result := NewFailureCounterValidationResult()

	if err := ValidateFailureCounterUserID(f.UserID); err != nil {
		result.AddError(err)
	}

	if err := ValidateFailureCounterLiftID(f.LiftID); err != nil {
		result.AddError(err)
	}

	if err := ValidateFailureCounterProgressionID(f.ProgressionID); err != nil {
		result.AddError(err)
	}

	if err := ValidateConsecutiveFailures(f.ConsecutiveFailures); err != nil {
		result.AddError(err)
	}

	return result
}

// IncrementFailure increments the consecutive failure count and updates the last failure timestamp.
// Returns the new failure count.
func (f *FailureCounter) IncrementFailure() int {
	f.ConsecutiveFailures++
	now := time.Now()
	f.LastFailureAt = &now
	f.UpdatedAt = now
	return f.ConsecutiveFailures
}

// ResetOnSuccess resets the consecutive failure count to 0 and updates the last success timestamp.
func (f *FailureCounter) ResetOnSuccess() {
	f.ConsecutiveFailures = 0
	now := time.Now()
	f.LastSuccessAt = &now
	f.UpdatedAt = now
}

// HasFailures returns true if there are any consecutive failures recorded.
func (f *FailureCounter) HasFailures() bool {
	return f.ConsecutiveFailures > 0
}

// MeetsThreshold returns true if consecutive failures meet or exceed the given threshold.
func (f *FailureCounter) MeetsThreshold(threshold int) bool {
	return f.ConsecutiveFailures >= threshold
}
