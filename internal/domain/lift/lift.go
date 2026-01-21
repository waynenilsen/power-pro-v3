// Package lift provides domain logic for the Lift entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package lift

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// MaxSlugLength is the maximum allowed length for lift slugs.
const MaxSlugLength = 100

// Validation errors
var (
	ErrNameRequired      = errors.New("lift name is required")
	ErrNameTooLong       = errors.New("lift name must be 100 characters or less")
	ErrCircularReference = errors.New("circular reference detected: lift cannot be its own ancestor")
	ErrSelfReference     = errors.New("lift cannot reference itself as parent")
	// Slug errors delegated to shared validation package
	ErrSlugEmpty   = validation.ErrSlugEmpty
	ErrSlugInvalid = validation.ErrSlugInvalid
	ErrSlugTooLong = validation.SlugTooLongError(MaxSlugLength)
)

// Lift represents a lift domain entity with all business rules.
type Lift struct {
	ID                string
	Name              string
	Slug              string
	IsCompetitionLift bool
	ParentLiftID      *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// LiftRepository defines the interface for lift persistence operations.
// This interface is used for circular reference detection.
type LiftRepository interface {
	// GetByID retrieves a lift by its ID. Returns nil if not found.
	GetByID(id string) (*Lift, error)
	// SlugExists checks if a slug already exists in the repository.
	SlugExists(slug string, excludeID *string) (bool, error)
}

// ValidationResult is an alias for the shared validation.Result type.
// This maintains backward compatibility while using the shared implementation.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateName validates the lift name according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrNameRequired
	}
	if len(name) > 100 {
		return ErrNameTooLong
	}
	return nil
}

// ValidateSlug validates the lift slug according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateSlug(slug string) error {
	return validation.ValidateSlug(slug, MaxSlugLength)
}

// GenerateSlug creates a URL-safe slug from a name.
// Delegates to the shared validation.GenerateSlug function.
func GenerateSlug(name string) string {
	return validation.GenerateSlug(name)
}

// ValidateParentLiftID validates the parent lift reference.
// It checks for self-reference and uses the repository to detect circular references.
func ValidateParentLiftID(liftID string, parentLiftID *string, repo LiftRepository) error {
	if parentLiftID == nil {
		return nil
	}

	// Check for self-reference
	if *parentLiftID == liftID {
		return ErrSelfReference
	}

	// Check for circular reference using repository
	if repo != nil {
		if err := detectCircularReference(liftID, *parentLiftID, repo); err != nil {
			return err
		}
	}

	return nil
}

// detectCircularReference recursively checks if setting parentLiftID as the parent
// of liftID would create a circular reference.
func detectCircularReference(liftID, parentLiftID string, repo LiftRepository) error {
	visited := make(map[string]bool)
	visited[liftID] = true

	currentID := parentLiftID
	for currentID != "" {
		if visited[currentID] {
			return ErrCircularReference
		}
		visited[currentID] = true

		parentLift, err := repo.GetByID(currentID)
		if err != nil {
			return fmt.Errorf("failed to check circular reference: %w", err)
		}
		if parentLift == nil {
			break
		}
		if parentLift.ParentLiftID == nil {
			break
		}
		currentID = *parentLift.ParentLiftID
	}

	return nil
}

// CreateLiftInput contains the input data for creating a new lift.
type CreateLiftInput struct {
	Name              string
	Slug              string  // Optional: auto-generated from Name if empty
	IsCompetitionLift bool    // Defaults to false
	ParentLiftID      *string // Optional
}

// CreateLift validates input and creates a new Lift domain entity.
// It auto-generates a slug from the name if not provided.
// Returns a validation result with errors if validation fails.
func CreateLift(input CreateLiftInput, id string, repo LiftRepository) (*Lift, *ValidationResult) {
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

	// Validate parent lift reference (circular reference detection)
	if err := ValidateParentLiftID(id, input.ParentLiftID, repo); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &Lift{
		ID:                id,
		Name:              input.Name,
		Slug:              slug,
		IsCompetitionLift: input.IsCompetitionLift,
		ParentLiftID:      input.ParentLiftID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, result
}

// UpdateLiftInput contains the input data for updating an existing lift.
type UpdateLiftInput struct {
	Name              *string // Optional: only update if provided
	Slug              *string // Optional: only update if provided
	IsCompetitionLift *bool   // Optional: only update if provided
	ParentLiftID      *string // Optional: use empty string to clear, nil to leave unchanged
	ClearParentLift   bool    // Set to true to explicitly clear the parent lift
}

// UpdateLift validates input and updates an existing Lift.
// Returns a validation result with errors if validation fails.
func UpdateLift(lift *Lift, input UpdateLiftInput, repo LiftRepository) *ValidationResult {
	result := NewValidationResult()

	// Validate name if provided
	if input.Name != nil {
		if err := ValidateName(*input.Name); err != nil {
			result.AddError(err)
		} else {
			lift.Name = *input.Name
		}
	}

	// Validate slug if provided
	if input.Slug != nil {
		if err := ValidateSlug(*input.Slug); err != nil {
			result.AddError(err)
		} else {
			lift.Slug = *input.Slug
		}
	}

	// Update competition lift flag if provided
	if input.IsCompetitionLift != nil {
		lift.IsCompetitionLift = *input.IsCompetitionLift
	}

	// Handle parent lift update
	if input.ClearParentLift {
		lift.ParentLiftID = nil
	} else if input.ParentLiftID != nil {
		if err := ValidateParentLiftID(lift.ID, input.ParentLiftID, repo); err != nil {
			result.AddError(err)
		} else {
			lift.ParentLiftID = input.ParentLiftID
		}
	}

	if result.Valid {
		lift.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing lift.
func (l *Lift) Validate(repo LiftRepository) *ValidationResult {
	result := NewValidationResult()

	if err := ValidateName(l.Name); err != nil {
		result.AddError(err)
	}

	if err := ValidateSlug(l.Slug); err != nil {
		result.AddError(err)
	}

	if err := ValidateParentLiftID(l.ID, l.ParentLiftID, repo); err != nil {
		result.AddError(err)
	}

	return result
}
