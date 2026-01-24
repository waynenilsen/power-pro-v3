// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// SessionLookup defines the interface for looking up logged sets within a session.
// This interface decouples the RelativeTo strategy from the persistence layer.
type SessionLookup interface {
	// GetLoggedSetByIndex retrieves a logged set by its index within a session.
	// setIndex is 0-based (0 = first/top set).
	// Returns nil if the set has not yet been logged.
	GetLoggedSetByIndex(ctx context.Context, sessionID string, liftID string, setIndex int) (*LoggedSetResult, error)
}

// LoggedSetResult represents the result of a logged set lookup.
type LoggedSetResult struct {
	Weight float64
	Reps   int
	RPE    *float64
}

// RelativeTo validation errors.
var (
	ErrReferenceSetIndexInvalid = errors.New("reference set index must be non-negative")
	ErrRelativeToPercentageInvalid = errors.New("percentage must be greater than 0")
	ErrSessionLookupRequired    = errors.New("session lookup is required for RELATIVE_TO strategy")
	ErrSessionIDRequired        = errors.New("session ID is required in context for RELATIVE_TO strategy")
	ErrReferenceSetNotFound     = errors.New("reference set not found (may not be logged yet)")
)

// RelativeToLoadStrategy calculates load as a percentage of a reference set from the current session.
// This is commonly used for back-off sets after finding a rep max or hitting an RPE-based top set.
//
// Example: 85% of top set weight with 5lb rounding
//   - Top set (index 0) weight: 400 lbs
//   - Calculated: 400 * 0.85 = 340 lbs
//   - Rounded to nearest 5: 340 lbs
type RelativeToLoadStrategy struct {
	// ReferenceSetIndex specifies which set to reference (0 = first/top set).
	ReferenceSetIndex int `json:"referenceSetIndex"`

	// Percentage is the percentage of the reference weight (e.g., 85 for 85%).
	// Must be > 0. Values > 100 are allowed.
	Percentage float64 `json:"percentage"`

	// RoundingIncrement is the weight increment for rounding (e.g., 2.5, 5.0).
	// Optional; defaults to 5.0 if not specified or <= 0.
	RoundingIncrement float64 `json:"roundingIncrement,omitempty"`

	// RoundingDirection specifies how to round (NEAREST, DOWN, UP).
	// Optional; defaults to NEAREST if not specified.
	RoundingDirection RoundingDirection `json:"roundingDirection,omitempty"`

	// sessionLookup is the repository for looking up logged sets.
	// This is injected and not serialized.
	sessionLookup SessionLookup `json:"-"`
}

// NewRelativeToLoadStrategy creates a new RelativeToLoadStrategy with the given parameters.
func NewRelativeToLoadStrategy(
	referenceSetIndex int,
	percentage float64,
	roundingIncrement float64,
	roundingDirection RoundingDirection,
	sessionLookup SessionLookup,
) *RelativeToLoadStrategy {
	return &RelativeToLoadStrategy{
		ReferenceSetIndex: referenceSetIndex,
		Percentage:        percentage,
		RoundingIncrement: roundingIncrement,
		RoundingDirection: roundingDirection,
		sessionLookup:     sessionLookup,
	}
}

// Type returns the strategy type discriminator.
func (s *RelativeToLoadStrategy) Type() LoadStrategyType {
	return TypeRelativeTo
}

// CalculateLoad calculates the target weight based on a percentage of a reference set.
// It looks up the referenced set from the current session and applies the percentage.
//
// The calculation follows these steps:
//  1. Validate parameters
//  2. Extract sessionID from params.Context
//  3. Look up the referenced set via sessionLookup
//  4. Calculate: weight = referenceWeight * (percentage / 100)
//  5. Round to configured increment
func (s *RelativeToLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
	// Validate params
	if err := params.Validate(); err != nil {
		return 0, err
	}

	// Validate strategy configuration
	if err := s.Validate(); err != nil {
		return 0, err
	}

	// Check that sessionLookup is available
	if s.sessionLookup == nil {
		return 0, ErrSessionLookupRequired
	}

	// Extract sessionID from context
	sessionID, err := s.extractSessionID(params.Context)
	if err != nil {
		return 0, err
	}

	// Look up the referenced set
	loggedSet, err := s.sessionLookup.GetLoggedSetByIndex(ctx, sessionID, params.LiftID, s.ReferenceSetIndex)
	if err != nil {
		return 0, fmt.Errorf("failed to lookup reference set: %w", err)
	}

	if loggedSet == nil {
		return 0, fmt.Errorf("%w: set index %d for session %s, lift %s",
			ErrReferenceSetNotFound, s.ReferenceSetIndex, sessionID, params.LiftID)
	}

	// Calculate raw weight
	rawWeight := loggedSet.Weight * (s.Percentage / 100)

	// Normalize rounding parameters
	increment := NormalizeRoundingIncrement(s.RoundingIncrement)
	direction := NormalizeRoundingDirection(s.RoundingDirection)

	// Apply rounding
	roundedWeight, err := RoundWeight(rawWeight, increment, direction)
	if err != nil {
		return 0, fmt.Errorf("failed to round weight: %w", err)
	}

	return roundedWeight, nil
}

// extractSessionID extracts the session ID from the params context.
func (s *RelativeToLoadStrategy) extractSessionID(context map[string]interface{}) (string, error) {
	if context == nil {
		return "", ErrSessionIDRequired
	}

	sessionIDRaw, ok := context["sessionID"]
	if !ok {
		return "", ErrSessionIDRequired
	}

	sessionID, ok := sessionIDRaw.(string)
	if !ok || sessionID == "" {
		return "", ErrSessionIDRequired
	}

	return sessionID, nil
}

// Validate validates the strategy's configuration parameters.
func (s *RelativeToLoadStrategy) Validate() error {
	// Validate reference set index (must be non-negative)
	if s.ReferenceSetIndex < 0 {
		return fmt.Errorf("%w: got %d", ErrReferenceSetIndexInvalid, s.ReferenceSetIndex)
	}

	// Validate percentage (must be > 0)
	if s.Percentage <= 0 {
		return fmt.Errorf("%w: got %.2f", ErrRelativeToPercentageInvalid, s.Percentage)
	}

	// Validate rounding direction if specified
	if s.RoundingDirection != "" {
		if err := ValidateRoundingDirection(s.RoundingDirection); err != nil {
			return err
		}
	}

	// Validate rounding increment if specified (must be > 0 if provided)
	// A value of 0 is allowed as it means "use default"
	// Only reject explicitly negative values
	if s.RoundingIncrement < 0 {
		return fmt.Errorf("%w: rounding increment cannot be negative", ErrInvalidParams)
	}

	return nil
}

// SetSessionLookup sets the session lookup repository.
// This is used after deserialization to inject the dependency.
func (s *RelativeToLoadStrategy) SetSessionLookup(lookup SessionLookup) {
	s.sessionLookup = lookup
}

// MarshalJSON implements json.Marshaler.
// Includes the type discriminator in the JSON output.
func (s *RelativeToLoadStrategy) MarshalJSON() ([]byte, error) {
	type Alias RelativeToLoadStrategy
	return json.Marshal(&struct {
		Type LoadStrategyType `json:"type"`
		*Alias
	}{
		Type:  TypeRelativeTo,
		Alias: (*Alias)(s),
	})
}

// UnmarshalRelativeTo deserializes a RelativeToLoadStrategy from JSON.
// This is a factory function that can be registered with StrategyFactory.
func UnmarshalRelativeTo(data json.RawMessage) (LoadStrategy, error) {
	var s RelativeToLoadStrategy
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RelativeTo strategy: %w", err)
	}

	// Validate the deserialized strategy
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid RelativeTo strategy: %w", err)
	}

	return &s, nil
}

// RegisterRelativeTo registers the RelativeTo strategy with a factory.
// This is a convenience function for setting up the factory.
func RegisterRelativeTo(factory *StrategyFactory) {
	factory.Register(TypeRelativeTo, UnmarshalRelativeTo)
}

// Ensure RelativeToLoadStrategy implements LoadStrategy.
var _ LoadStrategy = (*RelativeToLoadStrategy)(nil)
