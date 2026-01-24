// Package integration provides integration tests for cross-component behavior.
// This file tests the GreySkull LP system including A/B variant rotation,
// composite set scheme generation, and AMRAP-based progression.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/greyskull"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
)

// TestGreySkullWeekVariantRotation tests the A/B variant rotation pattern.
// Pattern: if weekParity == dayPositionParity then "A", else "B".
func TestGreySkullWeekVariantRotation(t *testing.T) {
	testCases := []struct {
		name        string
		weekNumber  int
		dayPosition int
		expected    greyskull.Variant
	}{
		// Week 1 (odd): A, B, A
		{"Week 1 Day 1", 1, 1, greyskull.VariantA},
		{"Week 1 Day 2", 1, 2, greyskull.VariantB},
		{"Week 1 Day 3", 1, 3, greyskull.VariantA},
		// Week 2 (even): B, A, B
		{"Week 2 Day 1", 2, 1, greyskull.VariantB},
		{"Week 2 Day 2", 2, 2, greyskull.VariantA},
		{"Week 2 Day 3", 2, 3, greyskull.VariantB},
		// Week 3 (odd): A, B, A
		{"Week 3 Day 1", 3, 1, greyskull.VariantA},
		{"Week 3 Day 2", 3, 2, greyskull.VariantB},
		{"Week 3 Day 3", 3, 3, greyskull.VariantA},
		// Week 4 (even): B, A, B
		{"Week 4 Day 1", 4, 1, greyskull.VariantB},
		{"Week 4 Day 2", 4, 2, greyskull.VariantA},
		{"Week 4 Day 3", 4, 3, greyskull.VariantB},
		// Week 5 (odd): A, B, A
		{"Week 5 Day 1", 5, 1, greyskull.VariantA},
		{"Week 5 Day 2", 5, 2, greyskull.VariantB},
		{"Week 5 Day 3", 5, 3, greyskull.VariantA},
		// Week 6 (even): B, A, B
		{"Week 6 Day 1", 6, 1, greyskull.VariantB},
		{"Week 6 Day 2", 6, 2, greyskull.VariantA},
		{"Week 6 Day 3", 6, 3, greyskull.VariantB},
		// Week 7 (odd): A, B, A
		{"Week 7 Day 1", 7, 1, greyskull.VariantA},
		{"Week 7 Day 2", 7, 2, greyskull.VariantB},
		{"Week 7 Day 3", 7, 3, greyskull.VariantA},
		// Week 8 (even): B, A, B
		{"Week 8 Day 1", 8, 1, greyskull.VariantB},
		{"Week 8 Day 2", 8, 2, greyskull.VariantA},
		{"Week 8 Day 3", 8, 3, greyskull.VariantB},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			variant, err := greyskull.GetVariantForDay(tc.weekNumber, tc.dayPosition)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if variant != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, variant)
			}
		})
	}
}

// TestGreySkullVariantAlternation verifies the A/B pattern alternates correctly
// when viewing training sessions chronologically.
func TestGreySkullVariantAlternation(t *testing.T) {
	// Simulate 4 weeks of training (12 sessions)
	// Sessions should alternate: A, B, A, B, A, B...
	expectedPattern := []greyskull.Variant{
		// Week 1
		greyskull.VariantA, greyskull.VariantB, greyskull.VariantA,
		// Week 2
		greyskull.VariantB, greyskull.VariantA, greyskull.VariantB,
		// Week 3
		greyskull.VariantA, greyskull.VariantB, greyskull.VariantA,
		// Week 4
		greyskull.VariantB, greyskull.VariantA, greyskull.VariantB,
	}

	sessionIdx := 0
	for week := 1; week <= 4; week++ {
		for day := 1; day <= 3; day++ {
			variant, err := greyskull.GetVariantForDay(week, day)
			if err != nil {
				t.Fatalf("Week %d Day %d: unexpected error: %v", week, day, err)
			}
			if variant != expectedPattern[sessionIdx] {
				t.Errorf("Session %d (Week %d Day %d): expected %s, got %s",
					sessionIdx+1, week, day, expectedPattern[sessionIdx], variant)
			}
			sessionIdx++
		}
	}
}

// TestGreySkullSetSchemeGeneration tests the composite set scheme (Fixed + AMRAP).
func TestGreySkullSetSchemeGeneration(t *testing.T) {
	testCases := []struct {
		name         string
		fixedSets    int
		fixedReps    int
		amrapSets    int
		minAmrapReps int
		weight       float64
		expectedSets []setscheme.GeneratedSet
	}{
		{
			name:         "main lift 2x5+1x5+",
			fixedSets:    2,
			fixedReps:    5,
			amrapSets:    1,
			minAmrapReps: 5,
			weight:       135.0,
			expectedSets: []setscheme.GeneratedSet{
				{SetNumber: 1, Weight: 135.0, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 2, Weight: 135.0, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 3, Weight: 135.0, TargetReps: 5, IsWorkSet: true}, // AMRAP
			},
		},
		{
			name:         "accessory 2x12+1x12+",
			fixedSets:    2,
			fixedReps:    12,
			amrapSets:    1,
			minAmrapReps: 12,
			weight:       50.0,
			expectedSets: []setscheme.GeneratedSet{
				{SetNumber: 1, Weight: 50.0, TargetReps: 12, IsWorkSet: true},
				{SetNumber: 2, Weight: 50.0, TargetReps: 12, IsWorkSet: true},
				{SetNumber: 3, Weight: 50.0, TargetReps: 12, IsWorkSet: true}, // AMRAP
			},
		},
		{
			name:         "heavier weight 2x5+1x5+",
			fixedSets:    2,
			fixedReps:    5,
			amrapSets:    1,
			minAmrapReps: 5,
			weight:       225.0,
			expectedSets: []setscheme.GeneratedSet{
				{SetNumber: 1, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 2, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 3, Weight: 225.0, TargetReps: 5, IsWorkSet: true}, // AMRAP
			},
		},
		{
			name:         "no fixed sets (pure AMRAP)",
			fixedSets:    0,
			fixedReps:    0,
			amrapSets:    1,
			minAmrapReps: 5,
			weight:       135.0,
			expectedSets: []setscheme.GeneratedSet{
				{SetNumber: 1, Weight: 135.0, TargetReps: 5, IsWorkSet: true}, // AMRAP only
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scheme, err := setscheme.NewGreySkullSetScheme(tc.fixedSets, tc.fixedReps, tc.amrapSets, tc.minAmrapReps)
			if err != nil {
				t.Fatalf("failed to create scheme: %v", err)
			}

			sets, err := scheme.GenerateSets(tc.weight, setscheme.DefaultSetGenerationContext())
			if err != nil {
				t.Fatalf("failed to generate sets: %v", err)
			}

			if len(sets) != len(tc.expectedSets) {
				t.Fatalf("expected %d sets, got %d", len(tc.expectedSets), len(sets))
			}

			for i, expected := range tc.expectedSets {
				if sets[i].SetNumber != expected.SetNumber {
					t.Errorf("set %d: SetNumber expected %d, got %d", i+1, expected.SetNumber, sets[i].SetNumber)
				}
				if sets[i].Weight != expected.Weight {
					t.Errorf("set %d: Weight expected %.1f, got %.1f", i+1, expected.Weight, sets[i].Weight)
				}
				if sets[i].TargetReps != expected.TargetReps {
					t.Errorf("set %d: TargetReps expected %d, got %d", i+1, expected.TargetReps, sets[i].TargetReps)
				}
				if sets[i].IsWorkSet != expected.IsWorkSet {
					t.Errorf("set %d: IsWorkSet expected %v, got %v", i+1, expected.IsWorkSet, sets[i].IsWorkSet)
				}
			}
		})
	}
}

// TestGreySkullProgression tests the three-tier progression logic.
func TestGreySkullProgression(t *testing.T) {
	testCases := []struct {
		name            string
		currentWeight   float64
		repsPerformed   int
		weightIncrement float64
		minReps         int
		doubleThreshold int
		deloadPercent   float64
		expectedWeight  float64
		expectApplied   bool
	}{
		// Main lift scenarios (minReps=5, doubleThreshold=10)
		{
			name:            "deload on failure (3 reps < 5 minReps)",
			currentWeight:   135.0,
			repsPerformed:   3,
			weightIncrement: 2.5,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  121.5, // 135 - 13.5 (10%)
			expectApplied:   true,
		},
		{
			name:            "deload on 4 reps",
			currentWeight:   135.0,
			repsPerformed:   4,
			weightIncrement: 2.5,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  121.5,
			expectApplied:   true,
		},
		{
			name:            "standard increment (5 reps = minReps)",
			currentWeight:   135.0,
			repsPerformed:   5,
			weightIncrement: 2.5,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  137.5,
			expectApplied:   true,
		},
		{
			name:            "standard increment (7 reps)",
			currentWeight:   135.0,
			repsPerformed:   7,
			weightIncrement: 2.5,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  137.5,
			expectApplied:   true,
		},
		{
			name:            "standard increment (9 reps, just below threshold)",
			currentWeight:   135.0,
			repsPerformed:   9,
			weightIncrement: 2.5,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  137.5,
			expectApplied:   true,
		},
		{
			name:            "double increment (10 reps = threshold)",
			currentWeight:   135.0,
			repsPerformed:   10,
			weightIncrement: 2.5,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  140.0, // 135 + 5.0 (2.5 * 2)
			expectApplied:   true,
		},
		{
			name:            "double increment (12 reps)",
			currentWeight:   135.0,
			repsPerformed:   12,
			weightIncrement: 2.5,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  140.0,
			expectApplied:   true,
		},
		// Lower body (5 lb increment)
		{
			name:            "lower body standard increment",
			currentWeight:   225.0,
			repsPerformed:   7,
			weightIncrement: 5.0,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  230.0,
			expectApplied:   true,
		},
		{
			name:            "lower body double increment",
			currentWeight:   225.0,
			repsPerformed:   10,
			weightIncrement: 5.0,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  235.0, // 225 + 10.0 (5.0 * 2)
			expectApplied:   true,
		},
		{
			name:            "lower body deload",
			currentWeight:   225.0,
			repsPerformed:   3,
			weightIncrement: 5.0,
			minReps:         5,
			doubleThreshold: 10,
			deloadPercent:   0.10,
			expectedWeight:  202.5, // 225 - 22.5 (10%)
			expectApplied:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gsp, err := greyskull.NewGreySkullProgression(
				uuid.New().String(),
				tc.name,
				tc.weightIncrement,
				tc.minReps,
				tc.doubleThreshold,
				tc.deloadPercent,
				progression.TrainingMax,
			)
			if err != nil {
				t.Fatalf("failed to create progression: %v", err)
			}

			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				MaxType:      progression.TrainingMax,
				CurrentValue: tc.currentWeight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(tc.repsPerformed),
					IsAMRAP:       true,
				},
			}

			result, err := gsp.Apply(context.Background(), params)
			if err != nil {
				t.Fatalf("Apply failed: %v", err)
			}

			if result.Applied != tc.expectApplied {
				t.Errorf("Applied: expected %v, got %v (reason: %s)", tc.expectApplied, result.Applied, result.Reason)
			}
			if result.Applied && result.NewValue != tc.expectedWeight {
				t.Errorf("NewValue: expected %.1f, got %.1f", tc.expectedWeight, result.NewValue)
			}
		})
	}
}

// TestGreySkullProgressionAccessory tests accessory lift progression rules.
func TestGreySkullProgressionAccessory(t *testing.T) {
	testCases := []struct {
		name           string
		currentWeight  float64
		repsPerformed  int
		expectedWeight float64
	}{
		// Accessory lifts: minReps=10, doubleThreshold=15
		{"deload (8 reps < 10)", 50.0, 8, 45.0},        // 50 - 5.0 (10%)
		{"deload (9 reps < 10)", 50.0, 9, 45.0},        // 50 - 5.0 (10%)
		{"standard (10 reps = minReps)", 50.0, 10, 52.5}, // 50 + 2.5
		{"standard (12 reps)", 50.0, 12, 52.5},          // 50 + 2.5
		{"standard (14 reps)", 50.0, 14, 52.5},          // 50 + 2.5
		{"double (15 reps = threshold)", 50.0, 15, 55.0}, // 50 + 5.0
		{"double (18 reps)", 50.0, 18, 55.0},            // 50 + 5.0
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gsp, err := greyskull.NewGreySkullAccessoryProgression(
				uuid.New().String(),
				"Accessory Test",
				2.5,
				progression.TrainingMax,
			)
			if err != nil {
				t.Fatalf("failed to create progression: %v", err)
			}

			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				MaxType:      progression.TrainingMax,
				CurrentValue: tc.currentWeight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(tc.repsPerformed),
					IsAMRAP:       true,
				},
			}

			result, err := gsp.Apply(context.Background(), params)
			if err != nil {
				t.Fatalf("Apply failed: %v", err)
			}

			if !result.Applied {
				t.Errorf("Expected progression to be applied, reason: %s", result.Reason)
			}
			if result.NewValue != tc.expectedWeight {
				t.Errorf("NewValue: expected %.1f, got %.1f", tc.expectedWeight, result.NewValue)
			}
		})
	}
}

// TestGreySkullProgressionNonAMRAPIgnored verifies progression does not apply to non-AMRAP sets.
func TestGreySkullProgressionNonAMRAPIgnored(t *testing.T) {
	gsp, err := greyskull.NewGreySkullMainLiftProgression(
		uuid.New().String(),
		"Test",
		2.5,
		progression.TrainingMax,
	)
	if err != nil {
		t.Fatalf("failed to create progression: %v", err)
	}

	params := progression.ProgressionContext{
		UserID:       uuid.New().String(),
		LiftID:       uuid.New().String(),
		MaxType:      progression.TrainingMax,
		CurrentValue: 135.0,
		TriggerEvent: progression.TriggerEvent{
			Type:          progression.TriggerAfterSet,
			Timestamp:     time.Now(),
			RepsPerformed: intPtr(10),
			IsAMRAP:       false, // Not an AMRAP set
		},
	}

	result, err := gsp.Apply(context.Background(), params)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if result.Applied {
		t.Error("Expected progression NOT to be applied for non-AMRAP set")
	}
	if result.Reason != "set is not marked as AMRAP" {
		t.Errorf("Unexpected reason: %s", result.Reason)
	}
}

// TestGreySkullProgressionWrongTriggerType verifies progression ignores wrong trigger types.
func TestGreySkullProgressionWrongTriggerType(t *testing.T) {
	gsp, err := greyskull.NewGreySkullMainLiftProgression(
		uuid.New().String(),
		"Test",
		2.5,
		progression.TrainingMax,
	)
	if err != nil {
		t.Fatalf("failed to create progression: %v", err)
	}

	wrongTriggers := []progression.TriggerType{
		progression.TriggerAfterSession,
		progression.TriggerAfterWeek,
		progression.TriggerAfterCycle,
		progression.TriggerOnFailure,
	}

	for _, triggerType := range wrongTriggers {
		t.Run(string(triggerType), func(t *testing.T) {
			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				MaxType:      progression.TrainingMax,
				CurrentValue: 135.0,
				TriggerEvent: progression.TriggerEvent{
					Type:          triggerType,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(10),
					IsAMRAP:       true,
				},
			}

			result, err := gsp.Apply(context.Background(), params)
			if err != nil {
				t.Fatalf("Apply failed: %v", err)
			}

			if result.Applied {
				t.Error("Expected progression NOT to be applied for wrong trigger type")
			}
		})
	}
}

// TestGreySkullDayTemplateSelection tests that the correct lifts appear for each variant.
func TestGreySkullDayTemplateSelection(t *testing.T) {
	testCases := []struct {
		variant       greyskull.Variant
		expectedLifts []string
	}{
		{
			variant:       greyskull.VariantA,
			expectedLifts: []string{"bench-press", "barbell-row", "squat", "tricep-extension", "ab-rollout"},
		},
		{
			variant:       greyskull.VariantB,
			expectedLifts: []string{"overhead-press", "chin-up", "deadlift", "bicep-curl", "shrug"},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.variant), func(t *testing.T) {
			lifts, err := greyskull.GetDayTemplate(tc.variant)
			if err != nil {
				t.Fatalf("GetDayTemplate failed: %v", err)
			}

			if len(lifts) != len(tc.expectedLifts) {
				t.Fatalf("expected %d lifts, got %d", len(tc.expectedLifts), len(lifts))
			}

			for i, expected := range tc.expectedLifts {
				if lifts[i].Slug != expected {
					t.Errorf("lift %d: expected %s, got %s", i+1, expected, lifts[i].Slug)
				}
			}
		})
	}
}

// TestGreySkullMainVsAccessoryLifts verifies correct categorization of lifts.
func TestGreySkullMainVsAccessoryLifts(t *testing.T) {
	testCases := []struct {
		variant              greyskull.Variant
		expectedMainLifts    []string
		expectedAccessories  []string
	}{
		{
			variant:              greyskull.VariantA,
			expectedMainLifts:    []string{"bench-press", "barbell-row", "squat"},
			expectedAccessories:  []string{"tricep-extension", "ab-rollout"},
		},
		{
			variant:              greyskull.VariantB,
			expectedMainLifts:    []string{"overhead-press", "chin-up", "deadlift"},
			expectedAccessories:  []string{"bicep-curl", "shrug"},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.variant)+" main lifts", func(t *testing.T) {
			mainLifts, err := greyskull.GetMainLiftSlugs(tc.variant)
			if err != nil {
				t.Fatalf("GetMainLiftSlugs failed: %v", err)
			}

			if len(mainLifts) != len(tc.expectedMainLifts) {
				t.Fatalf("expected %d main lifts, got %d", len(tc.expectedMainLifts), len(mainLifts))
			}

			for i, expected := range tc.expectedMainLifts {
				if mainLifts[i] != expected {
					t.Errorf("main lift %d: expected %s, got %s", i+1, expected, mainLifts[i])
				}
			}
		})

		t.Run(string(tc.variant)+" accessories", func(t *testing.T) {
			accessories, err := greyskull.GetAccessoryLiftSlugs(tc.variant)
			if err != nil {
				t.Fatalf("GetAccessoryLiftSlugs failed: %v", err)
			}

			if len(accessories) != len(tc.expectedAccessories) {
				t.Fatalf("expected %d accessories, got %d", len(tc.expectedAccessories), len(accessories))
			}

			for i, expected := range tc.expectedAccessories {
				if accessories[i] != expected {
					t.Errorf("accessory %d: expected %s, got %s", i+1, expected, accessories[i])
				}
			}
		})
	}
}

// TestGreySkullLiftInfoSetRepScheme verifies the set/rep scheme for each lift type.
func TestGreySkullLiftInfoSetRepScheme(t *testing.T) {
	for _, variant := range []greyskull.Variant{greyskull.VariantA, greyskull.VariantB} {
		t.Run(string(variant), func(t *testing.T) {
			lifts, err := greyskull.GetDayTemplate(variant)
			if err != nil {
				t.Fatalf("GetDayTemplate failed: %v", err)
			}

			// First 3 lifts are main lifts (3 sets, 5 reps)
			for i := 0; i < 3; i++ {
				lift := lifts[i]
				if lift.Sets != 3 {
					t.Errorf("Main lift %s: expected 3 sets, got %d", lift.Slug, lift.Sets)
				}
				if lift.Reps != 5 {
					t.Errorf("Main lift %s: expected 5 reps, got %d", lift.Slug, lift.Reps)
				}
				if !lift.IsAMRAP {
					t.Errorf("Main lift %s: expected IsAMRAP=true", lift.Slug)
				}
			}

			// Remaining lifts are accessories (3 sets, 10-12 reps)
			for i := 3; i < len(lifts); i++ {
				lift := lifts[i]
				if lift.Sets != 3 {
					t.Errorf("Accessory %s: expected 3 sets, got %d", lift.Slug, lift.Sets)
				}
				if lift.Reps != 10 && lift.Reps != 12 {
					t.Errorf("Accessory %s: expected 10 or 12 reps, got %d", lift.Slug, lift.Reps)
				}
				if !lift.IsAMRAP {
					t.Errorf("Accessory %s: expected IsAMRAP=true", lift.Slug)
				}
			}
		})
	}
}

// TestGreySkullFullWorkoutFlow tests the complete workout generation flow.
func TestGreySkullFullWorkoutFlow(t *testing.T) {
	// Simulate Week 1, Day 1 (Variant A)
	t.Run("Week 1 Day 1 - Variant A Workout", func(t *testing.T) {
		week := 1
		day := 1

		// Get the variant for this day
		variant, err := greyskull.GetVariantForDay(week, day)
		if err != nil {
			t.Fatalf("GetVariantForDay failed: %v", err)
		}
		if variant != greyskull.VariantA {
			t.Fatalf("Expected variant A for Week 1 Day 1, got %s", variant)
		}

		// Get the lifts for this variant
		lifts, err := greyskull.GetDayTemplate(variant)
		if err != nil {
			t.Fatalf("GetDayTemplate failed: %v", err)
		}

		// Create set scheme for main lifts
		mainLiftScheme, err := setscheme.NewGreySkullSetScheme(2, 5, 1, 5)
		if err != nil {
			t.Fatalf("Failed to create main lift scheme: %v", err)
		}

		// Verify main lift (bench press) set generation
		benchSets, err := mainLiftScheme.GenerateSets(135.0, setscheme.DefaultSetGenerationContext())
		if err != nil {
			t.Fatalf("Failed to generate bench sets: %v", err)
		}

		if len(benchSets) != 3 {
			t.Errorf("Expected 3 bench sets, got %d", len(benchSets))
		}

		// Simulate completing the AMRAP set with 8 reps
		benchProgression, err := greyskull.NewGreySkullMainLiftProgression(
			uuid.New().String(), "Bench Progression", 2.5, progression.TrainingMax)
		if err != nil {
			t.Fatalf("Failed to create bench progression: %v", err)
		}

		benchParams := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			MaxType:      progression.TrainingMax,
			CurrentValue: 135.0,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(8),
				IsAMRAP:       true,
			},
		}

		benchResult, err := benchProgression.Apply(context.Background(), benchParams)
		if err != nil {
			t.Fatalf("Bench progression failed: %v", err)
		}

		// Standard increment: 135 + 2.5 = 137.5
		if !benchResult.Applied {
			t.Error("Expected progression to apply")
		}
		if benchResult.NewValue != 137.5 {
			t.Errorf("Expected new bench weight 137.5, got %.1f", benchResult.NewValue)
		}

		// Verify all 5 lifts are present
		if len(lifts) != 5 {
			t.Errorf("Expected 5 lifts for variant A, got %d", len(lifts))
		}
	})

	// Simulate Week 1, Day 2 (Variant B)
	t.Run("Week 1 Day 2 - Variant B Workout", func(t *testing.T) {
		week := 1
		day := 2

		variant, err := greyskull.GetVariantForDay(week, day)
		if err != nil {
			t.Fatalf("GetVariantForDay failed: %v", err)
		}
		if variant != greyskull.VariantB {
			t.Fatalf("Expected variant B for Week 1 Day 2, got %s", variant)
		}

		lifts, err := greyskull.GetDayTemplate(variant)
		if err != nil {
			t.Fatalf("GetDayTemplate failed: %v", err)
		}

		// OHP should be first lift on variant B
		if lifts[0].Slug != "overhead-press" {
			t.Errorf("Expected overhead-press as first lift, got %s", lifts[0].Slug)
		}

		// Deadlift should be the third lift on variant B
		if lifts[2].Slug != "deadlift" {
			t.Errorf("Expected deadlift as third lift, got %s", lifts[2].Slug)
		}
	})
}

// TestGreySkullMultiWeekCycle tests a full multi-week cycle with progression.
func TestGreySkullMultiWeekCycle(t *testing.T) {
	type sessionResult struct {
		week          int
		day           int
		variant       greyskull.Variant
		benchWeight   float64
		amrapReps     int
		nextWeight    float64
	}

	// Simulate bench press progression over 4 weeks
	// Bench press only occurs on Variant A days
	// Week 1: Days 1, 3 = Variant A (2 bench sessions)
	// Week 2: Day 2 = Variant A (1 bench session)
	// Week 3: Days 1, 3 = Variant A (2 bench sessions)
	// Week 4: Day 2 = Variant A (1 bench session)
	// Total: 6 bench sessions in 4 weeks
	benchProgression, err := greyskull.NewGreySkullMainLiftProgression(
		uuid.New().String(), "Bench Progression", 2.5, progression.TrainingMax)
	if err != nil {
		t.Fatalf("Failed to create progression: %v", err)
	}

	userID := uuid.New().String()
	liftID := uuid.New().String()
	currentWeight := 135.0

	sessionResults := []sessionResult{}

	// Define AMRAP results for each bench session
	// Session 1: Week 1 Day 1 - 8 reps (standard)
	// Session 2: Week 1 Day 3 - 7 reps (standard)
	// Session 3: Week 2 Day 2 - 11 reps (double)
	// Session 4: Week 3 Day 1 - 6 reps (standard)
	// Session 5: Week 3 Day 3 - 3 reps (deload!)
	// Session 6: Week 4 Day 2 - 9 reps (standard after deload)
	amrapResultsBySession := []int{8, 7, 11, 6, 3, 9}
	sessionIndex := 0

	for week := 1; week <= 4; week++ {
		for day := 1; day <= 3; day++ {
			variant, _ := greyskull.GetVariantForDay(week, day)

			// Only progress bench on Variant A days
			if variant == greyskull.VariantA && sessionIndex < len(amrapResultsBySession) {
				amrapReps := amrapResultsBySession[sessionIndex]

				params := progression.ProgressionContext{
					UserID:       userID,
					LiftID:       liftID,
					MaxType:      progression.TrainingMax,
					CurrentValue: currentWeight,
					TriggerEvent: progression.TriggerEvent{
						Type:          progression.TriggerAfterSet,
						Timestamp:     time.Now(),
						RepsPerformed: intPtr(amrapReps),
						IsAMRAP:       true,
					},
				}

				result, err := benchProgression.Apply(context.Background(), params)
				if err != nil {
					t.Fatalf("Week %d Day %d progression failed: %v", week, day, err)
				}

				sessionResults = append(sessionResults, sessionResult{
					week:        week,
					day:         day,
					variant:     variant,
					benchWeight: currentWeight,
					amrapReps:   amrapReps,
					nextWeight:  result.NewValue,
				})

				currentWeight = result.NewValue
				sessionIndex++
			}
		}
	}

	// Verify progression history
	// Session 1: 135.0 -> 137.5 (+2.5 standard, 8 reps)
	// Session 2: 137.5 -> 140.0 (+2.5 standard, 7 reps)
	// Session 3: 140.0 -> 145.0 (+5.0 double, 11 reps)
	// Session 4: 145.0 -> 147.5 (+2.5 standard, 6 reps)
	// Session 5: 147.5 -> 132.75 (-14.75 deload 10%, 3 reps)
	// Session 6: 132.75 -> 135.25 (+2.5 standard, 9 reps)
	expectedWeights := []float64{
		135.0,   // Session 1: Start
		137.5,   // Session 2: After +2.5
		140.0,   // Session 3: After +2.5
		145.0,   // Session 4: After +5.0 (double)
		147.5,   // Session 5: After +2.5
		132.75,  // Session 6: After deload 10%
	}

	if len(sessionResults) != len(expectedWeights) {
		t.Fatalf("Expected %d sessions, got %d", len(expectedWeights), len(sessionResults))
	}

	// Verify each session's starting weight
	for i, session := range sessionResults {
		if session.benchWeight != expectedWeights[i] {
			t.Errorf("Session %d (Week %d Day %d): expected starting weight %.2f, got %.2f",
				i+1, session.week, session.day, expectedWeights[i], session.benchWeight)
		}
	}

	// Verify the final weight after all progressions
	expectedFinalWeight := 135.25 // 132.75 + 2.5
	if currentWeight != expectedFinalWeight {
		t.Errorf("Final weight: expected %.2f, got %.2f", expectedFinalWeight, currentWeight)
	}
}

// TestGreySkullStateAdvancement tests state advancement through a GreySkull LP program.
func TestGreySkullStateAdvancement(t *testing.T) {
	userID := uuid.New().String()
	programID := uuid.New().String()

	state, _ := userprogramstate.EnrollUser(
		userprogramstate.EnrollUserInput{
			UserID:    userID,
			ProgramID: programID,
		},
		uuid.New().String(),
	)

	// GreySkull LP: 3 training days per week, ongoing (no fixed cycle length)
	// Using 4 weeks for testing purposes
	ctx := userprogramstate.AdvancementContext{
		DaysInCurrentWeek: 3,
		CycleLengthWeeks:  4,
	}

	// Advance through 4 weeks (12 training days)
	weekVariants := make([][]greyskull.Variant, 4)
	for w := range weekVariants {
		weekVariants[w] = make([]greyskull.Variant, 3)
	}

	for day := 0; day < 12; day++ {
		// Get current week and day before advancing
		currentWeek := state.CurrentWeek
		var currentDay int
		if state.CurrentDayIndex != nil {
			currentDay = *state.CurrentDayIndex + 1
		} else {
			currentDay = 1
		}

		// Record variant for this position
		if currentWeek >= 1 && currentWeek <= 4 && currentDay >= 1 && currentDay <= 3 {
			variant, err := greyskull.GetVariantForDay(currentWeek, currentDay)
			if err != nil {
				t.Fatalf("Day %d: GetVariantForDay failed: %v", day+1, err)
			}
			weekVariants[currentWeek-1][currentDay-1] = variant
		}

		// Advance state
		result, valResult := userprogramstate.AdvanceState(state, ctx)
		if !valResult.Valid {
			t.Fatalf("Day %d: AdvanceState failed: %v", day+1, valResult.Errors)
		}
		state = result.NewState
	}

	// Verify the A/B pattern for each week
	t.Run("Week variant patterns", func(t *testing.T) {
		expectedPatterns := [][]greyskull.Variant{
			{greyskull.VariantA, greyskull.VariantB, greyskull.VariantA}, // Week 1
			{greyskull.VariantB, greyskull.VariantA, greyskull.VariantB}, // Week 2
			{greyskull.VariantA, greyskull.VariantB, greyskull.VariantA}, // Week 3
			{greyskull.VariantB, greyskull.VariantA, greyskull.VariantB}, // Week 4
		}

		for w, weekPattern := range weekVariants {
			for d, variant := range weekPattern {
				expected := expectedPatterns[w][d]
				if variant != expected {
					t.Errorf("Week %d Day %d: expected %s, got %s", w+1, d+1, expected, variant)
				}
			}
		}
	})

	// Verify cycle completed after 4 weeks
	t.Run("Cycle completion", func(t *testing.T) {
		if state.CurrentWeek != 1 {
			t.Errorf("CurrentWeek: expected 1 (reset), got %d", state.CurrentWeek)
		}
		if state.CyclesSinceStart != 1 {
			t.Errorf("CyclesSinceStart: expected 1, got %d", state.CyclesSinceStart)
		}
	})
}

// TestGreySkullSetSchemeType verifies the set scheme type discriminator.
func TestGreySkullSetSchemeType(t *testing.T) {
	scheme, err := setscheme.NewGreySkullSetScheme(2, 5, 1, 5)
	if err != nil {
		t.Fatalf("Failed to create scheme: %v", err)
	}

	if scheme.Type() != setscheme.TypeGreySkull {
		t.Errorf("Expected type %s, got %s", setscheme.TypeGreySkull, scheme.Type())
	}
}

// TestGreySkullProgressionType verifies the progression type discriminator.
func TestGreySkullProgressionType(t *testing.T) {
	gsp, err := greyskull.NewGreySkullMainLiftProgression(
		uuid.New().String(), "Test", 2.5, progression.TrainingMax)
	if err != nil {
		t.Fatalf("Failed to create progression: %v", err)
	}

	if gsp.Type() != progression.TypeGreySkull {
		t.Errorf("Expected type %s, got %s", progression.TypeGreySkull, gsp.Type())
	}

	if gsp.TriggerType() != progression.TriggerAfterSet {
		t.Errorf("Expected trigger type %s, got %s", progression.TriggerAfterSet, gsp.TriggerType())
	}
}

// TestGreySkullDayTemplateForWeek tests the convenience function combining variant + template.
func TestGreySkullDayTemplateForWeek(t *testing.T) {
	testCases := []struct {
		week           int
		day            int
		expectedFirst  string
	}{
		{1, 1, "bench-press"},     // Week 1 Day 1 = Variant A
		{1, 2, "overhead-press"},  // Week 1 Day 2 = Variant B
		{1, 3, "bench-press"},     // Week 1 Day 3 = Variant A
		{2, 1, "overhead-press"},  // Week 2 Day 1 = Variant B
		{2, 2, "bench-press"},     // Week 2 Day 2 = Variant A
		{2, 3, "overhead-press"},  // Week 2 Day 3 = Variant B
	}

	for _, tc := range testCases {
		t.Run("Week "+string(rune('0'+tc.week))+" Day "+string(rune('0'+tc.day)), func(t *testing.T) {
			lifts, err := greyskull.GetDayTemplateForWeek(tc.week, tc.day)
			if err != nil {
				t.Fatalf("GetDayTemplateForWeek failed: %v", err)
			}

			if lifts[0].Slug != tc.expectedFirst {
				t.Errorf("Expected first lift %s, got %s", tc.expectedFirst, lifts[0].Slug)
			}
		})
	}
}

// TestGreySkullEdgeCases tests edge cases in the GreySkull LP implementation.
func TestGreySkullEdgeCases(t *testing.T) {
	t.Run("deload from very low weight", func(t *testing.T) {
		gsp, err := greyskull.NewGreySkullMainLiftProgression(
			uuid.New().String(), "Test", 2.5, progression.TrainingMax)
		if err != nil {
			t.Fatalf("Failed to create progression: %v", err)
		}

		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			MaxType:      progression.TrainingMax,
			CurrentValue: 45.0, // Just the bar
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(3), // Failed
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(context.Background(), params)
		if err != nil {
			t.Fatalf("Apply failed: %v", err)
		}

		// 45 - 4.5 (10%) = 40.5
		if result.NewValue != 40.5 {
			t.Errorf("Expected 40.5, got %.1f", result.NewValue)
		}
	})

	t.Run("progression with large rep count", func(t *testing.T) {
		gsp, err := greyskull.NewGreySkullMainLiftProgression(
			uuid.New().String(), "Test", 2.5, progression.TrainingMax)
		if err != nil {
			t.Fatalf("Failed to create progression: %v", err)
		}

		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			MaxType:      progression.TrainingMax,
			CurrentValue: 95.0,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(20), // Very high reps
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(context.Background(), params)
		if err != nil {
			t.Fatalf("Apply failed: %v", err)
		}

		// Still double increment (not triple or more)
		// 95 + 5.0 = 100
		if result.NewValue != 100.0 {
			t.Errorf("Expected 100.0, got %.1f", result.NewValue)
		}
	})

	t.Run("invalid week number", func(t *testing.T) {
		_, err := greyskull.GetVariantForDay(0, 1)
		if err == nil {
			t.Error("Expected error for week 0")
		}

		_, err = greyskull.GetVariantForDay(-1, 1)
		if err == nil {
			t.Error("Expected error for negative week")
		}
	})

	t.Run("invalid day position", func(t *testing.T) {
		_, err := greyskull.GetVariantForDay(1, 0)
		if err == nil {
			t.Error("Expected error for day 0")
		}

		_, err = greyskull.GetVariantForDay(1, 4)
		if err == nil {
			t.Error("Expected error for day 4")
		}
	})
}

// TestGreySkullConsistencyAcrossWeeks verifies the A/B pattern is consistent across many weeks.
func TestGreySkullConsistencyAcrossWeeks(t *testing.T) {
	// Verify pattern consistency for 52 weeks (one year)
	for week := 1; week <= 52; week++ {
		isOddWeek := week%2 == 1
		expectedPattern := []greyskull.Variant{
			greyskull.VariantA, greyskull.VariantB, greyskull.VariantA, // odd week
		}
		if !isOddWeek {
			expectedPattern = []greyskull.Variant{
				greyskull.VariantB, greyskull.VariantA, greyskull.VariantB, // even week
			}
		}

		for day := 1; day <= 3; day++ {
			variant, err := greyskull.GetVariantForDay(week, day)
			if err != nil {
				t.Fatalf("Week %d Day %d: unexpected error: %v", week, day, err)
			}
			if variant != expectedPattern[day-1] {
				t.Errorf("Week %d Day %d: expected %s, got %s",
					week, day, expectedPattern[day-1], variant)
			}
		}
	}
}
