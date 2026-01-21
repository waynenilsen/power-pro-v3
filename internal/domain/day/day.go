// Package day provides domain logic for the Day entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package day

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// MaxSlugLength is the maximum allowed length for day slugs.
const MaxSlugLength = 50

// Validation errors
var (
	ErrNameRequired         = errors.New("day name is required")
	ErrNameTooLong          = errors.New("day name must be 50 characters or less")
	ErrMetadataInvalidJSON  = errors.New("day metadata must be valid JSON")
	ErrOrderNegative        = errors.New("order must be >= 0")
	ErrPrescriptionNotFound = errors.New("prescription not found in day")
	// Slug errors delegated to shared validation package
	ErrSlugEmpty   = validation.ErrSlugEmpty
	ErrSlugInvalid = validation.ErrSlugInvalid
	ErrSlugTooLong = validation.SlugTooLongError(MaxSlugLength)
)

// Metadata keys
const (
	MetadataKeyIntensityLevel = "intensityLevel"
	MetadataKeyFocus          = "focus"
)

// Valid intensity levels
const (
	IntensityLevelHeavy  = "HEAVY"
	IntensityLevelLight  = "LIGHT"
	IntensityLevelMedium = "MEDIUM"
)


// Day represents a day domain entity with all business rules.
type Day struct {
	ID        string
	Name      string
	Slug      string
	Metadata  map[string]interface{}
	ProgramID *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DayPrescription represents a prescription associated with a day.
type DayPrescription struct {
	ID             string
	DayID          string
	PrescriptionID string
	Order          int
	CreatedAt      time.Time
}

// DayWithPrescriptions represents a day with its associated prescriptions.
type DayWithPrescriptions struct {
	Day           *Day
	Prescriptions []DayPrescription
}

// ValidationResult is an alias for the shared validation.Result type.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateName validates the day name according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrNameRequired
	}
	if len(name) > 50 {
		return ErrNameTooLong
	}
	return nil
}

// ValidateSlug validates the day slug according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateSlug(slug string) error {
	return validation.ValidateSlug(slug, MaxSlugLength)
}

// ValidateMetadata validates the day metadata.
// Returns an error if validation fails, nil otherwise.
func ValidateMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		return nil
	}
	// Validate JSON serialization
	_, err := json.Marshal(metadata)
	if err != nil {
		return ErrMetadataInvalidJSON
	}
	return nil
}

// ValidateOrder validates the prescription order.
// Returns an error if validation fails, nil otherwise.
func ValidateOrder(order int) error {
	if order < 0 {
		return ErrOrderNegative
	}
	return nil
}

// GenerateSlug creates a URL-safe slug from a name.
// Delegates to the shared validation.GenerateSlug function.
func GenerateSlug(name string) string {
	return validation.GenerateSlug(name)
}

// CreateDayInput contains the input data for creating a new day.
type CreateDayInput struct {
	Name      string
	Slug      string                 // Optional: auto-generated from Name if empty
	Metadata  map[string]interface{} // Optional
	ProgramID *string                // Optional
}

// CreateDay validates input and creates a new Day domain entity.
// It auto-generates a slug from the name if not provided.
// Returns a validation result with errors if validation fails.
func CreateDay(input CreateDayInput, id string) (*Day, *ValidationResult) {
	result := NewValidationResult()

	// Validate name
	if err := ValidateName(input.Name); err != nil {
		result.AddError(err)
	}

	// Auto-generate slug if not provided
	slug := input.Slug
	if slug == "" {
		slug = GenerateSlug(input.Name)
	}

	// Validate slug
	if err := ValidateSlug(slug); err != nil {
		result.AddError(err)
	}

	// Validate metadata
	if err := ValidateMetadata(input.Metadata); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &Day{
		ID:        id,
		Name:      input.Name,
		Slug:      slug,
		Metadata:  input.Metadata,
		ProgramID: input.ProgramID,
		CreatedAt: now,
		UpdatedAt: now,
	}, result
}

// UpdateDayInput contains the input data for updating an existing day.
type UpdateDayInput struct {
	Name           *string                 // Optional: only update if provided
	Slug           *string                 // Optional: only update if provided
	Metadata       map[string]interface{}  // Optional: only update if provided
	ClearMetadata  bool                    // Set to true to explicitly clear metadata
	ProgramID      *string                 // Optional: only update if provided
	ClearProgramID bool                    // Set to true to explicitly clear program ID
}

// UpdateDay validates input and updates an existing Day.
// Returns a validation result with errors if validation fails.
func UpdateDay(day *Day, input UpdateDayInput) *ValidationResult {
	result := NewValidationResult()

	// Validate name if provided
	if input.Name != nil {
		if err := ValidateName(*input.Name); err != nil {
			result.AddError(err)
		} else {
			day.Name = *input.Name
		}
	}

	// Validate slug if provided
	if input.Slug != nil {
		if err := ValidateSlug(*input.Slug); err != nil {
			result.AddError(err)
		} else {
			day.Slug = *input.Slug
		}
	}

	// Handle metadata update
	if input.ClearMetadata {
		day.Metadata = nil
	} else if input.Metadata != nil {
		if err := ValidateMetadata(input.Metadata); err != nil {
			result.AddError(err)
		} else {
			day.Metadata = input.Metadata
		}
	}

	// Handle program ID update
	if input.ClearProgramID {
		day.ProgramID = nil
	} else if input.ProgramID != nil {
		day.ProgramID = input.ProgramID
	}

	if result.Valid {
		day.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing day.
func (d *Day) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateName(d.Name); err != nil {
		result.AddError(err)
	}

	if err := ValidateSlug(d.Slug); err != nil {
		result.AddError(err)
	}

	if err := ValidateMetadata(d.Metadata); err != nil {
		result.AddError(err)
	}

	return result
}

// CreateDayPrescriptionInput contains the input data for adding a prescription to a day.
type CreateDayPrescriptionInput struct {
	DayID          string
	PrescriptionID string
	Order          *int // Optional: if nil, will be set to max order + 1
}

// CreateDayPrescription validates input and creates a new DayPrescription.
// Returns a validation result with errors if validation fails.
func CreateDayPrescription(input CreateDayPrescriptionInput, id string, order int) (*DayPrescription, *ValidationResult) {
	result := NewValidationResult()

	if input.DayID == "" {
		result.AddError(errors.New("day ID is required"))
	}

	if input.PrescriptionID == "" {
		result.AddError(errors.New("prescription ID is required"))
	}

	if err := ValidateOrder(order); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	return &DayPrescription{
		ID:             id,
		DayID:          input.DayID,
		PrescriptionID: input.PrescriptionID,
		Order:          order,
		CreatedAt:      time.Now(),
	}, result
}

// ReorderPrescriptionsInput contains the input for reordering prescriptions in a day.
type ReorderPrescriptionsInput struct {
	DayID          string
	PrescriptionIDs []string // Ordered list of prescription IDs
}

// ValidateReorderInput validates the reorder input.
func ValidateReorderInput(input ReorderPrescriptionsInput) *ValidationResult {
	result := NewValidationResult()

	if input.DayID == "" {
		result.AddError(errors.New("day ID is required"))
	}

	if len(input.PrescriptionIDs) == 0 {
		result.AddError(errors.New("prescription IDs are required"))
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, id := range input.PrescriptionIDs {
		if id == "" {
			result.AddError(errors.New("prescription ID cannot be empty"))
			continue
		}
		if seen[id] {
			result.AddError(fmt.Errorf("duplicate prescription ID: %s", id))
		}
		seen[id] = true
	}

	return result
}
