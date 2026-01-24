// Package e1rm provides domain logic for Estimated 1-Rep Maximum (E1RM) calculations.
// E1RM is calculated from a performed set using RPE (Rate of Perceived Exertion) data.
package e1rm

import (
	"errors"
	"fmt"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/rpechart"
)

// Validation errors
var (
	ErrWeightMustBePositive = errors.New("weight must be greater than 0")
	ErrRepsOutOfRange       = errors.New("reps must be between 1 and 12")
	ErrRPEOutOfRange        = errors.New("RPE must be between 7.0 and 10.0")
)

// E1RMRoundingIncrement is the standard rounding increment for E1RM calculations (2.5 lbs/kg).
const E1RMRoundingIncrement = 2.5

// Calculator estimates 1RM from performed sets using an RPE chart.
type Calculator struct {
	rpeChart *rpechart.RPEChart
}

// NewCalculator creates a new E1RM Calculator with the given RPE chart.
func NewCalculator(chart *rpechart.RPEChart) *Calculator {
	return &Calculator{rpeChart: chart}
}

// Calculate estimates the 1RM from a performed set.
// Formula: E1RM = Weight / RPEChart.GetPercentage(RepsPerformed, RPE)
//
// Parameters:
//   - weight: the weight lifted (must be > 0)
//   - reps: the number of repetitions performed (must be 1-12)
//   - rpe: the Rate of Perceived Exertion (must be 7.0-10.0)
//
// Returns the estimated 1RM rounded to 2.5 lb increments, or an error if inputs are invalid.
func (c *Calculator) Calculate(weight float64, reps int, rpe float64) (float64, error) {
	// Validate weight
	if weight <= 0 {
		return 0, fmt.Errorf("%w: got %.2f", ErrWeightMustBePositive, weight)
	}

	// Validate reps
	if reps < 1 || reps > 12 {
		return 0, fmt.Errorf("%w: got %d", ErrRepsOutOfRange, reps)
	}

	// Validate RPE
	if rpe < 7.0 || rpe > 10.0 {
		return 0, fmt.Errorf("%w: got %.1f", ErrRPEOutOfRange, rpe)
	}

	// Look up the percentage from the RPE chart
	percentage, err := c.rpeChart.GetPercentage(reps, rpe)
	if err != nil {
		return 0, fmt.Errorf("RPE chart lookup failed: %w", err)
	}

	// Calculate E1RM: weight / percentage
	e1rm := weight / percentage

	// Round to 2.5 lb increments
	rounded, err := loadstrategy.RoundWeightNearest(e1rm, E1RMRoundingIncrement)
	if err != nil {
		return 0, fmt.Errorf("rounding failed: %w", err)
	}

	return rounded, nil
}
