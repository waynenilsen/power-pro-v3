// Package liftmax provides domain logic for the LiftMax entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package liftmax

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// MaxType represents the type of max (1RM or Training Max)
type MaxType string

const (
	OneRM       MaxType = "ONE_RM"
	TrainingMax MaxType = "TRAINING_MAX"
	E1RM        MaxType = "E1RM"
)

// Default TM percentage when converting between 1RM and TM
const DefaultTMPercentage = 90.0

// TM validation thresholds (as percentage of 1RM)
const (
	TMWarningLowerBound = 80.0  // TM should not be below 80% of 1RM
	TMWarningUpperBound = 95.0  // TM should not be above 95% of 1RM
)

// Validation errors
var (
	ErrValueRequired        = errors.New("lift max value is required")
	ErrValueNotPositive     = errors.New("lift max value must be positive")
	ErrValueInvalidPrecision = errors.New("lift max value must have precision of 0.25 (e.g., 315, 315.25, 315.5, 315.75)")
	ErrTypeRequired         = errors.New("max type is required")
	ErrTypeInvalid          = errors.New("max type must be ONE_RM, TRAINING_MAX, or E1RM")
	ErrEffectiveDateRequired = errors.New("effective date is required")
	ErrUserIDRequired       = errors.New("user ID is required")
	ErrLiftIDRequired       = errors.New("lift ID is required")
	ErrConversionPercentageInvalid = errors.New("conversion percentage must be between 1 and 100")
)

// LiftMax represents a lift max domain entity with all business rules.
type LiftMax struct {
	ID            string
	UserID        string
	LiftID        string
	Type          MaxType
	Value         float64
	EffectiveDate time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// LiftMaxRepository defines the interface for lift max persistence operations.
// This interface is used for TM validation warnings.
type LiftMaxRepository interface {
	// GetCurrentOneRM retrieves the most recent 1RM for a user and lift.
	// Returns nil if no 1RM exists.
	GetCurrentOneRM(userID, liftID string) (*LiftMax, error)
}

// ValidationResult is an alias for the shared validation.Result type.
// This includes support for warnings used during TM validation.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateValue validates the lift max value according to business rules.
// Value must be positive and have precision of 0.25.
func ValidateValue(value float64) error {
	if value <= 0 {
		return ErrValueNotPositive
	}

	// Check precision: value must be divisible by 0.25
	// Multiply by 4 and check if it's a whole number
	scaled := value * 4
	if math.Abs(scaled-math.Round(scaled)) > 0.0001 {
		return ErrValueInvalidPrecision
	}

	return nil
}

// ValidateType validates the max type according to business rules.
func ValidateType(maxType MaxType) error {
	if maxType == "" {
		return ErrTypeRequired
	}
	if maxType != OneRM && maxType != TrainingMax && maxType != E1RM {
		return ErrTypeInvalid
	}
	return nil
}

// ValidateUserID validates that user ID is provided.
func ValidateUserID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrUserIDRequired
	}
	return nil
}

// ValidateLiftID validates that lift ID is provided.
func ValidateLiftID(liftID string) error {
	if strings.TrimSpace(liftID) == "" {
		return ErrLiftIDRequired
	}
	return nil
}

// ValidateTMAgainstOneRM checks if a Training Max is within acceptable range of the 1RM.
// Returns a warning message if the TM is outside the 80%-95% range of the 1RM.
// This is a warning, not an error - the operation should proceed regardless.
func ValidateTMAgainstOneRM(tmValue float64, oneRM *LiftMax) string {
	if oneRM == nil {
		return ""
	}

	percentage := (tmValue / oneRM.Value) * 100

	if percentage < TMWarningLowerBound {
		return fmt.Sprintf("Training Max (%.2f) is below %.0f%% of 1RM (%.2f). Current: %.1f%%",
			tmValue, TMWarningLowerBound, oneRM.Value, percentage)
	}

	if percentage > TMWarningUpperBound {
		return fmt.Sprintf("Training Max (%.2f) is above %.0f%% of 1RM (%.2f). Current: %.1f%%",
			tmValue, TMWarningUpperBound, oneRM.Value, percentage)
	}

	return ""
}

// CreateLiftMaxInput contains the input data for creating a new lift max.
type CreateLiftMaxInput struct {
	UserID        string
	LiftID        string
	Type          MaxType
	Value         float64
	EffectiveDate *time.Time // Optional: defaults to current time
}

// CreateLiftMax validates input and creates a new LiftMax domain entity.
// When creating a Training Max, it checks against existing 1RM and generates warnings.
// Returns a validation result with errors if validation fails.
func CreateLiftMax(input CreateLiftMaxInput, id string, repo LiftMaxRepository) (*LiftMax, *ValidationResult) {
	result := NewValidationResult()

	// Validate user ID
	if err := ValidateUserID(input.UserID); err != nil {
		result.AddError(err)
	}

	// Validate lift ID
	if err := ValidateLiftID(input.LiftID); err != nil {
		result.AddError(err)
	}

	// Validate type
	if err := ValidateType(input.Type); err != nil {
		result.AddError(err)
	}

	// Validate value
	if err := ValidateValue(input.Value); err != nil {
		result.AddError(err)
	}

	// If basic validation failed, return early
	if !result.Valid {
		return nil, result
	}

	// For Training Max, check against existing 1RM and generate warning if needed
	if input.Type == TrainingMax && repo != nil {
		oneRM, err := repo.GetCurrentOneRM(input.UserID, input.LiftID)
		if err != nil {
			// Repository error - log but don't fail validation
			result.AddWarning(fmt.Sprintf("Unable to validate TM against 1RM: %v", err))
		} else if warning := ValidateTMAgainstOneRM(input.Value, oneRM); warning != "" {
			result.AddWarning(warning)
		}
	}

	// Set effective date
	effectiveDate := time.Now()
	if input.EffectiveDate != nil {
		effectiveDate = *input.EffectiveDate
	}

	now := time.Now()
	return &LiftMax{
		ID:            id,
		UserID:        input.UserID,
		LiftID:        input.LiftID,
		Type:          input.Type,
		Value:         input.Value,
		EffectiveDate: effectiveDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, result
}

// UpdateLiftMaxInput contains the input data for updating an existing lift max.
type UpdateLiftMaxInput struct {
	Type          *MaxType    // Optional: only update if provided
	Value         *float64    // Optional: only update if provided
	EffectiveDate *time.Time  // Optional: only update if provided
}

// UpdateLiftMax validates input and updates an existing LiftMax.
// Returns a validation result with errors if validation fails.
func UpdateLiftMax(liftMax *LiftMax, input UpdateLiftMaxInput, repo LiftMaxRepository) *ValidationResult {
	result := NewValidationResult()

	newType := liftMax.Type
	newValue := liftMax.Value

	// Validate type if provided
	if input.Type != nil {
		if err := ValidateType(*input.Type); err != nil {
			result.AddError(err)
		} else {
			newType = *input.Type
		}
	}

	// Validate value if provided
	if input.Value != nil {
		if err := ValidateValue(*input.Value); err != nil {
			result.AddError(err)
		} else {
			newValue = *input.Value
		}
	}

	// If basic validation failed, don't update the entity
	if !result.Valid {
		return result
	}

	// Apply changes
	if input.Type != nil {
		liftMax.Type = newType
	}
	if input.Value != nil {
		liftMax.Value = newValue
	}
	if input.EffectiveDate != nil {
		liftMax.EffectiveDate = *input.EffectiveDate
	}

	// For Training Max, check against existing 1RM and generate warning if needed
	if liftMax.Type == TrainingMax && repo != nil {
		oneRM, err := repo.GetCurrentOneRM(liftMax.UserID, liftMax.LiftID)
		if err != nil {
			result.AddWarning(fmt.Sprintf("Unable to validate TM against 1RM: %v", err))
		} else if warning := ValidateTMAgainstOneRM(liftMax.Value, oneRM); warning != "" {
			result.AddWarning(warning)
		}
	}

	liftMax.UpdatedAt = time.Now()

	return result
}

// Validate performs full validation on an existing lift max.
func (l *LiftMax) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateUserID(l.UserID); err != nil {
		result.AddError(err)
	}

	if err := ValidateLiftID(l.LiftID); err != nil {
		result.AddError(err)
	}

	if err := ValidateType(l.Type); err != nil {
		result.AddError(err)
	}

	if err := ValidateValue(l.Value); err != nil {
		result.AddError(err)
	}

	return result
}

// MaxCalculator provides conversion logic between 1RM and Training Max.
type MaxCalculator struct{}

// NewMaxCalculator creates a new MaxCalculator.
func NewMaxCalculator() *MaxCalculator {
	return &MaxCalculator{}
}

// ConvertToTM converts a 1RM value to Training Max.
// Uses the specified percentage or defaults to 90%.
// Result is rounded to nearest 0.25.
func (c *MaxCalculator) ConvertToTM(oneRM float64, percentage *float64) (float64, error) {
	pct := DefaultTMPercentage
	if percentage != nil {
		pct = *percentage
	}

	if pct <= 0 || pct > 100 {
		return 0, ErrConversionPercentageInvalid
	}

	tm := oneRM * (pct / 100)
	return RoundToQuarter(tm), nil
}

// ConvertToOneRM converts a Training Max value to estimated 1RM.
// Uses the specified percentage or defaults to 90%.
// Result is rounded to nearest 0.25.
func (c *MaxCalculator) ConvertToOneRM(tm float64, percentage *float64) (float64, error) {
	pct := DefaultTMPercentage
	if percentage != nil {
		pct = *percentage
	}

	if pct <= 0 || pct > 100 {
		return 0, ErrConversionPercentageInvalid
	}

	oneRM := tm / (pct / 100)
	return RoundToQuarter(oneRM), nil
}

// RoundToQuarter rounds a value to the nearest 0.25.
func RoundToQuarter(value float64) float64 {
	return math.Round(value*4) / 4
}

// ConversionResult contains the result of a max conversion.
type ConversionResult struct {
	OriginalValue float64
	OriginalType  MaxType
	ConvertedValue float64
	ConvertedType  MaxType
	Percentage     float64
}

// Convert performs a conversion and returns a detailed result.
func (c *MaxCalculator) Convert(value float64, fromType MaxType, percentage *float64) (*ConversionResult, error) {
	pct := DefaultTMPercentage
	if percentage != nil {
		pct = *percentage
	}

	var convertedValue float64
	var toType MaxType
	var err error

	switch fromType {
	case OneRM:
		convertedValue, err = c.ConvertToTM(value, &pct)
		toType = TrainingMax
	case TrainingMax:
		convertedValue, err = c.ConvertToOneRM(value, &pct)
		toType = OneRM
	default:
		return nil, ErrTypeInvalid
	}

	if err != nil {
		return nil, err
	}

	return &ConversionResult{
		OriginalValue:  value,
		OriginalType:   fromType,
		ConvertedValue: convertedValue,
		ConvertedType:  toType,
		Percentage:     pct,
	}, nil
}
