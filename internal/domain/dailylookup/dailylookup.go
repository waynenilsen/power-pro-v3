// Package dailylookup provides domain logic for the DailyLookup entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package dailylookup

import (
	"errors"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// Validation errors
var (
	ErrNameRequired             = errors.New("name is required")
	ErrNameTooLong              = errors.New("name must be at most 100 characters")
	ErrEntriesRequired          = errors.New("entries are required")
	ErrDayIdentifierRequired    = errors.New("day identifier is required")
	ErrDuplicateDayIdentifier   = errors.New("duplicate day identifier in entries")
	ErrIntensityLevelInvalid    = errors.New("intensity level must be one of: HEAVY, LIGHT, MEDIUM")
)

// MaxNameLength is the maximum length for a daily lookup name.
const MaxNameLength = 100

// IntensityLevel represents valid intensity levels.
type IntensityLevel string

const (
	IntensityHeavy  IntensityLevel = "HEAVY"
	IntensityLight  IntensityLevel = "LIGHT"
	IntensityMedium IntensityLevel = "MEDIUM"
)

// ValidIntensityLevels contains all valid intensity levels.
var ValidIntensityLevels = []IntensityLevel{IntensityHeavy, IntensityLight, IntensityMedium}

// IsValidIntensityLevel checks if a string is a valid intensity level.
func IsValidIntensityLevel(level string) bool {
	for _, valid := range ValidIntensityLevels {
		if IntensityLevel(strings.ToUpper(level)) == valid {
			return true
		}
	}
	return false
}

// DailyLookup represents a daily lookup domain entity with all business rules.
type DailyLookup struct {
	ID        string
	Name      string
	Entries   []DailyLookupEntry
	ProgramID *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DailyLookupEntry represents an entry in a daily lookup table.
// It maps a day identifier (e.g., "heavy", "light", "medium", or day slugs) to parameters.
type DailyLookupEntry struct {
	DayIdentifier      string   `json:"dayIdentifier"`
	PercentageModifier float64  `json:"percentageModifier"`
	IntensityLevel     *string  `json:"intensityLevel,omitempty"`
}

// ValidationResult is an alias for the shared validation.Result type.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateName validates the daily lookup name according to business rules.
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

// ValidateEntries validates the daily lookup entries.
// Returns an error if validation fails, nil otherwise.
func ValidateEntries(entries []DailyLookupEntry) error {
	if len(entries) == 0 {
		return ErrEntriesRequired
	}

	dayIdentifiers := make(map[string]bool)
	for _, entry := range entries {
		if strings.TrimSpace(entry.DayIdentifier) == "" {
			return ErrDayIdentifierRequired
		}
		normalizedID := strings.ToLower(entry.DayIdentifier)
		if dayIdentifiers[normalizedID] {
			return ErrDuplicateDayIdentifier
		}
		dayIdentifiers[normalizedID] = true

		// Validate intensity level if provided
		if entry.IntensityLevel != nil && *entry.IntensityLevel != "" {
			if !IsValidIntensityLevel(*entry.IntensityLevel) {
				return ErrIntensityLevelInvalid
			}
		}
	}
	return nil
}

// CreateDailyLookupInput contains the input data for creating a new daily lookup.
type CreateDailyLookupInput struct {
	Name      string
	Entries   []DailyLookupEntry
	ProgramID *string
}

// CreateDailyLookup validates input and creates a new DailyLookup domain entity.
// Returns a validation result with errors if validation fails.
func CreateDailyLookup(input CreateDailyLookupInput, id string) (*DailyLookup, *ValidationResult) {
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
	return &DailyLookup{
		ID:        id,
		Name:      strings.TrimSpace(input.Name),
		Entries:   input.Entries,
		ProgramID: input.ProgramID,
		CreatedAt: now,
		UpdatedAt: now,
	}, result
}

// UpdateDailyLookupInput contains the input data for updating an existing daily lookup.
type UpdateDailyLookupInput struct {
	Name      *string             // Optional: only update if provided
	Entries   *[]DailyLookupEntry // Optional: only update if provided
	ProgramID **string            // Optional: only update if provided (double pointer to distinguish nil from clearing)
}

// UpdateDailyLookup validates input and updates an existing DailyLookup.
// Returns a validation result with errors if validation fails.
func UpdateDailyLookup(d *DailyLookup, input UpdateDailyLookupInput) *ValidationResult {
	result := NewValidationResult()

	// Validate name if provided
	if input.Name != nil {
		if err := ValidateName(*input.Name); err != nil {
			result.AddError(err)
		} else {
			d.Name = strings.TrimSpace(*input.Name)
		}
	}

	// Validate entries if provided
	if input.Entries != nil {
		if err := ValidateEntries(*input.Entries); err != nil {
			result.AddError(err)
		} else {
			d.Entries = *input.Entries
		}
	}

	// Update program_id if provided
	if input.ProgramID != nil {
		d.ProgramID = *input.ProgramID
	}

	if result.Valid {
		d.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing daily lookup.
func (d *DailyLookup) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateName(d.Name); err != nil {
		result.AddError(err)
	}

	if err := ValidateEntries(d.Entries); err != nil {
		result.AddError(err)
	}

	return result
}

// GetByDayIdentifier returns the entry for a specific day identifier, or nil if not found.
// The lookup is case-insensitive.
func (d *DailyLookup) GetByDayIdentifier(daySlug string) *DailyLookupEntry {
	normalizedSlug := strings.ToLower(daySlug)
	for i := range d.Entries {
		if strings.ToLower(d.Entries[i].DayIdentifier) == normalizedSlug {
			return &d.Entries[i]
		}
	}
	return nil
}
