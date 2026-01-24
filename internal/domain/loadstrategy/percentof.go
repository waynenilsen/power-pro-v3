// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// ReferenceType identifies the type of max used as a reference for percentage calculations.
type ReferenceType string

const (
	// ReferenceOneRM uses the lifter's 1 rep max as the reference.
	ReferenceOneRM ReferenceType = "ONE_RM"
	// ReferenceTrainingMax uses the lifter's Training Max as the reference.
	ReferenceTrainingMax ReferenceType = "TRAINING_MAX"
	// ReferenceE1RM uses the lifter's estimated 1 rep max as the reference.
	ReferenceE1RM ReferenceType = "E1RM"
)

// Alias constants for backwards compatibility and convenience.
const (
	OneRM       = ReferenceOneRM
	TrainingMax = ReferenceTrainingMax
	E1RM        = ReferenceE1RM
)

// ValidReferenceTypes contains all valid reference type values.
var ValidReferenceTypes = map[ReferenceType]bool{
	ReferenceOneRM:       true,
	ReferenceTrainingMax: true,
	ReferenceE1RM:        true,
}

// PercentOf validation errors.
var (
	ErrPercentageRequired    = errors.New("percentage is required")
	ErrPercentageNotPositive = errors.New("percentage must be greater than 0")
	ErrReferenceTypeRequired = errors.New("reference type is required")
	ErrReferenceTypeInvalid  = errors.New("reference type must be ONE_RM, TRAINING_MAX, or E1RM")
)

// PercentOfLoadStrategy calculates load as a percentage of a reference max.
// This is the most common load calculation method in powerlifting programs.
//
// Example: 85% of Training Max with 5lb rounding
//   - User's TM for squat: 315 lbs
//   - Calculated: 315 * 0.85 = 267.75
//   - Rounded to nearest 5: 270 lbs
type PercentOfLoadStrategy struct {
	// ReferenceType specifies which max to use (ONE_RM or TRAINING_MAX).
	ReferenceType ReferenceType `json:"referenceType"`

	// Percentage is the percentage of the reference max (e.g., 85 for 85%).
	// Must be > 0. Values > 100 are allowed for overload work.
	Percentage float64 `json:"percentage"`

	// RoundingIncrement is the weight increment for rounding (e.g., 2.5, 5.0).
	// Optional; defaults to 5.0 if not specified or <= 0.
	RoundingIncrement float64 `json:"roundingIncrement,omitempty"`

	// RoundingDirection specifies how to round (NEAREST, DOWN, UP).
	// Optional; defaults to NEAREST if not specified.
	RoundingDirection RoundingDirection `json:"roundingDirection,omitempty"`

	// maxLookup is the repository for looking up user maxes.
	// This is injected and not serialized.
	maxLookup MaxLookup `json:"-"`
}

// NewPercentOfLoadStrategy creates a new PercentOfLoadStrategy with the given parameters.
// The maxLookup parameter is used to fetch the user's max for calculations.
func NewPercentOfLoadStrategy(
	referenceType ReferenceType,
	percentage float64,
	roundingIncrement float64,
	roundingDirection RoundingDirection,
	maxLookup MaxLookup,
) *PercentOfLoadStrategy {
	return &PercentOfLoadStrategy{
		ReferenceType:     referenceType,
		Percentage:        percentage,
		RoundingIncrement: roundingIncrement,
		RoundingDirection: roundingDirection,
		maxLookup:         maxLookup,
	}
}

// Type returns the strategy type discriminator.
func (s *PercentOfLoadStrategy) Type() LoadStrategyType {
	return TypePercentOf
}

// CalculateLoad calculates the target weight based on a percentage of the user's max.
// It fetches the current max for the specified user/lift/reference type,
// applies the percentage (with optional lookup modifiers), and rounds to the configured increment.
//
// If params.LookupContext is provided, lookup modifiers are applied to the base percentage:
//   - Weekly lookup: set-specific percentages or percentage modifier
//   - Daily lookup: percentage modifier
func (s *PercentOfLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
	// Validate params
	if err := params.Validate(); err != nil {
		return 0, err
	}

	// Validate strategy configuration
	if err := s.Validate(); err != nil {
		return 0, err
	}

	// Check that maxLookup is available
	if s.maxLookup == nil {
		return 0, fmt.Errorf("%w: max lookup not configured", ErrInvalidParams)
	}

	// Map reference type to the string expected by the max repository
	maxType := string(s.ReferenceType)

	// Fetch the current max
	maxValue, err := s.maxLookup.GetCurrentMax(ctx, params.UserID, params.LiftID, maxType)
	if err != nil {
		return 0, fmt.Errorf("failed to lookup max: %w", err)
	}

	if maxValue == nil {
		return 0, fmt.Errorf("%w: no %s found for user %s, lift %s",
			ErrMaxNotFound, s.ReferenceType, params.UserID, params.LiftID)
	}

	// Start with the base percentage from the strategy
	effectivePercentage := s.Percentage

	// Apply lookup modifiers if lookup context is provided
	if params.LookupContext != nil {
		effectivePercentage = params.LookupContext.ApplyModifiers(s.Percentage)
	}

	// Calculate the raw weight using the effective percentage
	rawWeight := maxValue.Value * (effectivePercentage / 100)

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

// Validate validates the strategy's configuration parameters.
func (s *PercentOfLoadStrategy) Validate() error {
	// Validate reference type
	if s.ReferenceType == "" {
		return ErrReferenceTypeRequired
	}
	if !ValidReferenceTypes[s.ReferenceType] {
		return fmt.Errorf("%w: got %s", ErrReferenceTypeInvalid, s.ReferenceType)
	}

	// Validate percentage
	if s.Percentage <= 0 {
		return ErrPercentageNotPositive
	}
	// Note: Percentage > 100 is explicitly allowed for overload work

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

// SetMaxLookup sets the max lookup repository.
// This is used after deserialization to inject the dependency.
func (s *PercentOfLoadStrategy) SetMaxLookup(maxLookup MaxLookup) {
	s.maxLookup = maxLookup
}

// MarshalJSON implements json.Marshaler.
// Includes the type discriminator in the JSON output.
func (s *PercentOfLoadStrategy) MarshalJSON() ([]byte, error) {
	type Alias PercentOfLoadStrategy
	return json.Marshal(&struct {
		Type LoadStrategyType `json:"type"`
		*Alias
	}{
		Type:  TypePercentOf,
		Alias: (*Alias)(s),
	})
}

// UnmarshalPercentOf deserializes a PercentOfLoadStrategy from JSON.
// This is a factory function that can be registered with StrategyFactory.
func UnmarshalPercentOf(data json.RawMessage) (LoadStrategy, error) {
	var s PercentOfLoadStrategy
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PercentOf strategy: %w", err)
	}

	// Validate the deserialized strategy
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid PercentOf strategy: %w", err)
	}

	return &s, nil
}

// RegisterPercentOf registers the PercentOf strategy with a factory.
// This is a convenience function for setting up the factory.
func RegisterPercentOf(factory *StrategyFactory) {
	factory.Register(TypePercentOf, UnmarshalPercentOf)
}

// Ensure PercentOfLoadStrategy implements LoadStrategy.
var _ LoadStrategy = (*PercentOfLoadStrategy)(nil)
