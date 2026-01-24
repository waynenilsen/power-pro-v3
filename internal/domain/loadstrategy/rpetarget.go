// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/waynenilsen/power-pro-v3/internal/domain/rpechart"
)

// RPETarget validation errors.
var (
	ErrTargetRepsInvalid = errors.New("target reps must be between 1 and 12")
	ErrTargetRPEInvalid  = errors.New("target RPE must be between 7.0 and 10.0")
	ErrRPEChartRequired  = errors.New("RPE chart is required for RPE_TARGET strategy")
)

// Valid RPE values (7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0)
var validRPEValues = map[float64]bool{
	7.0: true, 7.5: true, 8.0: true, 8.5: true, 9.0: true, 9.5: true, 10.0: true,
}

// RPETargetLoadStrategy calculates load based on target RPE and rep count.
// This strategy uses the RPE chart to convert (reps, RPE) to a percentage of 1RM,
// then applies that percentage to the user's 1RM.
//
// Example: 5 reps @ RPE 8 with user's 1RM of 400 lbs
//   - RPE chart lookup: 5 reps @ RPE 8 = 77% (0.77)
//   - Calculated: 400 * 0.77 = 308 lbs
//   - Rounded to nearest 5: 310 lbs
type RPETargetLoadStrategy struct {
	// TargetReps is the number of reps prescribed (1-12).
	TargetReps int `json:"targetReps"`

	// TargetRPE is the target rate of perceived exertion (7.0-10.0 in 0.5 increments).
	TargetRPE float64 `json:"targetRpe"`

	// RoundingIncrement is the weight increment for rounding (e.g., 2.5, 5.0).
	// Optional; defaults to 5.0 if not specified or <= 0.
	RoundingIncrement float64 `json:"roundingIncrement,omitempty"`

	// RoundingDirection specifies how to round (NEAREST, DOWN, UP).
	// Optional; defaults to NEAREST if not specified.
	RoundingDirection RoundingDirection `json:"roundingDirection,omitempty"`

	// maxLookup is the repository for looking up user maxes.
	// This is injected and not serialized.
	maxLookup MaxLookup `json:"-"`

	// rpeChart is the RPE chart for converting (reps, RPE) to percentage.
	// This is injected and not serialized.
	rpeChart *rpechart.RPEChart `json:"-"`
}

// NewRPETargetLoadStrategy creates a new RPETargetLoadStrategy with the given parameters.
func NewRPETargetLoadStrategy(
	targetReps int,
	targetRPE float64,
	roundingIncrement float64,
	roundingDirection RoundingDirection,
	maxLookup MaxLookup,
	rpeChart *rpechart.RPEChart,
) *RPETargetLoadStrategy {
	return &RPETargetLoadStrategy{
		TargetReps:        targetReps,
		TargetRPE:         targetRPE,
		RoundingIncrement: roundingIncrement,
		RoundingDirection: roundingDirection,
		maxLookup:         maxLookup,
		rpeChart:          rpeChart,
	}
}

// Type returns the strategy type discriminator.
func (s *RPETargetLoadStrategy) Type() LoadStrategyType {
	return TypeRPETarget
}

// CalculateLoad calculates the target weight based on RPE chart lookup.
// It uses the user's 1RM and the RPE chart to determine the appropriate weight.
//
// The calculation follows these steps:
//  1. Validate parameters
//  2. Get user's 1RM for the lift
//  3. Get percentage from RPE chart (using params.LookupContext if available, otherwise injected chart)
//  4. Calculate: weight = 1RM * percentage
//  5. Round to configured increment
func (s *RPETargetLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
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

	// Determine which RPE chart to use: prefer lookup context, fall back to injected
	var chart *rpechart.RPEChart
	if params.LookupContext != nil && params.LookupContext.HasRPEChart() {
		chart = params.LookupContext.RPEChart
	} else if s.rpeChart != nil {
		chart = s.rpeChart
	} else {
		return 0, fmt.Errorf("%w: no RPE chart available", ErrRPEChartRequired)
	}

	// Fetch the user's 1RM (RPE-based strategies always use 1RM as the reference)
	maxValue, err := s.maxLookup.GetCurrentMax(ctx, params.UserID, params.LiftID, string(ReferenceOneRM))
	if err != nil {
		return 0, fmt.Errorf("failed to lookup max: %w", err)
	}

	if maxValue == nil {
		return 0, fmt.Errorf("%w: no %s found for user %s, lift %s",
			ErrMaxNotFound, ReferenceOneRM, params.UserID, params.LiftID)
	}

	// Look up percentage from RPE chart
	percentage, err := chart.GetPercentage(s.TargetReps, s.TargetRPE)
	if err != nil {
		return 0, fmt.Errorf("RPE chart lookup failed: %w", err)
	}

	// Calculate raw weight (percentage is already a decimal, e.g., 0.77 for 77%)
	rawWeight := maxValue.Value * percentage

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
func (s *RPETargetLoadStrategy) Validate() error {
	// Validate target reps (1-12)
	if s.TargetReps < 1 || s.TargetReps > 12 {
		return fmt.Errorf("%w: got %d", ErrTargetRepsInvalid, s.TargetReps)
	}

	// Validate target RPE (7.0-10.0 in 0.5 increments)
	if !validRPEValues[s.TargetRPE] {
		return fmt.Errorf("%w: got %.1f", ErrTargetRPEInvalid, s.TargetRPE)
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

// SetMaxLookup sets the max lookup repository.
// This is used after deserialization to inject the dependency.
func (s *RPETargetLoadStrategy) SetMaxLookup(maxLookup MaxLookup) {
	s.maxLookup = maxLookup
}

// SetRPEChart sets the RPE chart.
// This is used after deserialization to inject the dependency.
func (s *RPETargetLoadStrategy) SetRPEChart(chart *rpechart.RPEChart) {
	s.rpeChart = chart
}

// MarshalJSON implements json.Marshaler.
// Includes the type discriminator in the JSON output.
func (s *RPETargetLoadStrategy) MarshalJSON() ([]byte, error) {
	type Alias RPETargetLoadStrategy
	return json.Marshal(&struct {
		Type LoadStrategyType `json:"type"`
		*Alias
	}{
		Type:  TypeRPETarget,
		Alias: (*Alias)(s),
	})
}

// UnmarshalRPETarget deserializes an RPETargetLoadStrategy from JSON.
// This is a factory function that can be registered with StrategyFactory.
func UnmarshalRPETarget(data json.RawMessage) (LoadStrategy, error) {
	var s RPETargetLoadStrategy
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RPETarget strategy: %w", err)
	}

	// Validate the deserialized strategy
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid RPETarget strategy: %w", err)
	}

	return &s, nil
}

// RegisterRPETarget registers the RPETarget strategy with a factory.
// This is a convenience function for setting up the factory.
func RegisterRPETarget(factory *StrategyFactory) {
	factory.Register(TypeRPETarget, UnmarshalRPETarget)
}

// Ensure RPETargetLoadStrategy implements LoadStrategy.
var _ LoadStrategy = (*RPETargetLoadStrategy)(nil)
