// Package program provides domain logic for the Program entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package program

import (
	"errors"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// Validation errors
var (
	ErrNameRequired           = errors.New("name is required")
	ErrNameTooLong            = errors.New("name must be at most 100 characters")
	ErrSlugRequired           = errors.New("slug is required")
	ErrSlugTooLong            = errors.New("slug must be at most 100 characters")
	ErrSlugInvalid            = errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	ErrCycleIDRequired        = errors.New("cycle_id is required")
	ErrDefaultRoundingInvalid = errors.New("default_rounding must be positive if provided")
	ErrInvalidDifficulty      = errors.New("difficulty must be one of: beginner, intermediate, advanced")
	ErrInvalidDaysPerWeek     = errors.New("days_per_week must be between 1 and 7")
	ErrInvalidFocus           = errors.New("focus must be one of: strength, hypertrophy, peaking")
)

// Valid values for filter fields
var (
	ValidDifficulties = []string{"beginner", "intermediate", "advanced"}
	ValidFocusValues  = []string{"strength", "hypertrophy", "peaking"}
)

// MaxNameLength is the maximum length for a program name.
const MaxNameLength = 100

// MaxSlugLength is the maximum length for a program slug.
const MaxSlugLength = 100

// Program represents a program domain entity with all business rules.
type Program struct {
	ID              string
	Name            string
	Slug            string
	Description     *string
	CycleID         string
	WeeklyLookupID  *string
	DailyLookupID   *string
	DefaultRounding *float64
	Difficulty      string
	DaysPerWeek     int
	Focus           string
	HasAmrap        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ProgramWithCycle represents a program with its associated cycle data (for detailed responses).
type ProgramWithCycle struct {
	Program *Program
	Cycle   *ProgramCycle
}

// ProgramCycle represents embedded cycle info in a program response.
type ProgramCycle struct {
	ID          string
	Name        string
	LengthWeeks int
	Weeks       []ProgramCycleWeek
}

// ProgramCycleWeek represents a week reference within a program's cycle.
type ProgramCycleWeek struct {
	ID         string
	WeekNumber int
}

// LookupReference represents a lookup table reference (not full entries).
type LookupReference struct {
	ID   string
	Name string
}

// FilterOptions represents optional filters for listing programs.
type FilterOptions struct {
	Difficulty  *string
	DaysPerWeek *int
	Focus       *string
	HasAmrap    *bool
}

// ValidateDifficulty validates the difficulty value.
func ValidateDifficulty(difficulty string) error {
	for _, valid := range ValidDifficulties {
		if difficulty == valid {
			return nil
		}
	}
	return ErrInvalidDifficulty
}

// ValidateDaysPerWeek validates the days_per_week value.
func ValidateDaysPerWeek(days int) error {
	if days < 1 || days > 7 {
		return ErrInvalidDaysPerWeek
	}
	return nil
}

// ValidateFocus validates the focus value.
func ValidateFocus(focus string) error {
	for _, valid := range ValidFocusValues {
		if focus == valid {
			return nil
		}
	}
	return ErrInvalidFocus
}

// Validate validates all filter options and returns a validation result.
func (f *FilterOptions) Validate() *ValidationResult {
	result := NewValidationResult()

	if f.Difficulty != nil {
		if err := ValidateDifficulty(*f.Difficulty); err != nil {
			result.AddError(err)
		}
	}

	if f.DaysPerWeek != nil {
		if err := ValidateDaysPerWeek(*f.DaysPerWeek); err != nil {
			result.AddError(err)
		}
	}

	if f.Focus != nil {
		if err := ValidateFocus(*f.Focus); err != nil {
			result.AddError(err)
		}
	}

	// HasAmrap is boolean, no validation needed (true/false always valid)

	return result
}

// ValidationResult is an alias for the shared validation.Result type.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateName validates the program name according to business rules.
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

// ValidateSlug validates the program slug according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateSlug(slug string) error {
	if strings.TrimSpace(slug) == "" {
		return ErrSlugRequired
	}
	err := validation.ValidateSlug(slug, MaxSlugLength)
	if err == nil {
		return nil
	}
	// Map shared validation errors to package-specific errors for backward compatibility
	if errors.Is(err, validation.ErrSlugInvalid) {
		return ErrSlugInvalid
	}
	// Check for slug too long error by checking error message prefix
	if strings.HasPrefix(err.Error(), "slug must be") {
		return ErrSlugTooLong
	}
	return ErrSlugInvalid
}

// ValidateCycleID validates the cycle_id field.
// Returns an error if validation fails, nil otherwise.
func ValidateCycleID(cycleID string) error {
	if strings.TrimSpace(cycleID) == "" {
		return ErrCycleIDRequired
	}
	return nil
}

// ValidateDefaultRounding validates the default_rounding field.
// Returns an error if validation fails, nil otherwise.
func ValidateDefaultRounding(rounding *float64) error {
	if rounding != nil && *rounding <= 0 {
		return ErrDefaultRoundingInvalid
	}
	return nil
}

// CreateProgramInput contains the input data for creating a new program.
type CreateProgramInput struct {
	Name            string
	Slug            string
	Description     *string
	CycleID         string
	WeeklyLookupID  *string
	DailyLookupID   *string
	DefaultRounding *float64
}

// CreateProgram validates input and creates a new Program domain entity.
// Returns a validation result with errors if validation fails.
func CreateProgram(input CreateProgramInput, id string) (*Program, *ValidationResult) {
	result := NewValidationResult()

	// Validate name
	if err := ValidateName(input.Name); err != nil {
		result.AddError(err)
	}

	// Validate slug
	if err := ValidateSlug(input.Slug); err != nil {
		result.AddError(err)
	}

	// Validate cycle_id
	if err := ValidateCycleID(input.CycleID); err != nil {
		result.AddError(err)
	}

	// Validate default_rounding
	if err := ValidateDefaultRounding(input.DefaultRounding); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &Program{
		ID:              id,
		Name:            strings.TrimSpace(input.Name),
		Slug:            strings.TrimSpace(input.Slug),
		Description:     trimStringPtr(input.Description),
		CycleID:         strings.TrimSpace(input.CycleID),
		WeeklyLookupID:  trimStringPtr(input.WeeklyLookupID),
		DailyLookupID:   trimStringPtr(input.DailyLookupID),
		DefaultRounding: input.DefaultRounding,
		Difficulty:      "beginner",  // Default per schema
		DaysPerWeek:     3,           // Default per schema
		Focus:           "strength",  // Default per schema
		HasAmrap:        false,       // Default per schema
		CreatedAt:       now,
		UpdatedAt:       now,
	}, result
}

// UpdateProgramInput contains the input data for updating an existing program.
type UpdateProgramInput struct {
	Name            *string  // Optional: only update if provided
	Slug            *string  // Optional: only update if provided
	Description     **string // Double pointer: nil = no change, *nil = clear, *value = set
	CycleID         *string  // Optional: only update if provided
	WeeklyLookupID  **string // Double pointer: nil = no change, *nil = clear, *value = set
	DailyLookupID   **string // Double pointer: nil = no change, *nil = clear, *value = set
	DefaultRounding **float64 // Double pointer: nil = no change, *nil = clear, *value = set
}

// UpdateProgram validates input and updates an existing Program.
// Returns a validation result with errors if validation fails.
func UpdateProgram(p *Program, input UpdateProgramInput) *ValidationResult {
	result := NewValidationResult()

	// Validate name if provided
	if input.Name != nil {
		if err := ValidateName(*input.Name); err != nil {
			result.AddError(err)
		} else {
			p.Name = strings.TrimSpace(*input.Name)
		}
	}

	// Validate slug if provided
	if input.Slug != nil {
		if err := ValidateSlug(*input.Slug); err != nil {
			result.AddError(err)
		} else {
			p.Slug = strings.TrimSpace(*input.Slug)
		}
	}

	// Validate cycle_id if provided
	if input.CycleID != nil {
		if err := ValidateCycleID(*input.CycleID); err != nil {
			result.AddError(err)
		} else {
			p.CycleID = strings.TrimSpace(*input.CycleID)
		}
	}

	// Handle description (double pointer for nullable field)
	if input.Description != nil {
		p.Description = trimStringPtr(*input.Description)
	}

	// Handle weekly_lookup_id (double pointer for nullable field)
	if input.WeeklyLookupID != nil {
		p.WeeklyLookupID = trimStringPtr(*input.WeeklyLookupID)
	}

	// Handle daily_lookup_id (double pointer for nullable field)
	if input.DailyLookupID != nil {
		p.DailyLookupID = trimStringPtr(*input.DailyLookupID)
	}

	// Handle default_rounding (double pointer for nullable field)
	if input.DefaultRounding != nil {
		newRounding := *input.DefaultRounding
		if err := ValidateDefaultRounding(newRounding); err != nil {
			result.AddError(err)
		} else {
			p.DefaultRounding = newRounding
		}
	}

	if result.Valid {
		p.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing program.
func (p *Program) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateName(p.Name); err != nil {
		result.AddError(err)
	}

	if err := ValidateSlug(p.Slug); err != nil {
		result.AddError(err)
	}

	if err := ValidateCycleID(p.CycleID); err != nil {
		result.AddError(err)
	}

	if err := ValidateDefaultRounding(p.DefaultRounding); err != nil {
		result.AddError(err)
	}

	return result
}

// trimStringPtr trims whitespace from a string pointer.
func trimStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
