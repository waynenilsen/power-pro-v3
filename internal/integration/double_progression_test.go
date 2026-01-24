// Package integration provides integration tests for cross-component behavior.
// This file tests the DoubleProgression strategy combined with RepRangeSetScheme.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// intPtr creates a pointer to an int value.
func intPtr(i int) *int {
	return &i
}

// TestDoubleProgressionBasicCycle tests a complete double progression cycle where
// a lifter progresses reps from 8 to 12 across sessions before adding weight.
func TestDoubleProgressionBasicCycle(t *testing.T) {
	ctx := context.Background()

	// Create a standard double progression: add 5lb when hitting 12 reps
	dp, err := progression.NewDoubleProgression(
		uuid.New().String(),
		"3x8-12 Curls",
		5.0, // Weight increment
		progression.TrainingMax,
		progression.TriggerAfterSet,
	)
	if err != nil {
		t.Fatalf("Failed to create double progression: %v", err)
	}

	userID := uuid.New().String()
	liftID := uuid.New().String()
	startingWeight := 50.0
	ceiling := 12

	// Session 1: 8, 8, 8 reps - no progression (not at ceiling)
	t.Run("session 1 - starting reps, no progression", func(t *testing.T) {
		repsPerSet := []int{8, 8, 8}

		for setNum, reps := range repsPerSet {
			params := progression.ProgressionContext{
				UserID:       userID,
				LiftID:       liftID,
				MaxType:      progression.TrainingMax,
				CurrentValue: startingWeight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(reps),
					MaxReps:       intPtr(ceiling),
				},
			}

			result, err := dp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("Set %d: unexpected error: %v", setNum+1, err)
			}
			if result.Applied {
				t.Errorf("Set %d: expected no progression at 8 reps, but Applied=true", setNum+1)
			}
			if result.NewValue != startingWeight {
				t.Errorf("Set %d: expected weight unchanged at %.1f, got %.1f", setNum+1, startingWeight, result.NewValue)
			}
		}
	})

	// Session 2: 10, 10, 10 reps - no progression (still below ceiling)
	t.Run("session 2 - progressing reps, no weight change", func(t *testing.T) {
		repsPerSet := []int{10, 10, 10}

		for setNum, reps := range repsPerSet {
			params := progression.ProgressionContext{
				UserID:       userID,
				LiftID:       liftID,
				MaxType:      progression.TrainingMax,
				CurrentValue: startingWeight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(reps),
					MaxReps:       intPtr(ceiling),
				},
			}

			result, err := dp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("Set %d: unexpected error: %v", setNum+1, err)
			}
			if result.Applied {
				t.Errorf("Set %d: expected no progression at 10 reps, but Applied=true", setNum+1)
			}
		}
	})

	// Session 3: 12, 12, 12 reps - progression triggered on each set hitting ceiling
	t.Run("session 3 - hitting ceiling triggers progression", func(t *testing.T) {
		repsPerSet := []int{12, 12, 12}

		for setNum, reps := range repsPerSet {
			params := progression.ProgressionContext{
				UserID:       userID,
				LiftID:       liftID,
				MaxType:      progression.TrainingMax,
				CurrentValue: startingWeight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(reps),
					MaxReps:       intPtr(ceiling),
				},
			}

			result, err := dp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("Set %d: unexpected error: %v", setNum+1, err)
			}
			if !result.Applied {
				t.Errorf("Set %d: expected progression at 12 reps, reason: %s", setNum+1, result.Reason)
			}
			if result.Delta != 5.0 {
				t.Errorf("Set %d: expected delta 5.0, got %.1f", setNum+1, result.Delta)
			}
			if result.NewValue != 55.0 {
				t.Errorf("Set %d: expected new weight 55.0, got %.1f", setNum+1, result.NewValue)
			}
		}
	})
}

// TestDoubleProgressionPartialCeilingHit tests the scenario where a user hits the
// ceiling on some but not all sets. Each set is evaluated independently.
func TestDoubleProgressionPartialCeilingHit(t *testing.T) {
	ctx := context.Background()

	dp, err := progression.NewDoubleProgression(
		uuid.New().String(),
		"3x8-12 Tricep Pushdowns",
		2.5,
		progression.TrainingMax,
		progression.TriggerAfterSet,
	)
	if err != nil {
		t.Fatalf("Failed to create double progression: %v", err)
	}

	userID := uuid.New().String()
	liftID := uuid.New().String()
	weight := 30.0
	ceiling := 12

	// Simulate mixed set: 12, 10, 11 reps
	// Only first set hits ceiling, others don't
	testCases := []struct {
		setNum          int
		reps            int
		expectedApplied bool
	}{
		{1, 12, true},  // Hit ceiling - progression applies
		{2, 10, false}, // Below ceiling - no progression
		{3, 11, false}, // Below ceiling - no progression
	}

	for _, tc := range testCases {
		t.Run("set "+string(rune('0'+tc.setNum))+" with "+string(rune('0'+tc.reps/10))+string(rune('0'+tc.reps%10))+" reps", func(t *testing.T) {
			params := progression.ProgressionContext{
				UserID:       userID,
				LiftID:       liftID,
				MaxType:      progression.TrainingMax,
				CurrentValue: weight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(tc.reps),
					MaxReps:       intPtr(ceiling),
				},
			}

			result, err := dp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Applied != tc.expectedApplied {
				if tc.expectedApplied {
					t.Errorf("expected progression for %d reps, reason: %s", tc.reps, result.Reason)
				} else {
					t.Errorf("expected no progression for %d reps, but Applied=true", tc.reps)
				}
			}
		})
	}
}

// TestRepRangeSetSchemeGeneration tests that RepRangeSetScheme generates
// correct sets for double progression workouts.
func TestRepRangeSetSchemeGeneration(t *testing.T) {
	// Create a standard 3x8-12 scheme
	scheme, err := setscheme.NewRepRangeSetScheme(3, 8, 12)
	if err != nil {
		t.Fatalf("Failed to create rep range scheme: %v", err)
	}

	t.Run("generates correct number of sets", func(t *testing.T) {
		sets, err := scheme.GenerateSets(50.0, setscheme.DefaultSetGenerationContext())
		if err != nil {
			t.Fatalf("Failed to generate sets: %v", err)
		}

		if len(sets) != 3 {
			t.Errorf("Expected 3 sets, got %d", len(sets))
		}
	})

	t.Run("sets have correct weight", func(t *testing.T) {
		baseWeight := 75.0
		sets, _ := scheme.GenerateSets(baseWeight, setscheme.DefaultSetGenerationContext())

		for i, set := range sets {
			if set.Weight != baseWeight {
				t.Errorf("Set %d: expected weight %.1f, got %.1f", i+1, baseWeight, set.Weight)
			}
		}
	})

	t.Run("target reps is minimum of range", func(t *testing.T) {
		sets, _ := scheme.GenerateSets(50.0, setscheme.DefaultSetGenerationContext())

		for i, set := range sets {
			if set.TargetReps != 8 {
				t.Errorf("Set %d: expected target reps 8 (minimum), got %d", i+1, set.TargetReps)
			}
		}
	})

	t.Run("all sets are work sets", func(t *testing.T) {
		sets, _ := scheme.GenerateSets(50.0, setscheme.DefaultSetGenerationContext())

		for i, set := range sets {
			if !set.IsWorkSet {
				t.Errorf("Set %d: expected IsWorkSet=true", i+1)
			}
		}
	})

	t.Run("set numbers are 1-indexed", func(t *testing.T) {
		sets, _ := scheme.GenerateSets(50.0, setscheme.DefaultSetGenerationContext())

		for i, set := range sets {
			expected := i + 1
			if set.SetNumber != expected {
				t.Errorf("Set %d: expected SetNumber %d, got %d", i+1, expected, set.SetNumber)
			}
		}
	})
}

// TestDoubleProgressionWithRepRangeScheme tests the integration between
// RepRangeSetScheme and DoubleProgression as used in real workout logging.
func TestDoubleProgressionWithRepRangeScheme(t *testing.T) {
	ctx := context.Background()

	// Create rep range scheme: 3 sets, 8-12 reps
	scheme, err := setscheme.NewRepRangeSetScheme(3, 8, 12)
	if err != nil {
		t.Fatalf("Failed to create scheme: %v", err)
	}

	// Create double progression that adds 5lb when hitting ceiling
	dp, err := progression.NewDoubleProgression(
		uuid.New().String(),
		"Accessory Double Progression",
		5.0,
		progression.TrainingMax,
		progression.TriggerAfterSet,
	)
	if err != nil {
		t.Fatalf("Failed to create progression: %v", err)
	}

	userID := uuid.New().String()
	liftID := uuid.New().String()
	startingWeight := 40.0

	// Generate prescribed sets
	prescribedSets, _ := scheme.GenerateSets(startingWeight, setscheme.DefaultSetGenerationContext())
	if len(prescribedSets) != 3 {
		t.Fatalf("Expected 3 prescribed sets, got %d", len(prescribedSets))
	}

	// Verify scheme provides min reps as target
	for _, set := range prescribedSets {
		if set.TargetReps != scheme.MinReps {
			t.Errorf("Prescribed set target reps should be MinReps=%d, got %d", scheme.MinReps, set.TargetReps)
		}
	}

	// Simulate logging where user hits ceiling on all sets
	t.Run("all sets at ceiling triggers progression", func(t *testing.T) {
		progressionTriggered := 0

		for _, set := range prescribedSets {
			// User performs MaxReps (ceiling)
			repsPerformed := scheme.MaxReps

			params := progression.ProgressionContext{
				UserID:       userID,
				LiftID:       liftID,
				MaxType:      progression.TrainingMax,
				CurrentValue: set.Weight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(repsPerformed),
					MaxReps:       intPtr(scheme.MaxReps), // Ceiling from scheme
				},
			}

			result, err := dp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("Set %d: unexpected error: %v", set.SetNumber, err)
			}

			if result.Applied {
				progressionTriggered++
			}
		}

		// Each set hitting ceiling should trigger progression
		if progressionTriggered != 3 {
			t.Errorf("Expected 3 progressions triggered, got %d", progressionTriggered)
		}
	})

	// Simulate logging where user is at minimum reps (no progression)
	t.Run("all sets at minimum does not trigger progression", func(t *testing.T) {
		for _, set := range prescribedSets {
			// User performs MinReps (floor, not ceiling)
			repsPerformed := scheme.MinReps

			params := progression.ProgressionContext{
				UserID:       userID,
				LiftID:       liftID,
				MaxType:      progression.TrainingMax,
				CurrentValue: set.Weight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(repsPerformed),
					MaxReps:       intPtr(scheme.MaxReps),
				},
			}

			result, err := dp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("Set %d: unexpected error: %v", set.SetNumber, err)
			}

			if result.Applied {
				t.Errorf("Set %d: should not trigger progression at min reps", set.SetNumber)
			}
		}
	})
}

// TestRedditPPLAccessoryPattern tests the double progression pattern as used
// in the Reddit PPL 6-Day program for accessory exercises like bicep curls.
func TestRedditPPLAccessoryPattern(t *testing.T) {
	ctx := context.Background()

	// Reddit PPL accessories: 3 sets x 8-12 reps
	// Progress weight when all sets hit 12 reps
	scheme, err := setscheme.NewRepRangeSetScheme(3, 8, 12)
	if err != nil {
		t.Fatalf("Failed to create scheme: %v", err)
	}

	// Typical accessory progression: small 5lb jumps
	dp, err := progression.NewDoubleProgression(
		uuid.New().String(),
		"Reddit PPL Bicep Curls",
		5.0,
		progression.TrainingMax,
		progression.TriggerAfterSet,
	)
	if err != nil {
		t.Fatalf("Failed to create progression: %v", err)
	}

	userID := uuid.New().String()
	liftID := uuid.New().String()

	// Simulate multi-week progression
	type weekSession struct {
		week           int
		startWeight    float64
		repsPerSet     []int
		expectedWeight float64 // Weight after applying progressions
		description    string
	}

	sessions := []weekSession{
		{1, 20.0, []int{8, 8, 8}, 20.0, "week 1 - starting at floor"},
		{2, 20.0, []int{10, 9, 9}, 20.0, "week 2 - building reps"},
		{3, 20.0, []int{11, 11, 10}, 20.0, "week 3 - almost there"},
		{4, 20.0, []int{12, 12, 12}, 25.0, "week 4 - hit ceiling, add weight"},
		{5, 25.0, []int{8, 8, 7}, 25.0, "week 5 - reset to floor at new weight"},
		{6, 25.0, []int{10, 10, 10}, 25.0, "week 6 - building again"},
		{7, 25.0, []int{12, 12, 12}, 30.0, "week 7 - ceiling again, add weight"},
	}

	for _, session := range sessions {
		t.Run(session.description, func(t *testing.T) {
			currentWeight := session.startWeight
			sets, _ := scheme.GenerateSets(currentWeight, setscheme.DefaultSetGenerationContext())

			for i, set := range sets {
				params := progression.ProgressionContext{
					UserID:       userID,
					LiftID:       liftID,
					MaxType:      progression.TrainingMax,
					CurrentValue: set.Weight,
					TriggerEvent: progression.TriggerEvent{
						Type:          progression.TriggerAfterSet,
						Timestamp:     time.Now(),
						RepsPerformed: intPtr(session.repsPerSet[i]),
						MaxReps:       intPtr(scheme.MaxReps),
					},
				}

				result, err := dp.Apply(ctx, params)
				if err != nil {
					t.Fatalf("Set %d: unexpected error: %v", i+1, err)
				}

				// Only update weight if progression applied
				if result.Applied {
					currentWeight = result.NewValue
				}
			}

			// Verify final weight matches expected
			if currentWeight != session.expectedWeight {
				t.Errorf("Expected weight %.1f after session, got %.1f", session.expectedWeight, currentWeight)
			}
		})
	}
}

// TestDoubleProgressionExceedingCeiling tests that exceeding the ceiling
// (performing more reps than the max) still triggers progression.
func TestDoubleProgressionExceedingCeiling(t *testing.T) {
	ctx := context.Background()

	dp, err := progression.NewDoubleProgression(
		uuid.New().String(),
		"Test Progression",
		5.0,
		progression.TrainingMax,
		progression.TriggerAfterSet,
	)
	if err != nil {
		t.Fatalf("Failed to create progression: %v", err)
	}

	params := progression.ProgressionContext{
		UserID:       uuid.New().String(),
		LiftID:       uuid.New().String(),
		MaxType:      progression.TrainingMax,
		CurrentValue: 100.0,
		TriggerEvent: progression.TriggerEvent{
			Type:          progression.TriggerAfterSet,
			Timestamp:     time.Now(),
			RepsPerformed: intPtr(15), // Exceeded ceiling of 12
			MaxReps:       intPtr(12),
		},
	}

	result, err := dp.Apply(ctx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Applied {
		t.Errorf("Expected progression when exceeding ceiling, reason: %s", result.Reason)
	}
	if result.NewValue != 105.0 {
		t.Errorf("Expected new weight 105.0, got %.1f", result.NewValue)
	}
}

// TestDoubleProgressionVariousRepRanges tests double progression with different
// rep range configurations commonly used in bodybuilding programs.
func TestDoubleProgressionVariousRepRanges(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name            string
		minReps         int
		maxReps         int
		sets            int
		weightIncrement float64
	}{
		{"3x8-12 (standard)", 8, 12, 3, 5.0},
		{"4x6-10 (strength-hypertrophy)", 6, 10, 4, 5.0},
		{"3x12-15 (high rep)", 12, 15, 3, 2.5},
		{"2x15-20 (endurance)", 15, 20, 2, 2.5},
		{"3x10-12 (tight range)", 10, 12, 3, 5.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scheme, err := setscheme.NewRepRangeSetScheme(tc.sets, tc.minReps, tc.maxReps)
			if err != nil {
				t.Fatalf("Failed to create scheme: %v", err)
			}

			dp, err := progression.NewDoubleProgression(
				uuid.New().String(),
				tc.name,
				tc.weightIncrement,
				progression.TrainingMax,
				progression.TriggerAfterSet,
			)
			if err != nil {
				t.Fatalf("Failed to create progression: %v", err)
			}

			startWeight := 50.0
			sets, _ := scheme.GenerateSets(startWeight, setscheme.DefaultSetGenerationContext())

			// Verify scheme generates correct number of sets
			if len(sets) != tc.sets {
				t.Errorf("Expected %d sets, got %d", tc.sets, len(sets))
			}

			// Test progression at ceiling
			for _, set := range sets {
				params := progression.ProgressionContext{
					UserID:       uuid.New().String(),
					LiftID:       uuid.New().String(),
					MaxType:      progression.TrainingMax,
					CurrentValue: set.Weight,
					TriggerEvent: progression.TriggerEvent{
						Type:          progression.TriggerAfterSet,
						Timestamp:     time.Now(),
						RepsPerformed: intPtr(tc.maxReps), // Hit ceiling
						MaxReps:       intPtr(tc.maxReps),
					},
				}

				result, err := dp.Apply(ctx, params)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !result.Applied {
					t.Errorf("Expected progression at ceiling %d reps", tc.maxReps)
				}
				if result.Delta != tc.weightIncrement {
					t.Errorf("Expected delta %.1f, got %.1f", tc.weightIncrement, result.Delta)
				}
			}

			// Test no progression below ceiling
			belowCeiling := tc.maxReps - 1
			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				MaxType:      progression.TrainingMax,
				CurrentValue: startWeight,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(belowCeiling),
					MaxReps:       intPtr(tc.maxReps),
				},
			}

			result, _ := dp.Apply(ctx, params)
			if result.Applied {
				t.Errorf("Should not progress at %d reps (below ceiling %d)", belowCeiling, tc.maxReps)
			}
		})
	}
}

// TestRepRangeSchemeValidation tests validation of RepRangeSetScheme parameters.
func TestRepRangeSchemeValidation(t *testing.T) {
	t.Run("valid scheme", func(t *testing.T) {
		_, err := setscheme.NewRepRangeSetScheme(3, 8, 12)
		if err != nil {
			t.Errorf("Expected valid scheme, got error: %v", err)
		}
	})

	t.Run("invalid - zero sets", func(t *testing.T) {
		_, err := setscheme.NewRepRangeSetScheme(0, 8, 12)
		if err == nil {
			t.Error("Expected error for zero sets")
		}
	})

	t.Run("invalid - zero min reps", func(t *testing.T) {
		_, err := setscheme.NewRepRangeSetScheme(3, 0, 12)
		if err == nil {
			t.Error("Expected error for zero min reps")
		}
	})

	t.Run("invalid - max reps less than min reps", func(t *testing.T) {
		_, err := setscheme.NewRepRangeSetScheme(3, 12, 8)
		if err == nil {
			t.Error("Expected error when maxReps < minReps")
		}
	})

	t.Run("valid - equal min and max reps", func(t *testing.T) {
		scheme, err := setscheme.NewRepRangeSetScheme(3, 10, 10)
		if err != nil {
			t.Errorf("Equal min/max should be valid, got error: %v", err)
		}
		if scheme.MinReps != scheme.MaxReps {
			t.Error("MinReps should equal MaxReps")
		}
	})
}
