// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"github.com/waynenilsen/power-pro-v3/internal/domain/dailylookup"
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

// ApplyModifiers applies lookup modifiers to a base percentage.
// The modification order is:
//  1. Weekly lookup (set-specific percentage or percentage modifier)
//  2. Daily lookup (percentage modifier)
//
// If both lookups have percentage modifiers, they are applied multiplicatively.
// If a weekly lookup has a percentages array and SetNumber is valid, that specific
// percentage is used directly (ignoring the base percentage).
//
// Returns the modified percentage.
func (c *LookupContext) ApplyModifiers(basePercentage float64) float64 {
	if c == nil {
		return basePercentage
	}

	resultPercentage := basePercentage

	// Apply weekly lookup modifications first
	weeklyEntry := c.GetWeeklyEntry()
	if weeklyEntry != nil {
		// Check if we have set-specific percentages
		if len(weeklyEntry.Percentages) > 0 && c.SetNumber > 0 {
			// Use set-specific percentage (1-indexed, so subtract 1)
			setIndex := c.SetNumber - 1
			if setIndex < len(weeklyEntry.Percentages) {
				// Set-specific percentage overrides the base percentage entirely
				resultPercentage = weeklyEntry.Percentages[setIndex]
			}
		} else if weeklyEntry.PercentageModifier != nil {
			// Apply percentage modifier as a multiplier
			resultPercentage = resultPercentage * (*weeklyEntry.PercentageModifier / 100)
		}
	}

	// Apply daily lookup modifications second
	dailyEntry := c.GetDailyEntry()
	if dailyEntry != nil {
		// Daily lookup's PercentageModifier is always a multiplier
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
