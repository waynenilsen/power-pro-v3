// Package progression provides domain logic for progression strategies.
// This package defines the interface for polymorphic progression handling that can
// be extended with new strategies without modifying existing code.
package progression

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ProgressionType identifies the type of progression strategy.
// Uses string constants for JSON serialization compatibility.
type ProgressionType string

const (
	// TypeLinear adds a fixed increment at regular intervals (per-session or per-week).
	TypeLinear ProgressionType = "LINEAR_PROGRESSION"
	// TypeCycle adds a fixed increment at cycle completion.
	TypeCycle ProgressionType = "CYCLE_PROGRESSION"
	// TypeAMRAP adjusts weight based on AMRAP set performance (e.g., nSuns).
	TypeAMRAP ProgressionType = "AMRAP_PROGRESSION"
	// Future types documented for extensibility:
	// TypeDeloadOnFailure ProgressionType = "DELOAD_ON_FAILURE" - reduces weight after repeated failures
	// TypeRPEBased ProgressionType = "RPE_BASED_PROGRESSION" - adjusts based on RPE targets
	// TypeDouble ProgressionType = "DOUBLE_PROGRESSION" - increases reps first, then weight
)

// ValidProgressionTypes contains all currently implemented progression types.
var ValidProgressionTypes = map[ProgressionType]bool{
	TypeLinear: true,
	TypeCycle:  true,
	TypeAMRAP:  true,
}

// TriggerType identifies what event causes a progression to evaluate/apply.
type TriggerType string

const (
	// TriggerAfterSession fires after a training session is completed.
	TriggerAfterSession TriggerType = "AFTER_SESSION"
	// TriggerAfterWeek fires when advancing from week N to week N+1.
	TriggerAfterWeek TriggerType = "AFTER_WEEK"
	// TriggerAfterCycle fires when a cycle completes (wraps back to week 1).
	TriggerAfterCycle TriggerType = "AFTER_CYCLE"
	// TriggerAfterSet fires immediately after logging an AMRAP set.
	TriggerAfterSet TriggerType = "AFTER_SET"
)

// ValidTriggerTypes contains all valid trigger types for validation.
var ValidTriggerTypes = map[TriggerType]bool{
	TriggerAfterSession: true,
	TriggerAfterWeek:    true,
	TriggerAfterCycle:   true,
	TriggerAfterSet:     true,
}

// MaxType represents the type of max (mirrors liftmax.MaxType for decoupling).
type MaxType string

const (
	OneRM       MaxType = "ONE_RM"
	TrainingMax MaxType = "TRAINING_MAX"
)

// ValidMaxTypes contains all valid max types for validation.
var ValidMaxTypes = map[MaxType]bool{
	OneRM:       true,
	TrainingMax: true,
}

// Errors for progression operations.
var (
	ErrUnknownProgressionType    = errors.New("unknown progression type")
	ErrUnknownTriggerType        = errors.New("unknown trigger type")
	ErrUnknownMaxType            = errors.New("unknown max type")
	ErrInvalidParams             = errors.New("invalid progression parameters")
	ErrProgressionNotRegistered  = errors.New("progression type not registered in factory")
	ErrUserIDRequired            = errors.New("user ID is required")
	ErrLiftIDRequired            = errors.New("lift ID is required")
	ErrMaxTypeRequired           = errors.New("max type is required")
	ErrCurrentValueNotPositive   = errors.New("current value must be positive")
	ErrTriggerEventRequired      = errors.New("trigger event is required")
	ErrIncrementNotPositive      = errors.New("increment must be positive")
)

// TriggerEvent contains context about what triggered a progression evaluation.
type TriggerEvent struct {
	// Type is the type of trigger that fired.
	Type TriggerType `json:"type"`
	// Timestamp is when the trigger event occurred.
	Timestamp time.Time `json:"timestamp"`
	// SessionID is set when Type is AFTER_SESSION (optional).
	SessionID *string `json:"sessionId,omitempty"`
	// WeekNumber is set when Type is AFTER_WEEK or AFTER_CYCLE (optional).
	WeekNumber *int `json:"weekNumber,omitempty"`
	// CycleIteration is set when Type is AFTER_CYCLE (optional).
	CycleIteration *int `json:"cycleIteration,omitempty"`
	// DaySlug identifies which day was completed (optional).
	DaySlug *string `json:"daySlug,omitempty"`
	// LiftsPerformed lists the lift IDs that were part of the completed session (optional).
	LiftsPerformed []string `json:"liftsPerformed,omitempty"`

	// AMRAP-specific fields (for AFTER_SET trigger)
	// RepsPerformed is the number of reps achieved on the AMRAP set.
	RepsPerformed *int `json:"repsPerformed,omitempty"`
	// IsAMRAP indicates whether this was an AMRAP set.
	IsAMRAP bool `json:"isAMRAP,omitempty"`
	// SetWeight is the weight used for the set (optional, for logging context).
	SetWeight *float64 `json:"setWeight,omitempty"`
}

// Validate validates the TriggerEvent.
func (t TriggerEvent) Validate() error {
	if t.Type == "" {
		return fmt.Errorf("%w: trigger type is required", ErrInvalidParams)
	}
	if !ValidTriggerTypes[t.Type] {
		return fmt.Errorf("%w: %s", ErrUnknownTriggerType, t.Type)
	}
	if t.Timestamp.IsZero() {
		return fmt.Errorf("%w: trigger timestamp is required", ErrInvalidParams)
	}
	return nil
}

// ProgressionContext contains all parameters needed to evaluate and apply a progression.
type ProgressionContext struct {
	// UserID is the UUID of the user for LiftMax lookup.
	UserID string `json:"userId"`
	// LiftID is the UUID of the lift for LiftMax lookup.
	LiftID string `json:"liftId"`
	// MaxType specifies which max type to modify (ONE_RM or TRAINING_MAX).
	MaxType MaxType `json:"maxType"`
	// CurrentValue is the current LiftMax value.
	CurrentValue float64 `json:"currentValue"`
	// TriggerEvent provides context about what triggered this progression.
	TriggerEvent TriggerEvent `json:"triggerEvent"`
}

// Validate validates the ProgressionContext.
func (c ProgressionContext) Validate() error {
	if c.UserID == "" {
		return ErrUserIDRequired
	}
	if c.LiftID == "" {
		return ErrLiftIDRequired
	}
	if c.MaxType == "" {
		return ErrMaxTypeRequired
	}
	if !ValidMaxTypes[c.MaxType] {
		return fmt.Errorf("%w: %s", ErrUnknownMaxType, c.MaxType)
	}
	if c.CurrentValue <= 0 {
		return ErrCurrentValueNotPositive
	}
	if err := c.TriggerEvent.Validate(); err != nil {
		return fmt.Errorf("invalid trigger event: %w", err)
	}
	return nil
}

// ProgressionResult contains the outcome of applying a progression.
type ProgressionResult struct {
	// Applied indicates whether the progression was actually applied.
	// May be false if conditions weren't met (e.g., wrong trigger type).
	Applied bool `json:"applied"`
	// PreviousValue is the LiftMax value before progression.
	PreviousValue float64 `json:"previousValue"`
	// NewValue is the LiftMax value after progression.
	NewValue float64 `json:"newValue"`
	// Delta is the increment that was applied (NewValue - PreviousValue).
	Delta float64 `json:"delta"`
	// LiftID identifies which lift was modified.
	LiftID string `json:"liftId"`
	// MaxType identifies which max type was modified.
	MaxType MaxType `json:"maxType"`
	// AppliedAt is when the progression was applied.
	AppliedAt time.Time `json:"appliedAt"`
	// Reason provides context when Applied is false (optional).
	Reason string `json:"reason,omitempty"`
}

// Progression defines the interface for all progression strategies.
// This interface enables polymorphic progression handling using the strategy pattern.
// New progressions can be added by implementing this interface without modifying
// existing code (Open/Closed Principle).
type Progression interface {
	// Type returns the discriminator string for this progression.
	// This is used for JSON serialization/deserialization.
	Type() ProgressionType

	// Apply evaluates and applies the progression given the context.
	// Returns a ProgressionResult indicating whether progression occurred and the details.
	// The method should be idempotent when combined with ProgressionLog checks.
	Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error)

	// Validate validates the progression's configuration parameters.
	// Returns an error if the progression is misconfigured.
	Validate() error

	// TriggerType returns the trigger type this progression responds to.
	// Used to filter progressions by the event that fired.
	TriggerType() TriggerType
}

// ProgressionEnvelope is the JSON wrapper for polymorphic Progression serialization.
// It uses the discriminated union pattern with a "type" field.
type ProgressionEnvelope struct {
	Type ProgressionType `json:"type"`
	// Raw contains the progression-specific JSON data (excluding the type field).
	// This is used during unmarshaling to delegate to the concrete type.
	Raw json.RawMessage `json:"-"`
}

// MarshalJSON implements json.Marshaler for ProgressionEnvelope.
// This is typically not used directly; instead, use MarshalProgression.
func (e *ProgressionEnvelope) MarshalJSON() ([]byte, error) {
	type Alias ProgressionEnvelope
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

// UnmarshalJSON implements json.Unmarshaler for ProgressionEnvelope.
// It extracts the type field and stores the raw JSON for later parsing.
func (e *ProgressionEnvelope) UnmarshalJSON(data []byte) error {
	// First, extract just the type field
	var typeOnly struct {
		Type ProgressionType `json:"type"`
	}
	if err := json.Unmarshal(data, &typeOnly); err != nil {
		return fmt.Errorf("failed to parse progression type: %w", err)
	}
	e.Type = typeOnly.Type
	e.Raw = data
	return nil
}

// ProgressionFactory creates Progression instances from their type and JSON data.
// Progressions must be registered with the factory before they can be deserialized.
type ProgressionFactory struct {
	// creators maps progression types to their constructor functions.
	// The constructor receives the raw JSON and returns the concrete progression.
	creators map[ProgressionType]func(json.RawMessage) (Progression, error)
}

// NewProgressionFactory creates a new ProgressionFactory with no registered types.
func NewProgressionFactory() *ProgressionFactory {
	return &ProgressionFactory{
		creators: make(map[ProgressionType]func(json.RawMessage) (Progression, error)),
	}
}

// Register registers a progression constructor for a given type.
// The constructor function receives the raw JSON data and returns the concrete progression.
func (f *ProgressionFactory) Register(progressionType ProgressionType, creator func(json.RawMessage) (Progression, error)) {
	f.creators[progressionType] = creator
}

// Create creates a Progression from a type and raw JSON data.
// Returns ErrProgressionNotRegistered if the type is not registered.
func (f *ProgressionFactory) Create(progressionType ProgressionType, data json.RawMessage) (Progression, error) {
	creator, ok := f.creators[progressionType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProgressionNotRegistered, progressionType)
	}
	return creator(data)
}

// CreateFromJSON creates a Progression from raw JSON containing the type discriminator.
func (f *ProgressionFactory) CreateFromJSON(data json.RawMessage) (Progression, error) {
	var envelope ProgressionEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse progression envelope: %w", err)
	}
	return f.Create(envelope.Type, data)
}

// IsRegistered checks if a progression type is registered with the factory.
func (f *ProgressionFactory) IsRegistered(progressionType ProgressionType) bool {
	_, ok := f.creators[progressionType]
	return ok
}

// RegisteredTypes returns a slice of all registered progression types.
func (f *ProgressionFactory) RegisteredTypes() []ProgressionType {
	types := make([]ProgressionType, 0, len(f.creators))
	for t := range f.creators {
		types = append(types, t)
	}
	return types
}

// MarshalProgression serializes a Progression to JSON with the type discriminator.
// This is a convenience function that ensures the type field is always included.
func MarshalProgression(progression Progression) ([]byte, error) {
	// Most concrete progressions should embed their Type() in their JSON struct,
	// but this provides a fallback mechanism.
	return json.Marshal(progression)
}

// ValidateProgressionType checks if a progression type string is valid.
func ValidateProgressionType(progressionType ProgressionType) error {
	if progressionType == "" {
		return fmt.Errorf("%w: progression type is required", ErrUnknownProgressionType)
	}
	if !ValidProgressionTypes[progressionType] {
		return fmt.Errorf("%w: %s", ErrUnknownProgressionType, progressionType)
	}
	return nil
}

// ValidateTriggerType checks if a trigger type string is valid.
func ValidateTriggerType(triggerType TriggerType) error {
	if triggerType == "" {
		return fmt.Errorf("%w: trigger type is required", ErrUnknownTriggerType)
	}
	if !ValidTriggerTypes[triggerType] {
		return fmt.Errorf("%w: %s", ErrUnknownTriggerType, triggerType)
	}
	return nil
}

// ValidateMaxType checks if a max type string is valid.
func ValidateMaxType(maxType MaxType) error {
	if maxType == "" {
		return fmt.Errorf("%w: max type is required", ErrUnknownMaxType)
	}
	if !ValidMaxTypes[maxType] {
		return fmt.Errorf("%w: %s", ErrUnknownMaxType, maxType)
	}
	return nil
}

// NewProgression is a convenience function that creates a Progression from a type and raw params.
// This function requires a properly initialized factory with registered progression types.
// For a standalone factory creation, use NewProgressionFactory and register types manually,
// or use DefaultProgressionFactory which has all built-in types registered.
func NewProgression(factory *ProgressionFactory, progressionType string, params json.RawMessage) (Progression, error) {
	pt := ProgressionType(progressionType)
	if err := ValidateProgressionType(pt); err != nil {
		return nil, err
	}
	return factory.Create(pt, params)
}
