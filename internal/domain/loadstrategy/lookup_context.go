// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"github.com/waynenilsen/power-pro-v3/internal/domain/dailylookup"
	"github.com/waynenilsen/power-pro-v3/internal/domain/rotationlookup"
	"github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"
)

// LookupContext provides the context needed for lookup-based load modifications.
// This context is passed through the resolution chain to enable lookups to
// modify base percentages during prescription resolution.
type LookupContext struct {
	// WeekNumber is the current week number within the cycle (1-indexed).
	// Used for WeeklyLookup resolution.
	WeekNumber int

	// DaySlug is the day identifier (e.g., "heavy", "light", "day-a").
	// Used for DailyLookup resolution. Case-insensitive matching is applied.
	DaySlug string

	// SetNumber is the current set number (1-indexed) for set-specific percentages.
	// When a lookup entry has a percentages array, this determines which percentage to use.
	SetNumber int

	// WeeklyLookup is the weekly lookup table to use for week-based modifications.
	// Optional: if nil, no weekly lookup modifications are applied.
	WeeklyLookup *weeklylookup.WeeklyLookup

	// DailyLookup is the daily lookup table to use for day-based modifications.
	// Optional: if nil, no daily lookup modifications are applied.
	DailyLookup *dailylookup.DailyLookup

	// RotationPosition is the current position in the rotation (0-indexed).
	// Used for RotationLookup resolution to determine which lift is in focus.
	RotationPosition int

	// RotationLookup is the rotation lookup table to use for lift rotation patterns.
	// Optional: if nil, no rotation-based modifications are applied.
	// Programs like Conjugate/Westside use this to cycle through different lifts.
	RotationLookup *rotationlookup.RotationLookup
}

// HasWeeklyLookup returns true if a weekly lookup is configured and the week number is valid.
func (c *LookupContext) HasWeeklyLookup() bool {
	return c != nil && c.WeeklyLookup != nil && c.WeekNumber > 0
}

// HasDailyLookup returns true if a daily lookup is configured and the day slug is non-empty.
func (c *LookupContext) HasDailyLookup() bool {
	return c != nil && c.DailyLookup != nil && c.DaySlug != ""
}

// GetWeeklyEntry returns the weekly lookup entry for the current week number.
// Returns nil if no weekly lookup is configured or the week number is not found.
func (c *LookupContext) GetWeeklyEntry() *weeklylookup.WeeklyLookupEntry {
	if !c.HasWeeklyLookup() {
		return nil
	}
	return c.WeeklyLookup.GetByWeekNumber(c.WeekNumber)
}

// GetDailyEntry returns the daily lookup entry for the current day slug.
// Returns nil if no daily lookup is configured or the day slug is not found.
func (c *LookupContext) GetDailyEntry() *dailylookup.DailyLookupEntry {
	if !c.HasDailyLookup() {
		return nil
	}
	return c.DailyLookup.GetByDayIdentifier(c.DaySlug)
}

// ApplyModifiers applies lookup modifiers to a base percentage, enabling programs
// to vary intensity across weeks and days without changing the underlying prescription.
//
// This is a core periodization mechanism used by programs like:
//   - 5/3/1: Different percentages each week (65/75/85% → 70/80/90% → 75/85/95%)
//   - Texas Method: Different intensities by day (heavy/light/medium)
//   - Juggernaut: Different percentages per set within a week
//
// The modification order is:
//  1. Weekly lookup (set-specific percentage or percentage modifier)
//  2. Daily lookup (percentage modifier)
//
// Two distinct modification modes exist for weekly lookups:
//   - Set-specific percentages: Used when each set has a different prescribed percentage
//     (e.g., 5/3/1's "5x65%, 5x75%, 5+x85%"). These REPLACE the base percentage entirely.
//   - Percentage modifier: A multiplier applied to the base (e.g., deload week at 90%
//     would use modifier=90, reducing all weights by 10%).
//
// Daily modifiers are always multiplicative, allowing patterns like:
//   - Monday (Heavy): 100% modifier → full intensity
//   - Wednesday (Light): 80% modifier → reduced intensity
//   - Friday (Medium): 90% modifier → moderate intensity
//
// When both weekly and daily modifiers apply, they stack multiplicatively.
// Example: Base 85%, weekly modifier 95%, daily modifier 90% → 85 * 0.95 * 0.90 = 72.675%
func (c *LookupContext) ApplyModifiers(basePercentage float64) float64 {
	if c == nil {
		return basePercentage
	}

	resultPercentage := basePercentage

	// Weekly lookup modifications are applied first.
	// Programs define weekly intensity patterns through either:
	// 1. Explicit per-set percentages (replaces base entirely) - for programs like 5/3/1
	// 2. A percentage modifier (multiplies base) - for deload weeks or wave loading
	weeklyEntry := c.GetWeeklyEntry()
	if weeklyEntry != nil {
		if len(weeklyEntry.Percentages) > 0 && c.SetNumber > 0 {
			// Set-specific percentage mode: each set has an explicitly defined percentage.
			// SetNumber is 1-indexed (user-facing), array is 0-indexed, so we subtract 1.
			// This completely replaces the base percentage rather than modifying it.
			setIndex := c.SetNumber - 1
			if setIndex < len(weeklyEntry.Percentages) {
				resultPercentage = weeklyEntry.Percentages[setIndex]
			}
		} else if weeklyEntry.PercentageModifier != nil {
			// Percentage modifier mode: scale the base percentage.
			// Modifier is stored as a percentage (e.g., 90 means 90%), so divide by 100.
			resultPercentage = resultPercentage * (*weeklyEntry.PercentageModifier / 100)
		}
	}

	// Daily lookup modifications are applied second (multiplicatively).
	// This enables heavy/light/medium day patterns within a training week.
	dailyEntry := c.GetDailyEntry()
	if dailyEntry != nil {
		if dailyEntry.PercentageModifier != 0 {
			resultPercentage = resultPercentage * (dailyEntry.PercentageModifier / 100)
		}
	}

	return resultPercentage
}

// GetRepsForSet returns the target reps for a specific set from the weekly lookup.
// Returns -1 if no weekly lookup is configured, set number is invalid,
// or no reps are defined for the set.
func (c *LookupContext) GetRepsForSet() int {
	if c == nil || c.SetNumber <= 0 {
		return -1
	}

	weeklyEntry := c.GetWeeklyEntry()
	if weeklyEntry == nil {
		return -1
	}

	setIndex := c.SetNumber - 1
	if setIndex >= len(weeklyEntry.Reps) {
		return -1
	}

	return weeklyEntry.Reps[setIndex]
}

// HasRotationLookup returns true if a rotation lookup is configured.
// Note: RotationPosition of 0 is valid (it's the first position in the rotation).
func (c *LookupContext) HasRotationLookup() bool {
	return c != nil && c.RotationLookup != nil
}

// GetRotationEntry returns the rotation lookup entry for the current rotation position.
// Returns nil if no rotation lookup is configured or the position is not found.
func (c *LookupContext) GetRotationEntry() *rotationlookup.RotationLookupEntry {
	if !c.HasRotationLookup() {
		return nil
	}
	return c.RotationLookup.GetByPosition(c.RotationPosition)
}

// IsLiftInRotationFocus checks if the given lift identifier is the current focus
// based on the rotation position. This is used by programs like Conjugate/Westside
// where different lifts are emphasized on different training days/weeks.
//
// Returns true if:
//   - A rotation lookup is configured AND
//   - The current rotation position's entry has a matching lift identifier
//
// Returns false if:
//   - No rotation lookup is configured
//   - The rotation position doesn't exist in the lookup
//   - The lift identifier doesn't match the current rotation entry
func (c *LookupContext) IsLiftInRotationFocus(liftIdentifier string) bool {
	entry := c.GetRotationEntry()
	if entry == nil {
		return false
	}
	return entry.LiftIdentifier == liftIdentifier
}
