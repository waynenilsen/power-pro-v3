// Package loggedset provides domain logic for the LoggedSet entity.
// LoggedSet records actual performance data for sets, enabling AMRAP-driven progressions.
package loggedset

import (
	"errors"
	"time"
)

// Validation errors
var (
	ErrUserIDRequired         = errors.New("user_id is required")
	ErrSessionIDRequired      = errors.New("session_id is required")
	ErrPrescriptionIDRequired = errors.New("prescription_id is required")
	ErrLiftIDRequired         = errors.New("lift_id is required")
	ErrSetNumberInvalid       = errors.New("set_number must be positive")
	ErrWeightInvalid          = errors.New("weight must be non-negative")
	ErrTargetRepsInvalid      = errors.New("target_reps must be positive")
	ErrRepsPerformedInvalid   = errors.New("reps_performed must be non-negative")
	ErrRPEInvalid             = errors.New("rpe must be between 5.0 and 10.0")
)

// LoggedSet represents a single logged set from a workout session.
type LoggedSet struct {
	ID             string
	UserID         string
	SessionID      string
	PrescriptionID string
	LiftID         string
	SetNumber      int
	Weight         float64
	TargetReps     int
	RepsPerformed  int
	IsAMRAP        bool
	// RPE is the rate of perceived exertion (5.0-10.0).
	// Optional - nil means RPE was not recorded.
	RPE       *float64
	CreatedAt time.Time
}

// CreateLoggedSetInput contains the input data for creating a new logged set.
type CreateLoggedSetInput struct {
	UserID         string
	SessionID      string
	PrescriptionID string
	LiftID         string
	SetNumber      int
	Weight         float64
	TargetReps     int
	RepsPerformed  int
	IsAMRAP        bool
	RPE            *float64
}

// ValidationResult holds validation errors.
type ValidationResult struct {
	Valid  bool
	Errors []error
}

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return &ValidationResult{Valid: true}
}

// AddError adds an error to the result and marks it invalid.
func (r *ValidationResult) AddError(err error) {
	r.Valid = false
	r.Errors = append(r.Errors, err)
}

// ValidateUserID validates the user ID is not empty.
func ValidateUserID(userID string) error {
	if userID == "" {
		return ErrUserIDRequired
	}
	return nil
}

// ValidateSessionID validates the session ID is not empty.
func ValidateSessionID(sessionID string) error {
	if sessionID == "" {
		return ErrSessionIDRequired
	}
	return nil
}

// ValidatePrescriptionID validates the prescription ID is not empty.
func ValidatePrescriptionID(prescriptionID string) error {
	if prescriptionID == "" {
		return ErrPrescriptionIDRequired
	}
	return nil
}

// ValidateLiftID validates the lift ID is not empty.
func ValidateLiftID(liftID string) error {
	if liftID == "" {
		return ErrLiftIDRequired
	}
	return nil
}

// ValidateSetNumber validates the set number is positive.
func ValidateSetNumber(setNumber int) error {
	if setNumber < 1 {
		return ErrSetNumberInvalid
	}
	return nil
}

// ValidateWeight validates the weight is non-negative.
func ValidateWeight(weight float64) error {
	if weight < 0 {
		return ErrWeightInvalid
	}
	return nil
}

// ValidateTargetReps validates the target reps is positive.
func ValidateTargetReps(targetReps int) error {
	if targetReps < 1 {
		return ErrTargetRepsInvalid
	}
	return nil
}

// ValidateRepsPerformed validates the reps performed is non-negative.
func ValidateRepsPerformed(repsPerformed int) error {
	if repsPerformed < 0 {
		return ErrRepsPerformedInvalid
	}
	return nil
}

// ValidateRPE validates that RPE is within the valid range (5.0-10.0) if provided.
func ValidateRPE(rpe *float64) error {
	if rpe != nil {
		if *rpe < 5.0 || *rpe > 10.0 {
			return ErrRPEInvalid
		}
	}
	return nil
}

// NewLoggedSet validates input and creates a new LoggedSet domain entity.
// Returns a validation result with errors if validation fails.
func NewLoggedSet(input CreateLoggedSetInput, id string) (*LoggedSet, *ValidationResult) {
	result := NewValidationResult()

	if err := ValidateUserID(input.UserID); err != nil {
		result.AddError(err)
	}

	if err := ValidateSessionID(input.SessionID); err != nil {
		result.AddError(err)
	}

	if err := ValidatePrescriptionID(input.PrescriptionID); err != nil {
		result.AddError(err)
	}

	if err := ValidateLiftID(input.LiftID); err != nil {
		result.AddError(err)
	}

	if err := ValidateSetNumber(input.SetNumber); err != nil {
		result.AddError(err)
	}

	if err := ValidateWeight(input.Weight); err != nil {
		result.AddError(err)
	}

	if err := ValidateTargetReps(input.TargetReps); err != nil {
		result.AddError(err)
	}

	if err := ValidateRepsPerformed(input.RepsPerformed); err != nil {
		result.AddError(err)
	}

	if err := ValidateRPE(input.RPE); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	return &LoggedSet{
		ID:             id,
		UserID:         input.UserID,
		SessionID:      input.SessionID,
		PrescriptionID: input.PrescriptionID,
		LiftID:         input.LiftID,
		SetNumber:      input.SetNumber,
		Weight:         input.Weight,
		TargetReps:     input.TargetReps,
		RepsPerformed:  input.RepsPerformed,
		IsAMRAP:        input.IsAMRAP,
		RPE:            input.RPE,
		CreatedAt:      time.Now(),
	}, result
}

// Validate performs full validation on an existing logged set.
func (l *LoggedSet) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateUserID(l.UserID); err != nil {
		result.AddError(err)
	}

	if err := ValidateSessionID(l.SessionID); err != nil {
		result.AddError(err)
	}

	if err := ValidatePrescriptionID(l.PrescriptionID); err != nil {
		result.AddError(err)
	}

	if err := ValidateLiftID(l.LiftID); err != nil {
		result.AddError(err)
	}

	if err := ValidateSetNumber(l.SetNumber); err != nil {
		result.AddError(err)
	}

	if err := ValidateWeight(l.Weight); err != nil {
		result.AddError(err)
	}

	if err := ValidateTargetReps(l.TargetReps); err != nil {
		result.AddError(err)
	}

	if err := ValidateRepsPerformed(l.RepsPerformed); err != nil {
		result.AddError(err)
	}

	if err := ValidateRPE(l.RPE); err != nil {
		result.AddError(err)
	}

	return result
}

// ExceededTarget returns true if reps performed exceeded target reps.
func (l *LoggedSet) ExceededTarget() bool {
	return l.RepsPerformed > l.TargetReps
}

// RepsDifference returns the difference between reps performed and target.
// Positive means exceeded target, negative means fell short.
func (l *LoggedSet) RepsDifference() int {
	return l.RepsPerformed - l.TargetReps
}
