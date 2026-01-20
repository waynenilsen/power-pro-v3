// Package weeklylookup provides domain logic for the WeeklyLookup entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package weeklylookup

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Validation errors
var (
	ErrNameRequired          = errors.New("name is required")
	ErrNameTooLong           = errors.New("name must be at most 100 characters")
	ErrEntriesRequired       = errors.New("entries are required")
	ErrWeekNumberInvalid     = errors.New("week number must be >= 1")
	ErrPercentagesRequired   = errors.New("percentages are required for each entry")
	ErrRepsRequired          = errors.New("reps are required for each entry")
	ErrPercentageRepsLengthMismatch = errors.New("percentages and reps must have the same length")
	ErrDuplicateWeekNumber   = errors.New("duplicate week number in entries")
)

// MaxNameLength is the maximum length for a weekly lookup name.
const MaxNameLength = 100

// WeeklyLookup represents a weekly lookup domain entity with all business rules.
type WeeklyLookup struct {
	ID        string
	Name      string
	Entries   []WeeklyLookupEntry
	ProgramID *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WeeklyLookupEntry represents an entry in a weekly lookup table.
// It maps a week number to percentages, reps, and an optional percentage modifier.
type WeeklyLookupEntry struct {
	WeekNumber         int       `json:"weekNumber"`
	Percentages        []float64 `json:"percentages"`
	Reps               []int     `json:"reps"`
	PercentageModifier *float64  `json:"percentageModifier,omitempty"`
}

// ValidationResult contains the result of validating a weekly lookup.
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

// ValidateName validates the weekly lookup name according to business rules.
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

// ValidateEntries validates the weekly lookup entries.
// Returns an error if validation fails, nil otherwise.
func ValidateEntries(entries []WeeklyLookupEntry) error {
	if len(entries) == 0 {
		return ErrEntriesRequired
	}

	weekNumbers := make(map[int]bool)
	for _, entry := range entries {
		if entry.WeekNumber < 1 {
			return ErrWeekNumberInvalid
		}
		if weekNumbers[entry.WeekNumber] {
			return ErrDuplicateWeekNumber
		}
		weekNumbers[entry.WeekNumber] = true

		if len(entry.Percentages) == 0 {
			return ErrPercentagesRequired
		}
		if len(entry.Reps) == 0 {
			return ErrRepsRequired
		}
		if len(entry.Percentages) != len(entry.Reps) {
			return ErrPercentageRepsLengthMismatch
		}
	}
	return nil
}

// CreateWeeklyLookupInput contains the input data for creating a new weekly lookup.
type CreateWeeklyLookupInput struct {
	Name      string
	Entries   []WeeklyLookupEntry
	ProgramID *string
}

// CreateWeeklyLookup validates input and creates a new WeeklyLookup domain entity.
// Returns a validation result with errors if validation fails.
func CreateWeeklyLookup(input CreateWeeklyLookupInput, id string) (*WeeklyLookup, *ValidationResult) {
	result := NewValidationResult()

	// Validate name
	if err := ValidateName(input.Name); err != nil {
		result.AddError(err)
	}

	// Validate entries
	if err := ValidateEntries(input.Entries); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &WeeklyLookup{
		ID:        id,
		Name:      strings.TrimSpace(input.Name),
		Entries:   input.Entries,
		ProgramID: input.ProgramID,
		CreatedAt: now,
		UpdatedAt: now,
	}, result
}

// UpdateWeeklyLookupInput contains the input data for updating an existing weekly lookup.
type UpdateWeeklyLookupInput struct {
	Name      *string              // Optional: only update if provided
	Entries   *[]WeeklyLookupEntry // Optional: only update if provided
	ProgramID **string             // Optional: only update if provided (double pointer to distinguish nil from clearing)
}

// UpdateWeeklyLookup validates input and updates an existing WeeklyLookup.
// Returns a validation result with errors if validation fails.
func UpdateWeeklyLookup(w *WeeklyLookup, input UpdateWeeklyLookupInput) *ValidationResult {
	result := NewValidationResult()

	// Validate name if provided
	if input.Name != nil {
		if err := ValidateName(*input.Name); err != nil {
			result.AddError(err)
		} else {
			w.Name = strings.TrimSpace(*input.Name)
		}
	}

	// Validate entries if provided
	if input.Entries != nil {
		if err := ValidateEntries(*input.Entries); err != nil {
			result.AddError(err)
		} else {
			w.Entries = *input.Entries
		}
	}

	// Update program_id if provided
	if input.ProgramID != nil {
		w.ProgramID = *input.ProgramID
	}

	if result.Valid {
		w.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing weekly lookup.
func (w *WeeklyLookup) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateName(w.Name); err != nil {
		result.AddError(err)
	}

	if err := ValidateEntries(w.Entries); err != nil {
		result.AddError(err)
	}

	return result
}

// GetByWeekNumber returns the entry for a specific week number, or nil if not found.
func (w *WeeklyLookup) GetByWeekNumber(weekNum int) *WeeklyLookupEntry {
	for i := range w.Entries {
		if w.Entries[i].WeekNumber == weekNum {
			return &w.Entries[i]
		}
	}
	return nil
}
