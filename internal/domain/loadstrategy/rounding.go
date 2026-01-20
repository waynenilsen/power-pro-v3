// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"errors"
	"fmt"
	"math"
)

// RoundingDirection specifies how to round calculated weights.
type RoundingDirection string

const (
	// RoundNearest rounds to the nearest increment (standard rounding).
	RoundNearest RoundingDirection = "NEAREST"
	// RoundDown always rounds down (conservative/floor).
	RoundDown RoundingDirection = "DOWN"
	// RoundUp always rounds up (ceiling).
	RoundUp RoundingDirection = "UP"
)

// ValidRoundingDirections contains all valid rounding direction values.
var ValidRoundingDirections = map[RoundingDirection]bool{
	RoundNearest: true,
	RoundDown:    true,
	RoundUp:      true,
}

// Rounding errors.
var (
	ErrNegativeWeight    = errors.New("weight cannot be negative")
	ErrInvalidIncrement  = errors.New("rounding increment must be greater than zero")
	ErrInvalidRoundingDirection = errors.New("invalid rounding direction")
)

// DefaultRoundingIncrement is the default weight increment for rounding (5.0 lbs/kg).
const DefaultRoundingIncrement = 5.0

// DefaultRoundingDirection is the default rounding direction.
const DefaultRoundingDirection = RoundNearest

// RoundWeight rounds a weight to the specified increment using the given direction.
// Parameters:
//   - weight: the weight to round (must be non-negative)
//   - increment: the rounding increment (e.g., 2.5, 5.0); must be > 0
//   - direction: how to round (NEAREST, DOWN, UP)
//
// Returns the rounded weight or an error if parameters are invalid.
func RoundWeight(weight float64, increment float64, direction RoundingDirection) (float64, error) {
	// Validate weight
	if weight < 0 {
		return 0, fmt.Errorf("%w: got %.2f", ErrNegativeWeight, weight)
	}

	// Validate increment
	if increment <= 0 {
		return 0, fmt.Errorf("%w: got %.2f", ErrInvalidIncrement, increment)
	}

	// Validate direction
	if !ValidRoundingDirections[direction] {
		return 0, fmt.Errorf("%w: %s", ErrInvalidRoundingDirection, direction)
	}

	// If weight is zero, return zero
	if weight == 0 {
		return 0, nil
	}

	// Apply rounding based on direction
	switch direction {
	case RoundNearest:
		return math.Round(weight/increment) * increment, nil
	case RoundDown:
		return math.Floor(weight/increment) * increment, nil
	case RoundUp:
		return math.Ceil(weight/increment) * increment, nil
	default:
		// This shouldn't happen due to validation above, but handle gracefully
		return 0, fmt.Errorf("%w: %s", ErrInvalidRoundingDirection, direction)
	}
}

// RoundWeightNearest is a convenience function that rounds to the nearest increment.
func RoundWeightNearest(weight float64, increment float64) (float64, error) {
	return RoundWeight(weight, increment, RoundNearest)
}

// RoundWeightDown is a convenience function that rounds down to the nearest increment.
func RoundWeightDown(weight float64, increment float64) (float64, error) {
	return RoundWeight(weight, increment, RoundDown)
}

// RoundWeightUp is a convenience function that rounds up to the nearest increment.
func RoundWeightUp(weight float64, increment float64) (float64, error) {
	return RoundWeight(weight, increment, RoundUp)
}

// ValidateRoundingDirection checks if a rounding direction is valid.
func ValidateRoundingDirection(direction RoundingDirection) error {
	if direction == "" {
		// Empty is allowed; defaults to NEAREST
		return nil
	}
	if !ValidRoundingDirections[direction] {
		return fmt.Errorf("%w: %s", ErrInvalidRoundingDirection, direction)
	}
	return nil
}

// NormalizeRoundingDirection returns the effective rounding direction,
// using the default if the provided value is empty.
func NormalizeRoundingDirection(direction RoundingDirection) RoundingDirection {
	if direction == "" {
		return DefaultRoundingDirection
	}
	return direction
}

// NormalizeRoundingIncrement returns the effective rounding increment,
// using the default if the provided value is zero or negative.
func NormalizeRoundingIncrement(increment float64) float64 {
	if increment <= 0 {
		return DefaultRoundingIncrement
	}
	return increment
}
