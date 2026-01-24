// Package juggernaut provides domain logic for the Inverted Juggernaut program.
package juggernaut

// RepStandards maps wave index (0-3) to expected rep standard at AMRAP percentage.
// Wave 0 (10s): 10 reps at 75%
// Wave 1 (8s): 8 reps at 80%
// Wave 2 (5s): 5 reps at 85%
// Wave 3 (3s): 3 reps at 90%
var RepStandards = []int{10, 8, 5, 3}

// UpperBodyIncrement is the weight added per rep over standard for bench/OHP.
const UpperBodyIncrement = 2.5

// LowerBodyIncrement is the weight added per rep over standard for squat/deadlift.
const LowerBodyIncrement = 5.0

// CycleIncrementUpper is the base TM increase for upper body lifts after completing a full 16-week cycle.
const CycleIncrementUpper = 5.0

// CycleIncrementLower is the base TM increase for lower body lifts after completing a full 16-week cycle.
const CycleIncrementLower = 10.0

// CalculateNewTM calculates the new Training Max based on AMRAP performance.
//
// The formula is:
//
//	New TM = Current TM + ((AMRAP Reps - Rep Standard) × Weight Increment)
//
// Parameters:
//   - currentTM: Current training max
//   - waveIndex: 0-3 (determines rep standard: 10/8/5/3)
//   - amrapReps: Actual reps performed on realization AMRAP
//   - isUpperBody: true for bench/OHP (2.5 increment), false for squat/deadlift (5 increment)
//
// Returns the new training max. The result can be lower than currentTM if
// amrapReps < rep standard (underperformance).
func CalculateNewTM(currentTM float64, waveIndex int, amrapReps int, isUpperBody bool) float64 {
	if waveIndex < 0 || waveIndex > 3 {
		return currentTM
	}

	repStandard := RepStandards[waveIndex]
	excessReps := amrapReps - repStandard

	var increment float64
	if isUpperBody {
		increment = UpperBodyIncrement
	} else {
		increment = LowerBodyIncrement
	}

	delta := float64(excessReps) * increment
	return currentTM + delta
}

// CalculateCycleIncrement returns the base TM increase after completing a full 16-week cycle.
//
// Parameters:
//   - isUpperBody: true for bench/OHP (5 units), false for squat/deadlift (10 units)
//
// Returns the increment to add to the base TM for the next cycle.
func CalculateCycleIncrement(isUpperBody bool) float64 {
	if isUpperBody {
		return CycleIncrementUpper
	}
	return CycleIncrementLower
}
