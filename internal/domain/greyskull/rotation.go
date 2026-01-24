// Package greyskull provides domain logic for the GreySkull LP program.
// This file implements the A/B week variant rotation logic for determining
// which day template to use based on week number and day position.
package greyskull

import (
	"errors"
)

// Variant represents the A/B day variant in GreySkull LP.
type Variant string

const (
	// VariantA represents the "A" day (Bench Day).
	VariantA Variant = "A"
	// VariantB represents the "B" day (OHP Day).
	VariantB Variant = "B"
)

// Errors for rotation validation.
var (
	ErrInvalidWeekNumber  = errors.New("week number must be >= 1")
	ErrInvalidDayPosition = errors.New("day position must be 1, 2, or 3")
	ErrInvalidVariant     = errors.New("variant must be 'A' or 'B'")
)

// LiftInfo represents a lift in a day template with its set/rep scheme.
type LiftInfo struct {
	// Slug is the lift identifier (e.g., "bench-press", "overhead-press").
	Slug string
	// Sets is the number of sets.
	Sets int
	// Reps is the target reps per set (negative indicates AMRAP).
	Reps int
	// IsAMRAP indicates if the final set is AMRAP.
	IsAMRAP bool
}

// GetVariantForDay returns the variant ("A" or "B") based on week number and day position.
//
// The rotation pattern is:
//
//	Week 1 (odd):  Day 1 → A, Day 2 → B, Day 3 → A
//	Week 2 (even): Day 1 → B, Day 2 → A, Day 3 → B
//	Week 3 (odd):  Day 1 → A, Day 2 → B, Day 3 → A
//	Week 4 (even): Day 1 → B, Day 2 → A, Day 3 → B
//	...and so on
//
// The logic: if weekParity == dayPositionParity then "A", else "B".
//
// Parameters:
//   - weekNumber: The week number (1-indexed, must be >= 1)
//   - dayPosition: The day position within the week (1, 2, or 3)
//
// Returns:
//   - Variant: "A" or "B"
//   - error: if weekNumber < 1 or dayPosition is not 1, 2, or 3
func GetVariantForDay(weekNumber int, dayPosition int) (Variant, error) {
	if weekNumber < 1 {
		return "", ErrInvalidWeekNumber
	}
	if dayPosition < 1 || dayPosition > 3 {
		return "", ErrInvalidDayPosition
	}

	isOddWeek := weekNumber%2 == 1
	isOddPosition := dayPosition%2 == 1

	if isOddWeek == isOddPosition {
		return VariantA, nil
	}
	return VariantB, nil
}

// GetVariantString returns the variant as a string pointer for integration
// with the Week.Variant field.
func GetVariantString(weekNumber int, dayPosition int) (*string, error) {
	variant, err := GetVariantForDay(weekNumber, dayPosition)
	if err != nil {
		return nil, err
	}
	s := string(variant)
	return &s, nil
}

// GetDayTemplate returns the list of lifts for a given variant.
//
// Variant A (Bench Day):
//   - Bench Press: 2x5 + 1x5+ (AMRAP)
//   - Barbell Row: 2x5 + 1x5+ (AMRAP)
//   - Squat: 2x5 + 1x5+ (AMRAP)
//   - Tricep Extension: 3x12 (AMRAP final)
//   - Ab Rollout: 3x10+
//
// Variant B (OHP Day):
//   - Overhead Press: 2x5 + 1x5+ (AMRAP)
//   - Chin-ups/Pull-ups: 2x5 + 1x5+ (AMRAP)
//   - Deadlift: 2x5 + 1x5+ (AMRAP) (alternates with Squat)
//   - Bicep Curl: 3x12 (AMRAP final)
//   - Shrug: 3x12 (AMRAP final)
//
// Parameters:
//   - variant: "A" or "B"
//
// Returns:
//   - []LiftInfo: The lifts for the day
//   - error: if variant is invalid
func GetDayTemplate(variant Variant) ([]LiftInfo, error) {
	switch variant {
	case VariantA:
		return []LiftInfo{
			{Slug: "bench-press", Sets: 3, Reps: 5, IsAMRAP: true},
			{Slug: "barbell-row", Sets: 3, Reps: 5, IsAMRAP: true},
			{Slug: "squat", Sets: 3, Reps: 5, IsAMRAP: true},
			{Slug: "tricep-extension", Sets: 3, Reps: 12, IsAMRAP: true},
			{Slug: "ab-rollout", Sets: 3, Reps: 10, IsAMRAP: true},
		}, nil
	case VariantB:
		return []LiftInfo{
			{Slug: "overhead-press", Sets: 3, Reps: 5, IsAMRAP: true},
			{Slug: "chin-up", Sets: 3, Reps: 5, IsAMRAP: true},
			{Slug: "deadlift", Sets: 3, Reps: 5, IsAMRAP: true},
			{Slug: "bicep-curl", Sets: 3, Reps: 12, IsAMRAP: true},
			{Slug: "shrug", Sets: 3, Reps: 12, IsAMRAP: true},
		}, nil
	default:
		return nil, ErrInvalidVariant
	}
}

// GetDayTemplateForWeek is a convenience function that combines GetVariantForDay
// and GetDayTemplate to get the lifts for a specific week and day position.
func GetDayTemplateForWeek(weekNumber int, dayPosition int) ([]LiftInfo, error) {
	variant, err := GetVariantForDay(weekNumber, dayPosition)
	if err != nil {
		return nil, err
	}
	return GetDayTemplate(variant)
}

// GetLiftSlugs returns just the lift slugs for a given variant.
// This is useful when only the exercise identifiers are needed.
func GetLiftSlugs(variant Variant) ([]string, error) {
	lifts, err := GetDayTemplate(variant)
	if err != nil {
		return nil, err
	}

	slugs := make([]string, len(lifts))
	for i, lift := range lifts {
		slugs[i] = lift.Slug
	}
	return slugs, nil
}

// GetMainLiftSlugs returns the main lift slugs (first 3) for a variant.
// Main lifts are the compound movements that use the 2x5 + 1x5+ scheme.
func GetMainLiftSlugs(variant Variant) ([]string, error) {
	lifts, err := GetDayTemplate(variant)
	if err != nil {
		return nil, err
	}

	// First 3 lifts are main lifts
	if len(lifts) < 3 {
		return nil, errors.New("template has fewer than 3 lifts")
	}

	slugs := make([]string, 3)
	for i := 0; i < 3; i++ {
		slugs[i] = lifts[i].Slug
	}
	return slugs, nil
}

// GetAccessoryLiftSlugs returns the accessory lift slugs for a variant.
// Accessory lifts are the isolation movements that use the 3x12 scheme.
func GetAccessoryLiftSlugs(variant Variant) ([]string, error) {
	lifts, err := GetDayTemplate(variant)
	if err != nil {
		return nil, err
	}

	// Lifts after the first 3 are accessories
	if len(lifts) <= 3 {
		return []string{}, nil
	}

	slugs := make([]string, len(lifts)-3)
	for i := 3; i < len(lifts); i++ {
		slugs[i-3] = lifts[i].Slug
	}
	return slugs, nil
}
