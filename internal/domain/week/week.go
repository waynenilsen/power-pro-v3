// Package week provides domain logic for the Week entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package week

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Validation errors
var (
	ErrWeekNumberRequired     = errors.New("week number is required")
	ErrWeekNumberInvalid      = errors.New("week number must be >= 1")
	ErrCycleIDRequired        = errors.New("cycle_id is required")
	ErrVariantInvalid         = errors.New("variant must be 'A', 'B', or empty")
	ErrDayIDRequired          = errors.New("day_id is required")
	ErrDayOfWeekRequired      = errors.New("day_of_week is required")
	ErrDayOfWeekInvalid       = errors.New("day_of_week must be MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, or SUNDAY")
	ErrWeekDayNotFound        = errors.New("week day not found in week")
	ErrDuplicateDayInWeek     = errors.New("duplicate day in week")
)

// DayOfWeek represents valid days of the week.
type DayOfWeek string

const (
	Monday    DayOfWeek = "MONDAY"
	Tuesday   DayOfWeek = "TUESDAY"
	Wednesday DayOfWeek = "WEDNESDAY"
	Thursday  DayOfWeek = "THURSDAY"
	Friday    DayOfWeek = "FRIDAY"
	Saturday  DayOfWeek = "SATURDAY"
	Sunday    DayOfWeek = "SUNDAY"
)

// ValidDaysOfWeek contains all valid days of the week.
var ValidDaysOfWeek = []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}

// Variant represents valid week variants.
type Variant string

const (
	VariantA    Variant = "A"
	VariantB    Variant = "B"
	VariantNone Variant = ""
)

// Week represents a week domain entity with all business rules.
type Week struct {
	ID         string
	WeekNumber int
	Variant    *string
	CycleID    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// WeekDay represents a day associated with a week.
type WeekDay struct {
	ID        string
	WeekID    string
	DayID     string
	DayOfWeek DayOfWeek
	CreatedAt time.Time
}

// WeekWithDays represents a week with its associated days.
type WeekWithDays struct {
	Week *Week
	Days []WeekDay
}

// ValidationResult contains the result of validating a week.
type ValidationResult struct {
	Valid  bool
	Errors []error
}

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return &ValidationResult{Valid: true, Errors: []error{}}
}

// AddError adds an error to the validation result and marks it invalid.
func (v *ValidationResult) AddError(err error) {
	v.Valid = false
	v.Errors = append(v.Errors, err)
}

// Error returns a combined error message if there are validation errors.
func (v *ValidationResult) Error() error {
	if v.Valid {
		return nil
	}
	var msgs []string
	for _, err := range v.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Errorf("validation failed: %s", strings.Join(msgs, "; "))
}

// ValidateWeekNumber validates the week number according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateWeekNumber(weekNumber int) error {
	if weekNumber < 1 {
		return ErrWeekNumberInvalid
	}
	return nil
}

// ValidateCycleID validates the cycle ID according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateCycleID(cycleID string) error {
	if strings.TrimSpace(cycleID) == "" {
		return ErrCycleIDRequired
	}
	return nil
}

// ValidateVariant validates the variant field.
// Returns an error if validation fails, nil otherwise.
func ValidateVariant(variant *string) error {
	if variant == nil {
		return nil
	}
	v := *variant
	if v != string(VariantA) && v != string(VariantB) && v != string(VariantNone) {
		return ErrVariantInvalid
	}
	return nil
}

// ValidateDayOfWeek validates the day of week field.
// Returns an error if validation fails, nil otherwise.
func ValidateDayOfWeek(dayOfWeek string) error {
	if dayOfWeek == "" {
		return ErrDayOfWeekRequired
	}
	for _, valid := range ValidDaysOfWeek {
		if string(valid) == dayOfWeek {
			return nil
		}
	}
	return ErrDayOfWeekInvalid
}

// CreateWeekInput contains the input data for creating a new week.
type CreateWeekInput struct {
	WeekNumber int
	Variant    *string
	CycleID    string
}

// CreateWeek validates input and creates a new Week domain entity.
// Returns a validation result with errors if validation fails.
func CreateWeek(input CreateWeekInput, id string) (*Week, *ValidationResult) {
	result := NewValidationResult()

	// Validate week number
	if err := ValidateWeekNumber(input.WeekNumber); err != nil {
		result.AddError(err)
	}

	// Validate cycle ID
	if err := ValidateCycleID(input.CycleID); err != nil {
		result.AddError(err)
	}

	// Validate variant
	if err := ValidateVariant(input.Variant); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	// Normalize empty string variant to nil
	var variant *string
	if input.Variant != nil && *input.Variant != "" {
		variant = input.Variant
	}

	now := time.Now()
	return &Week{
		ID:         id,
		WeekNumber: input.WeekNumber,
		Variant:    variant,
		CycleID:    input.CycleID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, result
}

// UpdateWeekInput contains the input data for updating an existing week.
type UpdateWeekInput struct {
	WeekNumber   *int    // Optional: only update if provided
	Variant      *string // Optional: only update if provided
	ClearVariant bool    // Set to true to explicitly clear variant
	CycleID      *string // Optional: only update if provided
}

// UpdateWeek validates input and updates an existing Week.
// Returns a validation result with errors if validation fails.
func UpdateWeek(w *Week, input UpdateWeekInput) *ValidationResult {
	result := NewValidationResult()

	// Validate week number if provided
	if input.WeekNumber != nil {
		if err := ValidateWeekNumber(*input.WeekNumber); err != nil {
			result.AddError(err)
		} else {
			w.WeekNumber = *input.WeekNumber
		}
	}

	// Handle variant update
	if input.ClearVariant {
		w.Variant = nil
	} else if input.Variant != nil {
		if err := ValidateVariant(input.Variant); err != nil {
			result.AddError(err)
		} else {
			// Normalize empty string variant to nil
			if *input.Variant == "" {
				w.Variant = nil
			} else {
				w.Variant = input.Variant
			}
		}
	}

	// Validate cycle ID if provided
	if input.CycleID != nil {
		if err := ValidateCycleID(*input.CycleID); err != nil {
			result.AddError(err)
		} else {
			w.CycleID = *input.CycleID
		}
	}

	if result.Valid {
		w.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing week.
func (w *Week) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateWeekNumber(w.WeekNumber); err != nil {
		result.AddError(err)
	}

	if err := ValidateCycleID(w.CycleID); err != nil {
		result.AddError(err)
	}

	if err := ValidateVariant(w.Variant); err != nil {
		result.AddError(err)
	}

	return result
}

// CreateWeekDayInput contains the input data for adding a day to a week.
type CreateWeekDayInput struct {
	WeekID    string
	DayID     string
	DayOfWeek string
}

// CreateWeekDay validates input and creates a new WeekDay.
// Returns a validation result with errors if validation fails.
func CreateWeekDay(input CreateWeekDayInput, id string) (*WeekDay, *ValidationResult) {
	result := NewValidationResult()

	if input.WeekID == "" {
		result.AddError(errors.New("week ID is required"))
	}

	if input.DayID == "" {
		result.AddError(ErrDayIDRequired)
	}

	if err := ValidateDayOfWeek(input.DayOfWeek); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	return &WeekDay{
		ID:        id,
		WeekID:    input.WeekID,
		DayID:     input.DayID,
		DayOfWeek: DayOfWeek(input.DayOfWeek),
		CreatedAt: time.Now(),
	}, result
}

// DayOfWeekOrder returns the sort order of a day of week (Monday=1, Sunday=7).
func DayOfWeekOrder(dow DayOfWeek) int {
	switch dow {
	case Monday:
		return 1
	case Tuesday:
		return 2
	case Wednesday:
		return 3
	case Thursday:
		return 4
	case Friday:
		return 5
	case Saturday:
		return 6
	case Sunday:
		return 7
	default:
		return 8
	}
}
