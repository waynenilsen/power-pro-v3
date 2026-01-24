// Package loadstrategy provides domain logic for load calculation strategies.
// This package defines the interface for polymorphic load calculation that can
// be extended with new strategies without modifying existing code.
package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// LoadStrategyType identifies the type of load calculation strategy.
// Uses string constants for JSON serialization compatibility.
type LoadStrategyType string

const (
	// TypePercentOf calculates load as a percentage of a reference max.
	TypePercentOf LoadStrategyType = "PERCENT_OF"
	// TypeRPETarget calculates load based on RPE (future implementation).
	TypeRPETarget LoadStrategyType = "RPE_TARGET"
	// TypeFixedWeight uses a fixed weight value (future implementation).
	TypeFixedWeight LoadStrategyType = "FIXED_WEIGHT"
	// TypeRelativeTo calculates load relative to another lift (future implementation).
	TypeRelativeTo LoadStrategyType = "RELATIVE_TO"
	// TypeFindRM indicates the user works up to find their rep max (no prescribed weight).
	TypeFindRM LoadStrategyType = "FIND_RM"
	// TypeTaper applies a taper multiplier to reduce volume as meet approaches.
	// Note: TypeTaper constant is defined in taper.go to avoid circular reference.
)

// ValidStrategyTypes contains all valid strategy types for validation.
// Note: TypeTaper ("TAPER") is also valid but defined in taper.go to avoid import cycles.
var ValidStrategyTypes = map[LoadStrategyType]bool{
	TypePercentOf:   true,
	TypeRPETarget:   true,
	TypeFixedWeight: true,
	TypeRelativeTo:  true,
	TypeFindRM:      true,
	"TAPER":         true,
}

// Errors for load strategy operations.
var (
	ErrUnknownStrategyType = errors.New("unknown load strategy type")
	ErrInvalidParams       = errors.New("invalid load calculation parameters")
	ErrMaxNotFound         = errors.New("max not found for user/lift combination")
	ErrStrategyNotRegistered = errors.New("strategy type not registered in factory")
)

// LoadCalculationParams contains the parameters needed to calculate a load.
type LoadCalculationParams struct {
	// UserID is the UUID of the user for max lookup.
	UserID string
	// LiftID is the UUID of the lift for max lookup.
	LiftID string
	// Context contains additional strategy-specific parameters.
	// This allows strategies to access any context-specific data they need.
	Context map[string]interface{}
	// LookupContext provides week/day context for lookup-based load modifications.
	// Optional: if nil, no lookup modifications are applied.
	LookupContext *LookupContext
}

// Validate validates the LoadCalculationParams.
func (p LoadCalculationParams) Validate() error {
	if p.UserID == "" {
		return fmt.Errorf("%w: user ID is required", ErrInvalidParams)
	}
	if p.LiftID == "" {
		return fmt.Errorf("%w: lift ID is required", ErrInvalidParams)
	}
	return nil
}

// LoadStrategy defines the interface for all load calculation strategies.
// This interface enables polymorphic load calculation using the strategy pattern.
// New strategies can be added by implementing this interface without modifying
// existing code (Open/Closed Principle).
type LoadStrategy interface {
	// Type returns the discriminator string for this strategy.
	// This is used for JSON serialization/deserialization.
	Type() LoadStrategyType

	// CalculateLoad calculates the target weight for a given set of parameters.
	// Returns the calculated weight in the user's preferred unit or an error.
	// The calculation may involve looking up user maxes, applying percentages,
	// rounding to available plate increments, etc.
	CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error)

	// Validate validates the strategy's configuration parameters.
	// Returns an error if the strategy is misconfigured.
	Validate() error
}

// StrategyEnvelope is the JSON wrapper for polymorphic LoadStrategy serialization.
// It uses the discriminated union pattern with a "type" field.
type StrategyEnvelope struct {
	Type LoadStrategyType `json:"type"`
	// Raw contains the strategy-specific JSON data (excluding the type field).
	// This is used during unmarshaling to delegate to the concrete type.
	Raw json.RawMessage `json:"-"`
}

// MarshalJSON implements json.Marshaler for StrategyEnvelope.
// This is typically not used directly; instead, use MarshalStrategy.
func (e *StrategyEnvelope) MarshalJSON() ([]byte, error) {
	type Alias StrategyEnvelope
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

// UnmarshalJSON implements json.Unmarshaler for StrategyEnvelope.
// It extracts the type field and stores the raw JSON for later parsing.
func (e *StrategyEnvelope) UnmarshalJSON(data []byte) error {
	// First, extract just the type field
	var typeOnly struct {
		Type LoadStrategyType `json:"type"`
	}
	if err := json.Unmarshal(data, &typeOnly); err != nil {
		return fmt.Errorf("failed to parse strategy type: %w", err)
	}
	e.Type = typeOnly.Type
	e.Raw = data
	return nil
}

// StrategyFactory creates LoadStrategy instances from their type and JSON data.
// Strategies must be registered with the factory before they can be deserialized.
type StrategyFactory struct {
	// creators maps strategy types to their constructor functions.
	// The constructor receives the raw JSON and returns the concrete strategy.
	creators map[LoadStrategyType]func(json.RawMessage) (LoadStrategy, error)
}

// NewStrategyFactory creates a new StrategyFactory with no registered types.
func NewStrategyFactory() *StrategyFactory {
	return &StrategyFactory{
		creators: make(map[LoadStrategyType]func(json.RawMessage) (LoadStrategy, error)),
	}
}

// Register registers a strategy constructor for a given type.
// The constructor function receives the raw JSON data and returns the concrete strategy.
func (f *StrategyFactory) Register(strategyType LoadStrategyType, creator func(json.RawMessage) (LoadStrategy, error)) {
	f.creators[strategyType] = creator
}

// Create creates a LoadStrategy from a type and raw JSON data.
// Returns ErrStrategyNotRegistered if the type is not registered.
func (f *StrategyFactory) Create(strategyType LoadStrategyType, data json.RawMessage) (LoadStrategy, error) {
	creator, ok := f.creators[strategyType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrStrategyNotRegistered, strategyType)
	}
	return creator(data)
}

// CreateFromJSON creates a LoadStrategy from raw JSON containing the type discriminator.
func (f *StrategyFactory) CreateFromJSON(data json.RawMessage) (LoadStrategy, error) {
	var envelope StrategyEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse strategy envelope: %w", err)
	}
	return f.Create(envelope.Type, data)
}

// IsRegistered checks if a strategy type is registered with the factory.
func (f *StrategyFactory) IsRegistered(strategyType LoadStrategyType) bool {
	_, ok := f.creators[strategyType]
	return ok
}

// RegisteredTypes returns a slice of all registered strategy types.
func (f *StrategyFactory) RegisteredTypes() []LoadStrategyType {
	types := make([]LoadStrategyType, 0, len(f.creators))
	for t := range f.creators {
		types = append(types, t)
	}
	return types
}

// MarshalStrategy serializes a LoadStrategy to JSON with the type discriminator.
// This is a convenience function that ensures the type field is always included.
func MarshalStrategy(strategy LoadStrategy) ([]byte, error) {
	// Most concrete strategies should embed their Type() in their JSON struct,
	// but this provides a fallback mechanism.
	return json.Marshal(strategy)
}

// ValidateStrategyType checks if a strategy type string is valid.
func ValidateStrategyType(strategyType LoadStrategyType) error {
	if strategyType == "" {
		return fmt.Errorf("%w: strategy type is required", ErrUnknownStrategyType)
	}
	if !ValidStrategyTypes[strategyType] {
		return fmt.Errorf("%w: %s", ErrUnknownStrategyType, strategyType)
	}
	return nil
}

// MaxLookup defines the interface for looking up user maxes.
// This interface decouples the load strategy from the persistence layer.
type MaxLookup interface {
	// GetCurrentMax retrieves the most recent max for a user, lift, and max type.
	// Returns nil if no max exists for the combination.
	// maxType should be "ONE_RM" or "TRAINING_MAX".
	GetCurrentMax(ctx context.Context, userID, liftID, maxType string) (*MaxValue, error)
}

// MaxValue represents a max value returned from MaxLookup.
type MaxValue struct {
	Value         float64
	EffectiveDate string
}
