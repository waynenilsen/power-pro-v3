// Package cycle provides domain logic for the Cycle entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package cycle

import (
	"errors"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// Validation errors
var (
	ErrNameRequired       = errors.New("name is required")
	ErrNameTooLong        = errors.New("name must be at most 100 characters")
	ErrLengthWeeksInvalid = errors.New("length_weeks must be >= 1")
)

// MaxNameLength is the maximum length for a cycle name.
const MaxNameLength = 100

// Cycle represents a cycle domain entity with all business rules.
type Cycle struct {
	ID          string
	Name        string
	LengthWeeks int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CycleWeek represents a week within a cycle (for response data).
type CycleWeek struct {
	ID         string
	WeekNumber int
}

// CycleWithWeeks represents a cycle with its associated weeks.
type CycleWithWeeks struct {
	Cycle *Cycle
	Weeks []CycleWeek
}

// ValidationResult is an alias for the shared validation.Result type.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateName validates the cycle name according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrNameRequired
	}
	if len(name) > MaxNameLength {
		return ErrNameTooLong
	}
	return nil
}

// ValidateLengthWeeks validates the length_weeks field.
// Returns an error if validation fails, nil otherwise.
func ValidateLengthWeeks(lengthWeeks int) error {
	if lengthWeeks < 1 {
		return ErrLengthWeeksInvalid
	}
	return nil
}

// CreateCycleInput contains the input data for creating a new cycle.
type CreateCycleInput struct {
	Name        string
	LengthWeeks int
}

// CreateCycle validates input and creates a new Cycle domain entity.
// Returns a validation result with errors if validation fails.
func CreateCycle(input CreateCycleInput, id string) (*Cycle, *ValidationResult) {
	result := NewValidationResult()

	// Validate name
	if err := ValidateName(input.Name); err != nil {
		result.AddError(err)
	}

	// Validate length_weeks
	if err := ValidateLengthWeeks(input.LengthWeeks); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &Cycle{
		ID:          id,
		Name:        strings.TrimSpace(input.Name),
		LengthWeeks: input.LengthWeeks,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, result
}

// UpdateCycleInput contains the input data for updating an existing cycle.
type UpdateCycleInput struct {
	Name        *string // Optional: only update if provided
	LengthWeeks *int    // Optional: only update if provided
}

// UpdateCycle validates input and updates an existing Cycle.
// Returns a validation result with errors if validation fails.
func UpdateCycle(c *Cycle, input UpdateCycleInput) *ValidationResult {
	result := NewValidationResult()

	// Validate name if provided
	if input.Name != nil {
		if err := ValidateName(*input.Name); err != nil {
			result.AddError(err)
		} else {
			c.Name = strings.TrimSpace(*input.Name)
		}
	}

	// Validate length_weeks if provided
	if input.LengthWeeks != nil {
		if err := ValidateLengthWeeks(*input.LengthWeeks); err != nil {
			result.AddError(err)
		} else {
			c.LengthWeeks = *input.LengthWeeks
		}
	}

	if result.Valid {
		c.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing cycle.
func (c *Cycle) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateName(c.Name); err != nil {
		result.AddError(err)
	}

	if err := ValidateLengthWeeks(c.LengthWeeks); err != nil {
		result.AddError(err)
	}

	return result
}
