// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// FindRM validation errors.
var (
	ErrFindRMTargetRepsInvalid = errors.New("target reps must be between 1 and 12")
)

// FindRMLoadStrategy represents a strategy where the user works up to discover
// their rep max. No weight is prescribed - the user decides the weight and
// the system records what they achieve.
//
// This is commonly used in programs like GZCL Jacked & Tan 2.0 where users
// find their 10RM in week 1, 8RM in week 2, etc.
//
// Example usage:
//   - Week 1: "Find your 10RM" - user works up until they can only do 10 reps
//   - Week 2: "Find your 8RM" - user works up until they can only do 8 reps
//
// The weight discovered can then be used by RelativeTo for back-off sets.
type FindRMLoadStrategy struct {
	// TargetReps is the number of reps for the rep max to find (1-12).
	// For example, 10 means "find your 10RM".
	TargetReps int `json:"targetReps"`
}

// NewFindRMLoadStrategy creates a new FindRMLoadStrategy with the given target reps.
func NewFindRMLoadStrategy(targetReps int) *FindRMLoadStrategy {
	return &FindRMLoadStrategy{
		TargetReps: targetReps,
	}
}

// Type returns the strategy type discriminator.
func (s *FindRMLoadStrategy) Type() LoadStrategyType {
	return TypeFindRM
}

// CalculateLoad returns 0 for FindRM strategies because no weight is prescribed.
// The user decides the weight and works up to their rep max.
//
// A return value of 0 indicates "user decides" - the prescription display
// should show "Find 10RM" (or similar) rather than "X lbs × 10".
func (s *FindRMLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
	// Validate params (still required for context)
	if err := params.Validate(); err != nil {
		return 0, err
	}

	// Validate strategy configuration
	if err := s.Validate(); err != nil {
		return 0, err
	}

	// Return 0 to indicate no prescribed weight
	return 0, nil
}

// Validate validates the strategy's configuration parameters.
func (s *FindRMLoadStrategy) Validate() error {
	// Validate target reps (1-12, matching RPE chart range)
	if s.TargetReps < 1 || s.TargetReps > 12 {
		return fmt.Errorf("%w: got %d", ErrFindRMTargetRepsInvalid, s.TargetReps)
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
// Includes the type discriminator in the JSON output.
func (s *FindRMLoadStrategy) MarshalJSON() ([]byte, error) {
	type Alias FindRMLoadStrategy
	return json.Marshal(&struct {
		Type LoadStrategyType `json:"type"`
		*Alias
	}{
		Type:  TypeFindRM,
		Alias: (*Alias)(s),
	})
}

// UnmarshalFindRM deserializes a FindRMLoadStrategy from JSON.
// This is a factory function that can be registered with StrategyFactory.
func UnmarshalFindRM(data json.RawMessage) (LoadStrategy, error) {
	var s FindRMLoadStrategy
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FindRM strategy: %w", err)
	}

	// Validate the deserialized strategy
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid FindRM strategy: %w", err)
	}

	return &s, nil
}

// RegisterFindRM registers the FindRM strategy with a factory.
// This is a convenience function for setting up the factory.
func RegisterFindRM(factory *StrategyFactory) {
	factory.Register(TypeFindRM, UnmarshalFindRM)
}

// Ensure FindRMLoadStrategy implements LoadStrategy.
var _ LoadStrategy = (*FindRMLoadStrategy)(nil)
