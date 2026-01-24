// Package integration provides integration tests for cross-component behavior.
// This file tests the nSuns CAP3 rotation system using RotationLookup and WeeklyLookup.
package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/rotationlookup"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
	"github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"
)

// CAP3 rotation positions
const (
	CAP3PositionDeadlift = 0 // Week 1: Deadlift High Intensity
	CAP3PositionSquat    = 1 // Week 2: Squat High Intensity
	CAP3PositionBench    = 2 // Week 3: Bench High Intensity
)

// CAP3 lift identifiers
const (
	CAP3LiftDeadlift = "deadlift"
	CAP3LiftSquat    = "squat"
	CAP3LiftBench    = "bench"
)

// CAP3 percentages by intensity level
var (
	// High Intensity AMRAP percentages: 79.5% -> 83.5% -> 88.5%
	CAP3HighIntensityPercentages = []float64{79.5, 83.5, 88.5}
	CAP3HighIntensityReps        = []int{6, 4, 2}

	// Medium Day percentages: 77% base, +5lb for sets 5-6, +10lb for set 7-8
	// For simplicity in testing, we use percentage modifier pattern
	CAP3MediumPercentages = []float64{77, 77, 77, 77, 77, 77, 77, 77}
	CAP3MediumReps        = []int{3, 3, 3, 3, 3, 3, 3, 3}

	// Volume Day percentages: 73.5% base
	CAP3VolumePercentages = []float64{73.5, 73.5, 73.5, 73.5, 73.5, 73.5, 73.5}
	CAP3VolumeReps        = []int{4, 4, 4, 4, 4, 4, 4}
)

// createCAP3RotationLookup creates the 3-position rotation for CAP3.
// Position 0: Deadlift Focus (Week 1)
// Position 1: Squat Focus (Week 2)
// Position 2: Bench Focus (Week 3)
func createCAP3RotationLookup() *rotationlookup.RotationLookup {
	programID := uuid.New().String()
	rl, _ := rotationlookup.CreateRotationLookup(
		rotationlookup.CreateRotationLookupInput{
			Name: "CAP3 Rotation",
			Entries: []rotationlookup.RotationLookupEntry{
				{
					Position:       CAP3PositionDeadlift,
					LiftIdentifier: CAP3LiftDeadlift,
					Description:    "Deadlift High Intensity AMRAP Week",
				},
				{
					Position:       CAP3PositionSquat,
					LiftIdentifier: CAP3LiftSquat,
					Description:    "Squat High Intensity AMRAP Week",
				},
				{
					Position:       CAP3PositionBench,
					LiftIdentifier: CAP3LiftBench,
					Description:    "Bench High Intensity AMRAP Week",
				},
			},
			ProgramID: &programID,
		},
		uuid.New().String(),
	)
	return rl
}

// createCAP3HighIntensityWeeklyLookup creates weekly lookup for High Intensity AMRAP day.
func createCAP3HighIntensityWeeklyLookup() *weeklylookup.WeeklyLookup {
	programID := uuid.New().String()
	wl, _ := weeklylookup.CreateWeeklyLookup(
		weeklylookup.CreateWeeklyLookupInput{
			Name: "CAP3 High Intensity",
			Entries: []weeklylookup.WeeklyLookupEntry{
				{
					WeekNumber:  1,
					Percentages: CAP3HighIntensityPercentages,
					Reps:        CAP3HighIntensityReps,
				},
			},
			ProgramID: &programID,
		},
		uuid.New().String(),
	)
	return wl
}

// createCAP3MediumWeeklyLookup creates weekly lookup for Medium Volume day.
func createCAP3MediumWeeklyLookup() *weeklylookup.WeeklyLookup {
	programID := uuid.New().String()
	wl, _ := weeklylookup.CreateWeeklyLookup(
		weeklylookup.CreateWeeklyLookupInput{
			Name: "CAP3 Medium",
			Entries: []weeklylookup.WeeklyLookupEntry{
				{
					WeekNumber:  1,
					Percentages: CAP3MediumPercentages,
					Reps:        CAP3MediumReps,
				},
			},
			ProgramID: &programID,
		},
		uuid.New().String(),
	)
	return wl
}

// createCAP3VolumeWeeklyLookup creates weekly lookup for Volume day.
func createCAP3VolumeWeeklyLookup() *weeklylookup.WeeklyLookup {
	programID := uuid.New().String()
	wl, _ := weeklylookup.CreateWeeklyLookup(
		weeklylookup.CreateWeeklyLookupInput{
			Name: "CAP3 Volume",
			Entries: []weeklylookup.WeeklyLookupEntry{
				{
					WeekNumber:  1,
					Percentages: CAP3VolumePercentages,
					Reps:        CAP3VolumeReps,
				},
			},
			ProgramID: &programID,
		},
		uuid.New().String(),
	)
	return wl
}

// TestCAP3RotationLookupCreation tests that the CAP3 rotation lookup is created correctly.
func TestCAP3RotationLookupCreation(t *testing.T) {
	rl := createCAP3RotationLookup()

	t.Run("has 3 positions", func(t *testing.T) {
		if rl.Length() != 3 {
			t.Errorf("Expected 3 positions, got %d", rl.Length())
		}
	})

	t.Run("position 0 is deadlift", func(t *testing.T) {
		entry := rl.GetByPosition(CAP3PositionDeadlift)
		if entry == nil {
			t.Fatal("Position 0 entry is nil")
		}
		if entry.LiftIdentifier != CAP3LiftDeadlift {
			t.Errorf("Expected deadlift, got %s", entry.LiftIdentifier)
		}
	})

	t.Run("position 1 is squat", func(t *testing.T) {
		entry := rl.GetByPosition(CAP3PositionSquat)
		if entry == nil {
			t.Fatal("Position 1 entry is nil")
		}
		if entry.LiftIdentifier != CAP3LiftSquat {
			t.Errorf("Expected squat, got %s", entry.LiftIdentifier)
		}
	})

	t.Run("position 2 is bench", func(t *testing.T) {
		entry := rl.GetByPosition(CAP3PositionBench)
		if entry == nil {
			t.Fatal("Position 2 entry is nil")
		}
		if entry.LiftIdentifier != CAP3LiftBench {
			t.Errorf("Expected bench, got %s", entry.LiftIdentifier)
		}
	})

	t.Run("contains all three lifts", func(t *testing.T) {
		lifts := []string{CAP3LiftDeadlift, CAP3LiftSquat, CAP3LiftBench}
		for _, lift := range lifts {
			if !rl.ContainsLift(lift) {
				t.Errorf("Expected rotation to contain %s", lift)
			}
		}
	})
}

// TestCAP3RotationFocusDetection tests that we can correctly detect which lift is in focus.
func TestCAP3RotationFocusDetection(t *testing.T) {
	rl := createCAP3RotationLookup()

	testCases := []struct {
		name            string
		rotationPos     int
		liftIdentifier  string
		expectedInFocus bool
	}{
		// Week 1: Deadlift focus
		{"Week 1 - Deadlift is focus", CAP3PositionDeadlift, CAP3LiftDeadlift, true},
		{"Week 1 - Squat not focus", CAP3PositionDeadlift, CAP3LiftSquat, false},
		{"Week 1 - Bench not focus", CAP3PositionDeadlift, CAP3LiftBench, false},

		// Week 2: Squat focus
		{"Week 2 - Deadlift not focus", CAP3PositionSquat, CAP3LiftDeadlift, false},
		{"Week 2 - Squat is focus", CAP3PositionSquat, CAP3LiftSquat, true},
		{"Week 2 - Bench not focus", CAP3PositionSquat, CAP3LiftBench, false},

		// Week 3: Bench focus
		{"Week 3 - Deadlift not focus", CAP3PositionBench, CAP3LiftDeadlift, false},
		{"Week 3 - Squat not focus", CAP3PositionBench, CAP3LiftSquat, false},
		{"Week 3 - Bench is focus", CAP3PositionBench, CAP3LiftBench, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &loadstrategy.LookupContext{
				RotationPosition: tc.rotationPos,
				RotationLookup:   rl,
			}

			inFocus := ctx.IsLiftInRotationFocus(tc.liftIdentifier)
			if inFocus != tc.expectedInFocus {
				t.Errorf("Expected IsLiftInRotationFocus=%v, got %v", tc.expectedInFocus, inFocus)
			}
		})
	}
}

// TestCAP3HighIntensityWeeklyLookup tests the high intensity AMRAP percentages.
func TestCAP3HighIntensityWeeklyLookup(t *testing.T) {
	wl := createCAP3HighIntensityWeeklyLookup()

	t.Run("week 1 entry exists", func(t *testing.T) {
		entry := wl.GetByWeekNumber(1)
		if entry == nil {
			t.Fatal("Week 1 entry is nil")
		}
	})

	t.Run("has correct percentages", func(t *testing.T) {
		entry := wl.GetByWeekNumber(1)
		expectedPercentages := []float64{79.5, 83.5, 88.5}

		if len(entry.Percentages) != len(expectedPercentages) {
			t.Fatalf("Expected %d percentages, got %d", len(expectedPercentages), len(entry.Percentages))
		}

		for i, expected := range expectedPercentages {
			if entry.Percentages[i] != expected {
				t.Errorf("Set %d: expected %.1f%%, got %.1f%%", i+1, expected, entry.Percentages[i])
			}
		}
	})

	t.Run("has correct reps", func(t *testing.T) {
		entry := wl.GetByWeekNumber(1)
		expectedReps := []int{6, 4, 2}

		if len(entry.Reps) != len(expectedReps) {
			t.Fatalf("Expected %d rep values, got %d", len(expectedReps), len(entry.Reps))
		}

		for i, expected := range expectedReps {
			if entry.Reps[i] != expected {
				t.Errorf("Set %d: expected %d reps, got %d reps", i+1, expected, entry.Reps[i])
			}
		}
	})
}

// TestCAP3LookupContextPercentageApplication tests that percentages are correctly applied per set.
func TestCAP3LookupContextPercentageApplication(t *testing.T) {
	wl := createCAP3HighIntensityWeeklyLookup()

	testCases := []struct {
		setNumber          int
		basePercentage     float64
		expectedPercentage float64
	}{
		{1, 100.0, 79.5}, // Set 1: replaces with 79.5%
		{2, 100.0, 83.5}, // Set 2: replaces with 83.5%
		{3, 100.0, 88.5}, // Set 3: replaces with 88.5%
	}

	for _, tc := range testCases {
		t.Run("set "+string(rune('0'+tc.setNumber)), func(t *testing.T) {
			ctx := &loadstrategy.LookupContext{
				WeekNumber:   1,
				SetNumber:    tc.setNumber,
				WeeklyLookup: wl,
			}

			result := ctx.ApplyModifiers(tc.basePercentage)
			if result != tc.expectedPercentage {
				t.Errorf("Expected %.1f%%, got %.1f%%", tc.expectedPercentage, result)
			}
		})
	}
}

// TestCAP3LookupContextRepsForSet tests that target reps are correctly retrieved per set.
func TestCAP3LookupContextRepsForSet(t *testing.T) {
	wl := createCAP3HighIntensityWeeklyLookup()

	testCases := []struct {
		setNumber    int
		expectedReps int
	}{
		{1, 6}, // Set 1: 6 reps
		{2, 4}, // Set 2: 4 reps
		{3, 2}, // Set 3: 2+ reps (AMRAP)
	}

	for _, tc := range testCases {
		t.Run("set "+string(rune('0'+tc.setNumber)), func(t *testing.T) {
			ctx := &loadstrategy.LookupContext{
				WeekNumber:   1,
				SetNumber:    tc.setNumber,
				WeeklyLookup: wl,
			}

			reps := ctx.GetRepsForSet()
			if reps != tc.expectedReps {
				t.Errorf("Expected %d reps, got %d", tc.expectedReps, reps)
			}
		})
	}
}

// TestCAP3UserProgramStateRotation tests rotation advancement in UserProgramState.
func TestCAP3UserProgramStateRotation(t *testing.T) {
	userID := uuid.New().String()
	programID := uuid.New().String()

	state, _ := userprogramstate.EnrollUser(
		userprogramstate.EnrollUserInput{
			UserID:    userID,
			ProgramID: programID,
		},
		uuid.New().String(),
	)

	rotationLength := 3 // CAP3 has 3 positions

	t.Run("initial position is 0", func(t *testing.T) {
		if state.RotationPosition != 0 {
			t.Errorf("Expected initial rotation position 0, got %d", state.RotationPosition)
		}
	})

	t.Run("advances to position 1", func(t *testing.T) {
		state.AdvanceRotation(rotationLength)
		if state.RotationPosition != 1 {
			t.Errorf("Expected rotation position 1, got %d", state.RotationPosition)
		}
	})

	t.Run("advances to position 2", func(t *testing.T) {
		state.AdvanceRotation(rotationLength)
		if state.RotationPosition != 2 {
			t.Errorf("Expected rotation position 2, got %d", state.RotationPosition)
		}
	})

	t.Run("wraps around to position 0", func(t *testing.T) {
		state.AdvanceRotation(rotationLength)
		if state.RotationPosition != 0 {
			t.Errorf("Expected rotation position to wrap to 0, got %d", state.RotationPosition)
		}
	})
}

// TestCAP3FullRotationCycle tests a complete 3-week rotation cycle.
func TestCAP3FullRotationCycle(t *testing.T) {
	rl := createCAP3RotationLookup()
	highIntensityLookup := createCAP3HighIntensityWeeklyLookup()
	mediumLookup := createCAP3MediumWeeklyLookup()
	volumeLookup := createCAP3VolumeWeeklyLookup()

	// User enrolled in CAP3
	userID := uuid.New().String()
	programID := uuid.New().String()
	state, _ := userprogramstate.EnrollUser(
		userprogramstate.EnrollUserInput{
			UserID:    userID,
			ProgramID: programID,
		},
		uuid.New().String(),
	)

	// Simulate 3-week rotation cycle
	weekConfigs := []struct {
		weekNum            int
		focusLift          string
		deadliftLookup     *weeklylookup.WeeklyLookup
		squatLookup        *weeklylookup.WeeklyLookup
		benchLookup        *weeklylookup.WeeklyLookup
	}{
		// Week 1: Deadlift High Intensity, Squat Medium, Bench Volume
		{1, CAP3LiftDeadlift, highIntensityLookup, mediumLookup, volumeLookup},
		// Week 2: Squat High Intensity, Deadlift Medium, Bench Medium
		{2, CAP3LiftSquat, mediumLookup, highIntensityLookup, mediumLookup},
		// Week 3: Bench High Intensity, Deadlift Volume, Squat Volume
		{3, CAP3LiftBench, volumeLookup, volumeLookup, highIntensityLookup},
	}

	for weekIdx, config := range weekConfigs {
		t.Run("week "+string(rune('0'+config.weekNum)), func(t *testing.T) {
			// Verify rotation position matches expected focus lift
			ctx := &loadstrategy.LookupContext{
				RotationPosition: state.RotationPosition,
				RotationLookup:   rl,
			}

			if !ctx.IsLiftInRotationFocus(config.focusLift) {
				t.Errorf("Expected %s to be in focus for week %d", config.focusLift, config.weekNum)
			}

			// Verify deadlift percentage type
			deadliftInFocus := ctx.IsLiftInRotationFocus(CAP3LiftDeadlift)
			if deadliftInFocus && config.deadliftLookup != highIntensityLookup {
				t.Error("Deadlift should use high intensity lookup when in focus")
			}

			// Verify squat percentage type
			squatInFocus := ctx.IsLiftInRotationFocus(CAP3LiftSquat)
			if squatInFocus && config.squatLookup != highIntensityLookup {
				t.Error("Squat should use high intensity lookup when in focus")
			}

			// Verify bench percentage type
			benchInFocus := ctx.IsLiftInRotationFocus(CAP3LiftBench)
			if benchInFocus && config.benchLookup != highIntensityLookup {
				t.Error("Bench should use high intensity lookup when in focus")
			}

			// Advance rotation after each week (in real program, this happens after cycle completes)
			if weekIdx < len(weekConfigs)-1 {
				state.AdvanceRotation(rl.Length())
			}
		})
	}

	// Verify we've completed one rotation cycle
	t.Run("rotation returns to position 0 after cycle", func(t *testing.T) {
		state.AdvanceRotation(rl.Length())
		if state.RotationPosition != 0 {
			t.Errorf("Expected rotation to return to 0, got %d", state.RotationPosition)
		}
	})
}

// TestCAP3AMRAPSetSchemeIntegration tests AMRAP set scheme with CAP3 percentages.
func TestCAP3AMRAPSetSchemeIntegration(t *testing.T) {
	// Create AMRAP scheme for the final high intensity set (2+ reps)
	scheme, err := setscheme.NewAMRAPSetScheme(1, 2)
	if err != nil {
		t.Fatalf("Failed to create AMRAP scheme: %v", err)
	}

	t.Run("generates single AMRAP set", func(t *testing.T) {
		baseWeight := 100.0 // Assume 100lb at 88.5% TM

		sets, err := scheme.GenerateSets(baseWeight, setscheme.DefaultSetGenerationContext())
		if err != nil {
			t.Fatalf("Failed to generate sets: %v", err)
		}

		if len(sets) != 1 {
			t.Errorf("Expected 1 set, got %d", len(sets))
		}

		if sets[0].TargetReps != 2 {
			t.Errorf("Expected target reps 2, got %d", sets[0].TargetReps)
		}

		if sets[0].Weight != baseWeight {
			t.Errorf("Expected weight %.1f, got %.1f", baseWeight, sets[0].Weight)
		}

		if !sets[0].IsWorkSet {
			t.Error("AMRAP set should be marked as work set")
		}
	})
}

// mockLiftLookup implements prescription.LiftLookup for testing.
type mockLiftLookup struct {
	lifts map[string]*prescription.LiftInfo
}

func (m *mockLiftLookup) GetLiftByID(_ context.Context, liftID string) (*prescription.LiftInfo, error) {
	if lift, ok := m.lifts[liftID]; ok {
		return lift, nil
	}
	return nil, nil
}

// mockMaxLookup implements loadstrategy.MaxLookup for testing.
type mockMaxLookup struct {
	maxes map[string]float64 // liftID -> max value
}

func (m *mockMaxLookup) GetCurrentMax(_ context.Context, _ string, liftID string, _ string) (*loadstrategy.MaxValue, error) {
	if max, ok := m.maxes[liftID]; ok {
		return &loadstrategy.MaxValue{Value: max}, nil
	}
	return nil, loadstrategy.ErrMaxNotFound
}

// TestCAP3PrescriptionResolutionWithRotation tests prescription resolution integrating rotation context.
func TestCAP3PrescriptionResolutionWithRotation(t *testing.T) {
	// Setup: Create lift IDs and mocks
	deadliftID := uuid.New().String()
	squatID := uuid.New().String()
	benchID := uuid.New().String()

	liftLookup := &mockLiftLookup{
		lifts: map[string]*prescription.LiftInfo{
			deadliftID: {ID: deadliftID, Name: "Deadlift", Slug: "deadlift"},
			squatID:    {ID: squatID, Name: "Squat", Slug: "squat"},
			benchID:    {ID: benchID, Name: "Bench Press", Slug: "bench"},
		},
	}

	maxLookup := &mockMaxLookup{
		maxes: map[string]float64{
			deadliftID: 200.0, // 200lb TM
			squatID:    180.0, // 180lb TM
			benchID:    150.0, // 150lb TM
		},
	}

	rl := createCAP3RotationLookup()
	highIntensityLookup := createCAP3HighIntensityWeeklyLookup()

	userID := uuid.New().String()

	// Test Week 1: Deadlift focus with high intensity
	t.Run("Week 1 - Deadlift at High Intensity", func(t *testing.T) {
		// Create prescription for deadlift at high intensity
		loadStrategy := loadstrategy.NewPercentOfLoadStrategy(
			loadstrategy.ReferenceTrainingMax,
			100.0, // Base 100%, will be modified by lookup
			5.0,   // Rounding increment
			loadstrategy.RoundNearest,
			maxLookup,
		)

		amrapScheme, _ := setscheme.NewAMRAPSetScheme(1, 2)

		prx, result := prescription.CreatePrescription(
			prescription.CreatePrescriptionInput{
				LiftID:       deadliftID,
				LoadStrategy: loadStrategy,
				SetScheme:    amrapScheme,
				Order:        1,
				Notes:        "High Intensity AMRAP Day",
			},
			uuid.New().String(),
		)
		if !result.Valid {
			t.Fatalf("Failed to create prescription: %v", result.Errors)
		}

		// Create lookup context for Week 1, Set 3 (88.5% AMRAP)
		lookupCtx := &loadstrategy.LookupContext{
			WeekNumber:       1,
			SetNumber:        3, // Final AMRAP set
			WeeklyLookup:     highIntensityLookup,
			RotationPosition: CAP3PositionDeadlift,
			RotationLookup:   rl,
		}

		resCtx := prescription.ResolutionContext{
			LiftLookup:    liftLookup,
			SetGenContext: setscheme.DefaultSetGenerationContext(),
			LookupContext: lookupCtx,
		}

		resolved, err := prx.Resolve(context.Background(), userID, resCtx)
		if err != nil {
			t.Fatalf("Failed to resolve prescription: %v", err)
		}

		// Verify the lift is in focus
		if !lookupCtx.IsLiftInRotationFocus(CAP3LiftDeadlift) {
			t.Error("Deadlift should be in focus for Week 1")
		}

		// Verify set was generated
		if len(resolved.Sets) != 1 {
			t.Errorf("Expected 1 set, got %d", len(resolved.Sets))
		}

		// Weight should be 200 * 88.5% = 177, rounded to nearest 5 = 175
		expectedWeight := 175.0 // 200 * 0.885 = 177 → rounded to 175
		if resolved.Sets[0].Weight != expectedWeight {
			t.Errorf("Expected weight %.1f (88.5%% of 200, rounded to nearest 5), got %.1f", expectedWeight, resolved.Sets[0].Weight)
		}
	})

	// Test that squat is NOT in focus during Week 1
	t.Run("Week 1 - Squat not in focus", func(t *testing.T) {
		lookupCtx := &loadstrategy.LookupContext{
			RotationPosition: CAP3PositionDeadlift,
			RotationLookup:   rl,
		}

		if lookupCtx.IsLiftInRotationFocus(CAP3LiftSquat) {
			t.Error("Squat should NOT be in focus for Week 1")
		}
	})

	// Test Week 2: Squat in focus
	t.Run("Week 2 - Squat is in focus", func(t *testing.T) {
		lookupCtx := &loadstrategy.LookupContext{
			RotationPosition: CAP3PositionSquat,
			RotationLookup:   rl,
		}

		if !lookupCtx.IsLiftInRotationFocus(CAP3LiftSquat) {
			t.Error("Squat should be in focus for Week 2")
		}

		if lookupCtx.IsLiftInRotationFocus(CAP3LiftDeadlift) {
			t.Error("Deadlift should NOT be in focus for Week 2")
		}

		if lookupCtx.IsLiftInRotationFocus(CAP3LiftBench) {
			t.Error("Bench should NOT be in focus for Week 2")
		}
	})

	// Test Week 3: Bench in focus
	t.Run("Week 3 - Bench is in focus", func(t *testing.T) {
		lookupCtx := &loadstrategy.LookupContext{
			RotationPosition: CAP3PositionBench,
			RotationLookup:   rl,
		}

		if !lookupCtx.IsLiftInRotationFocus(CAP3LiftBench) {
			t.Error("Bench should be in focus for Week 3")
		}

		if lookupCtx.IsLiftInRotationFocus(CAP3LiftDeadlift) {
			t.Error("Deadlift should NOT be in focus for Week 3")
		}

		if lookupCtx.IsLiftInRotationFocus(CAP3LiftSquat) {
			t.Error("Squat should NOT be in focus for Week 3")
		}
	})
}

// TestCAP3VolumeAndMediumDays tests the volume and medium day configurations.
func TestCAP3VolumeAndMediumDays(t *testing.T) {
	volumeLookup := createCAP3VolumeWeeklyLookup()
	mediumLookup := createCAP3MediumWeeklyLookup()

	t.Run("volume day has 7 sets at 73.5%", func(t *testing.T) {
		entry := volumeLookup.GetByWeekNumber(1)
		if entry == nil {
			t.Fatal("Volume entry is nil")
		}

		if len(entry.Percentages) != 7 {
			t.Errorf("Expected 7 sets for volume day, got %d", len(entry.Percentages))
		}

		for i, pct := range entry.Percentages {
			if pct != 73.5 {
				t.Errorf("Set %d: expected 73.5%%, got %.1f%%", i+1, pct)
			}
		}

		for i, reps := range entry.Reps {
			if reps != 4 {
				t.Errorf("Set %d: expected 4 reps, got %d", i+1, reps)
			}
		}
	})

	t.Run("medium day has 8 sets at 77%", func(t *testing.T) {
		entry := mediumLookup.GetByWeekNumber(1)
		if entry == nil {
			t.Fatal("Medium entry is nil")
		}

		if len(entry.Percentages) != 8 {
			t.Errorf("Expected 8 sets for medium day, got %d", len(entry.Percentages))
		}

		for i, pct := range entry.Percentages {
			if pct != 77 {
				t.Errorf("Set %d: expected 77%%, got %.1f%%", i+1, pct)
			}
		}

		for i, reps := range entry.Reps {
			if reps != 3 {
				t.Errorf("Set %d: expected 3 reps, got %d", i+1, reps)
			}
		}
	})
}

// TestCAP3CyclesSinceStartTracking tests that cycles are properly counted.
func TestCAP3CyclesSinceStartTracking(t *testing.T) {
	userID := uuid.New().String()
	programID := uuid.New().String()

	state, _ := userprogramstate.EnrollUser(
		userprogramstate.EnrollUserInput{
			UserID:    userID,
			ProgramID: programID,
		},
		uuid.New().String(),
	)

	t.Run("initial cycles is 0", func(t *testing.T) {
		if state.CyclesSinceStart != 0 {
			t.Errorf("Expected 0 cycles initially, got %d", state.CyclesSinceStart)
		}
	})

	// Simulate completing 3-week cycle (CAP3 cycle length = 3 weeks)
	ctx := userprogramstate.AdvancementContext{
		DaysInCurrentWeek: 6, // 6 training days per week
		CycleLengthWeeks:  3, // CAP3 is 3 weeks
	}

	// Advance through all days of week 1-3
	for week := 1; week <= 3; week++ {
		for day := 0; day < 6; day++ {
			result, valResult := userprogramstate.AdvanceState(state, ctx)
			if !valResult.Valid {
				t.Fatalf("Failed to advance state: %v", valResult.Errors)
			}
			state = result.NewState
		}
	}

	t.Run("cycles incremented after 3 weeks", func(t *testing.T) {
		if state.CyclesSinceStart != 1 {
			t.Errorf("Expected 1 cycle after 3 weeks, got %d", state.CyclesSinceStart)
		}
	})

	// Complete another full cycle
	for week := 1; week <= 3; week++ {
		for day := 0; day < 6; day++ {
			result, valResult := userprogramstate.AdvanceState(state, ctx)
			if !valResult.Valid {
				t.Fatalf("Failed to advance state: %v", valResult.Errors)
			}
			state = result.NewState
		}
	}

	t.Run("cycles is 2 after second cycle", func(t *testing.T) {
		if state.CyclesSinceStart != 2 {
			t.Errorf("Expected 2 cycles after 6 weeks, got %d", state.CyclesSinceStart)
		}
	})
}

// TestCAP3RotationPositionIndependentOfWeek tests that rotation position and week number are independent.
func TestCAP3RotationPositionIndependentOfWeek(t *testing.T) {
	rl := createCAP3RotationLookup()

	// In CAP3, rotation position determines WHICH lift is in focus
	// Week number within the lookup determines SET STRUCTURE (not focus)
	// These are independent concerns that work together

	t.Run("rotation position determines lift focus", func(t *testing.T) {
		for pos := 0; pos < 3; pos++ {
			ctx := &loadstrategy.LookupContext{
				RotationPosition: pos,
				RotationLookup:   rl,
			}

			entry := ctx.GetRotationEntry()
			if entry == nil {
				t.Errorf("No entry for position %d", pos)
				continue
			}

			expectedLifts := []string{CAP3LiftDeadlift, CAP3LiftSquat, CAP3LiftBench}
			if entry.LiftIdentifier != expectedLifts[pos] {
				t.Errorf("Position %d: expected %s, got %s", pos, expectedLifts[pos], entry.LiftIdentifier)
			}
		}
	})

	t.Run("week number determines percentage structure within lookup", func(t *testing.T) {
		highIntensityLookup := createCAP3HighIntensityWeeklyLookup()

		// Week 1 of the high intensity lookup has specific percentages
		ctx := &loadstrategy.LookupContext{
			WeekNumber:   1,
			SetNumber:    1,
			WeeklyLookup: highIntensityLookup,
		}

		pct := ctx.ApplyModifiers(100.0)
		if pct != 79.5 {
			t.Errorf("Expected 79.5%% for set 1, got %.1f%%", pct)
		}
	})
}
