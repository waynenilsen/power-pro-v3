// Package rotationlookup provides domain logic for the RotationLookup entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
//
// RotationLookup maps rotation positions to lift identifiers, enabling
// programs like Conjugate/Westside that rotate through different lifts
// across training cycles.
package rotationlookup

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
	ErrEntriesRequired        = errors.New("entries are required")
	ErrPositionInvalid        = errors.New("position must be >= 0")
	ErrLiftIdentifierRequired = errors.New("lift identifier is required")
	ErrDuplicatePosition      = errors.New("duplicate position in entries")
)

// MaxNameLength is the maximum length for a rotation lookup name.
const MaxNameLength = 100

// RotationLookup represents a rotation lookup domain entity with all business rules.
// It maps rotation positions (0-based) to lift identifiers for programs that
// cycle through different lift focuses across training sessions.
type RotationLookup struct {
	ID        string
	Name      string
	Entries   []RotationLookupEntry
	ProgramID *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RotationLookupEntry represents an entry in a rotation lookup table.
// It maps a position to a lift identifier and optional description.
type RotationLookupEntry struct {
	Position       int    `json:"position"`       // 0-based position in rotation
	LiftIdentifier string `json:"liftIdentifier"` // e.g., "deadlift", "squat", "bench"
	Description    string `json:"description"`    // e.g., "Deadlift Focus - High Intensity AMRAP"
}

// ValidationResult is an alias for the shared validation.Result type.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateName validates the rotation lookup name according to business rules.
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

// ValidateEntry validates a single rotation lookup entry.
// Returns an error if validation fails, nil otherwise.
func ValidateEntry(entry RotationLookupEntry) error {
	if entry.Position < 0 {
		return ErrPositionInvalid
	}
	if strings.TrimSpace(entry.LiftIdentifier) == "" {
		return ErrLiftIdentifierRequired
	}
	return nil
}

// ValidateEntries validates the rotation lookup entries.
// Returns an error if validation fails, nil otherwise.
func ValidateEntries(entries []RotationLookupEntry) error {
	if len(entries) == 0 {
		return ErrEntriesRequired
	}

	positions := make(map[int]bool)
	for _, entry := range entries {
		if err := ValidateEntry(entry); err != nil {
			return err
		}
		if positions[entry.Position] {
			return ErrDuplicatePosition
		}
		positions[entry.Position] = true
	}
	return nil
}

// CreateRotationLookupInput contains the input data for creating a new rotation lookup.
type CreateRotationLookupInput struct {
	Name      string
	Entries   []RotationLookupEntry
	ProgramID *string
}

// CreateRotationLookup validates input and creates a new RotationLookup domain entity.
// Returns a validation result with errors if validation fails.
func CreateRotationLookup(input CreateRotationLookupInput, id string) (*RotationLookup, *ValidationResult) {
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
	return &RotationLookup{
		ID:        id,
		Name:      strings.TrimSpace(input.Name),
		Entries:   input.Entries,
		ProgramID: input.ProgramID,
		CreatedAt: now,
		UpdatedAt: now,
	}, result
}

// UpdateRotationLookupInput contains the input data for updating an existing rotation lookup.
type UpdateRotationLookupInput struct {
	Name      *string                 // Optional: only update if provided
	Entries   *[]RotationLookupEntry  // Optional: only update if provided
	ProgramID **string                // Optional: only update if provided (double pointer to distinguish nil from clearing)
}

// UpdateRotationLookup validates input and updates an existing RotationLookup.
// Returns a validation result with errors if validation fails.
func UpdateRotationLookup(r *RotationLookup, input UpdateRotationLookupInput) *ValidationResult {
	result := NewValidationResult()

	// Validate name if provided
	if input.Name != nil {
		if err := ValidateName(*input.Name); err != nil {
			result.AddError(err)
		} else {
			r.Name = strings.TrimSpace(*input.Name)
		}
	}

	// Validate entries if provided
	if input.Entries != nil {
		if err := ValidateEntries(*input.Entries); err != nil {
			result.AddError(err)
		} else {
			r.Entries = *input.Entries
		}
	}

	// Update program_id if provided
	if input.ProgramID != nil {
		r.ProgramID = *input.ProgramID
	}

	if result.Valid {
		r.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing rotation lookup.
func (r *RotationLookup) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateName(r.Name); err != nil {
		result.AddError(err)
	}

	if err := ValidateEntries(r.Entries); err != nil {
		result.AddError(err)
	}

	return result
}

// GetByPosition returns the entry for a specific position, or nil if not found.
func (r *RotationLookup) GetByPosition(position int) *RotationLookupEntry {
	for i := range r.Entries {
		if r.Entries[i].Position == position {
			return &r.Entries[i]
		}
	}
	return nil
}

// Length returns the number of entries in the rotation.
func (r *RotationLookup) Length() int {
	return len(r.Entries)
}

// ContainsLift checks if any entry in the rotation has the given lift identifier.
func (r *RotationLookup) ContainsLift(liftIdentifier string) bool {
	for i := range r.Entries {
		if r.Entries[i].LiftIdentifier == liftIdentifier {
			return true
		}
	}
	return false
}
