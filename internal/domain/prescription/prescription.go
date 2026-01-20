// Package prescription provides domain logic for the Prescription entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package prescription

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// Validation errors
var (
	ErrLiftIDRequired       = errors.New("lift ID is required")
	ErrLiftIDInvalid        = errors.New("lift ID must be a valid UUID format")
	ErrLoadStrategyRequired = errors.New("load strategy is required")
	ErrSetSchemeRequired    = errors.New("set scheme is required")
	ErrOrderNegative        = errors.New("order must be >= 0")
	ErrNotesTooLong         = errors.New("notes must be 500 characters or less")
	ErrRestSecondsNegative  = errors.New("rest seconds must be >= 0 when provided")
	ErrLiftNotFound         = errors.New("lift not found")
	ErrMaxNotFound          = errors.New("max not found for user/lift combination")
)

// Max length for notes field
const MaxNotesLength = 500

// LiftInfo contains the minimal lift information needed for resolved prescriptions.
type LiftInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// LiftLookup defines the interface for looking up lift information.
// This interface decouples the prescription domain from the persistence layer.
type LiftLookup interface {
	// GetLiftByID retrieves lift information by ID.
	// Returns nil if the lift does not exist.
	GetLiftByID(ctx context.Context, liftID string) (*LiftInfo, error)
}

// Prescription represents a prescription domain entity with all business rules.
// A prescription links a lift to load and set specifications.
type Prescription struct {
	ID           string
	LiftID       string
	LoadStrategy loadstrategy.LoadStrategy
	SetScheme    setscheme.SetScheme
	Order        int
	Notes        string
	RestSeconds  *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ValidationResult contains the result of validating a prescription.
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

// ValidateLiftID validates the lift ID according to business rules.
// Returns an error if validation fails, nil otherwise.
func ValidateLiftID(liftID string) error {
	trimmed := strings.TrimSpace(liftID)
	if trimmed == "" {
		return ErrLiftIDRequired
	}
	// Basic UUID format check (8-4-4-4-12 hex characters)
	if !isValidUUID(trimmed) {
		return ErrLiftIDInvalid
	}
	return nil
}

// isValidUUID performs basic UUID format validation.
// Accepts UUIDs in the standard format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	// Check hyphen positions
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return false
	}
	// Check that all other characters are hex
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue // Skip hyphens
		}
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// ValidateLoadStrategy validates the load strategy.
func ValidateLoadStrategy(strategy loadstrategy.LoadStrategy) error {
	if strategy == nil {
		return ErrLoadStrategyRequired
	}
	return strategy.Validate()
}

// ValidateSetScheme validates the set scheme.
func ValidateSetScheme(scheme setscheme.SetScheme) error {
	if scheme == nil {
		return ErrSetSchemeRequired
	}
	return scheme.Validate()
}

// ValidateOrder validates the order field.
func ValidateOrder(order int) error {
	if order < 0 {
		return ErrOrderNegative
	}
	return nil
}

// ValidateNotes validates the notes field.
func ValidateNotes(notes string) error {
	if len(notes) > MaxNotesLength {
		return ErrNotesTooLong
	}
	return nil
}

// ValidateRestSeconds validates the rest seconds field.
func ValidateRestSeconds(restSeconds *int) error {
	if restSeconds != nil && *restSeconds < 0 {
		return ErrRestSecondsNegative
	}
	return nil
}

// CreatePrescriptionInput contains the input data for creating a new prescription.
type CreatePrescriptionInput struct {
	LiftID       string
	LoadStrategy loadstrategy.LoadStrategy
	SetScheme    setscheme.SetScheme
	Order        int     // Defaults to 0
	Notes        string  // Optional
	RestSeconds  *int    // Optional
}

// CreatePrescription validates input and creates a new Prescription domain entity.
// Returns a validation result with errors if validation fails.
func CreatePrescription(input CreatePrescriptionInput, id string) (*Prescription, *ValidationResult) {
	result := NewValidationResult()

	// Validate lift ID
	if err := ValidateLiftID(input.LiftID); err != nil {
		result.AddError(err)
	}

	// Validate load strategy
	if err := ValidateLoadStrategy(input.LoadStrategy); err != nil {
		result.AddError(err)
	}

	// Validate set scheme
	if err := ValidateSetScheme(input.SetScheme); err != nil {
		result.AddError(err)
	}

	// Validate order
	if err := ValidateOrder(input.Order); err != nil {
		result.AddError(err)
	}

	// Validate notes
	if err := ValidateNotes(input.Notes); err != nil {
		result.AddError(err)
	}

	// Validate rest seconds
	if err := ValidateRestSeconds(input.RestSeconds); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &Prescription{
		ID:           id,
		LiftID:       input.LiftID,
		LoadStrategy: input.LoadStrategy,
		SetScheme:    input.SetScheme,
		Order:        input.Order,
		Notes:        input.Notes,
		RestSeconds:  input.RestSeconds,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, result
}

// UpdatePrescriptionInput contains the input data for updating an existing prescription.
type UpdatePrescriptionInput struct {
	LiftID          *string                       // Optional: only update if provided
	LoadStrategy    loadstrategy.LoadStrategy     // Optional: only update if non-nil
	SetScheme       setscheme.SetScheme           // Optional: only update if non-nil
	Order           *int                          // Optional: only update if provided
	Notes           *string                       // Optional: only update if provided
	RestSeconds     *int                          // Optional: only update if provided
	ClearRestSeconds bool                         // Set to true to explicitly clear rest seconds
}

// UpdatePrescription validates input and updates an existing Prescription.
// Returns a validation result with errors if validation fails.
func UpdatePrescription(prescription *Prescription, input UpdatePrescriptionInput) *ValidationResult {
	result := NewValidationResult()

	// Validate lift ID if provided
	if input.LiftID != nil {
		if err := ValidateLiftID(*input.LiftID); err != nil {
			result.AddError(err)
		} else {
			prescription.LiftID = *input.LiftID
		}
	}

	// Validate load strategy if provided
	if input.LoadStrategy != nil {
		if err := ValidateLoadStrategy(input.LoadStrategy); err != nil {
			result.AddError(err)
		} else {
			prescription.LoadStrategy = input.LoadStrategy
		}
	}

	// Validate set scheme if provided
	if input.SetScheme != nil {
		if err := ValidateSetScheme(input.SetScheme); err != nil {
			result.AddError(err)
		} else {
			prescription.SetScheme = input.SetScheme
		}
	}

	// Validate order if provided
	if input.Order != nil {
		if err := ValidateOrder(*input.Order); err != nil {
			result.AddError(err)
		} else {
			prescription.Order = *input.Order
		}
	}

	// Validate notes if provided
	if input.Notes != nil {
		if err := ValidateNotes(*input.Notes); err != nil {
			result.AddError(err)
		} else {
			prescription.Notes = *input.Notes
		}
	}

	// Handle rest seconds update
	if input.ClearRestSeconds {
		prescription.RestSeconds = nil
	} else if input.RestSeconds != nil {
		if err := ValidateRestSeconds(input.RestSeconds); err != nil {
			result.AddError(err)
		} else {
			prescription.RestSeconds = input.RestSeconds
		}
	}

	if result.Valid {
		prescription.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing prescription.
func (p *Prescription) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateLiftID(p.LiftID); err != nil {
		result.AddError(err)
	}

	if err := ValidateLoadStrategy(p.LoadStrategy); err != nil {
		result.AddError(err)
	}

	if err := ValidateSetScheme(p.SetScheme); err != nil {
		result.AddError(err)
	}

	if err := ValidateOrder(p.Order); err != nil {
		result.AddError(err)
	}

	if err := ValidateNotes(p.Notes); err != nil {
		result.AddError(err)
	}

	if err := ValidateRestSeconds(p.RestSeconds); err != nil {
		result.AddError(err)
	}

	return result
}

// ResolvedPrescription represents a fully resolved prescription with concrete sets.
// This is the output of the Resolve method.
type ResolvedPrescription struct {
	PrescriptionID string                 `json:"prescriptionId"`
	Lift           LiftInfo               `json:"lift"`
	Sets           []setscheme.GeneratedSet `json:"sets"`
	Notes          string                 `json:"notes,omitempty"`
	RestSeconds    *int                   `json:"restSeconds,omitempty"`
}

// ResolutionContext provides dependencies needed for prescription resolution.
type ResolutionContext struct {
	LiftLookup     LiftLookup
	SetGenContext  setscheme.SetGenerationContext
}

// DefaultResolutionContext returns a ResolutionContext with default values.
func DefaultResolutionContext(liftLookup LiftLookup) ResolutionContext {
	return ResolutionContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}
}

// Resolve executes the prescription resolution workflow:
// 1. Looks up the lift information
// 2. Calls LoadStrategy.CalculateLoad to get base weight
// 3. Calls SetScheme.GenerateSets with base weight
// 4. Returns ResolvedPrescription with sets, notes, and rest
func (p *Prescription) Resolve(ctx context.Context, userID string, resCtx ResolutionContext) (*ResolvedPrescription, error) {
	// Validate userID
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("user ID is required for prescription resolution")
	}

	// Look up lift information
	var liftInfo LiftInfo
	if resCtx.LiftLookup != nil {
		lift, err := resCtx.LiftLookup.GetLiftByID(ctx, p.LiftID)
		if err != nil {
			return nil, fmt.Errorf("failed to look up lift: %w", err)
		}
		if lift == nil {
			return nil, fmt.Errorf("%w: lift ID %s", ErrLiftNotFound, p.LiftID)
		}
		liftInfo = *lift
	} else {
		// If no lift lookup is provided, use minimal info
		liftInfo = LiftInfo{ID: p.LiftID}
	}

	// Calculate base weight using load strategy
	loadParams := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: p.LiftID,
	}
	baseWeight, err := p.LoadStrategy.CalculateLoad(ctx, loadParams)
	if err != nil {
		// Check if it's a "max not found" error for clearer messaging
		if errors.Is(err, loadstrategy.ErrMaxNotFound) {
			return nil, fmt.Errorf("%w: unable to calculate load for lift %s", ErrMaxNotFound, p.LiftID)
		}
		return nil, fmt.Errorf("failed to calculate load: %w", err)
	}

	// Generate sets using set scheme
	sets, err := p.SetScheme.GenerateSets(baseWeight, resCtx.SetGenContext)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sets: %w", err)
	}

	return &ResolvedPrescription{
		PrescriptionID: p.ID,
		Lift:           liftInfo,
		Sets:           sets,
		Notes:          p.Notes,
		RestSeconds:    p.RestSeconds,
	}, nil
}
