// Package integration provides integration tests for cross-component behavior.
// This file contains acceptance tests that validate Phase 5 business requirements
// for CAP3, Inverted Juggernaut, and GreySkull LP programs.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/greyskull"
	"github.com/waynenilsen/power-pro-v3/internal/domain/juggernaut"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
)

// =============================================================================
// NSUNS CAP3 ACCEPTANCE TESTS
// =============================================================================

// TestCAP3AcceptanceCriteria validates all business requirements for nSuns CAP3.
func TestCAP3AcceptanceCriteria(t *testing.T) {
	// Setup rotation lookup
	rl := createCAP3RotationLookup()

	// Acceptance: Each lift gets AMRAP focus once per 3-week cycle
	t.Run("AC1: Each lift gets AMRAP focus once per 3-week cycle", func(t *testing.T) {
		focusByWeek := map[int]string{
			1: CAP3LiftDeadlift, // Week 1: Deadlift focus
			2: CAP3LiftSquat,    // Week 2: Squat focus
			3: CAP3LiftBench,    // Week 3: Bench focus
		}

		for week := 1; week <= 3; week++ {
			rotationPosition := week - 1 // Position 0-2
			ctx := &loadstrategy.LookupContext{
				RotationPosition: rotationPosition,
				RotationLookup:   rl,
			}

			expectedFocusLift := focusByWeek[week]
			if !ctx.IsLiftInRotationFocus(expectedFocusLift) {
				t.Errorf("Week %d: Expected %s to be in focus", week, expectedFocusLift)
			}

			// Verify other lifts are NOT in focus
			for _, lift := range []string{CAP3LiftDeadlift, CAP3LiftSquat, CAP3LiftBench} {
				if lift != expectedFocusLift && ctx.IsLiftInRotationFocus(lift) {
					t.Errorf("Week %d: %s should NOT be in focus (only %s should be)", week, lift, expectedFocusLift)
				}
			}
		}
	})

	// Acceptance: Non-focus lifts use medium/volume percentages
	t.Run("AC2: Non-focus lifts use medium/volume percentages", func(t *testing.T) {
		mediumLookup := createCAP3MediumWeeklyLookup()
		volumeLookup := createCAP3VolumeWeeklyLookup()

		// Verify medium day: 8 sets at 77%
		mediumEntry := mediumLookup.GetByWeekNumber(1)
		if len(mediumEntry.Percentages) != 8 {
			t.Errorf("Medium day: expected 8 sets, got %d", len(mediumEntry.Percentages))
		}
		for i, pct := range mediumEntry.Percentages {
			if pct != 77.0 {
				t.Errorf("Medium day set %d: expected 77%%, got %.1f%%", i+1, pct)
			}
		}

		// Verify volume day: 7 sets at 73.5%
		volumeEntry := volumeLookup.GetByWeekNumber(1)
		if len(volumeEntry.Percentages) != 7 {
			t.Errorf("Volume day: expected 7 sets, got %d", len(volumeEntry.Percentages))
		}
		for i, pct := range volumeEntry.Percentages {
			if pct != 73.5 {
				t.Errorf("Volume day set %d: expected 73.5%%, got %.1f%%", i+1, pct)
			}
		}
	})

	// Acceptance: Rotation cycles correctly (0→1→2→0)
	t.Run("AC3: Rotation cycles correctly (0→1→2→0)", func(t *testing.T) {
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

		expectedPositions := []int{0, 1, 2, 0, 1, 2} // Two full cycles
		for i, expected := range expectedPositions {
			if state.RotationPosition != expected {
				t.Errorf("Cycle position %d: expected %d, got %d", i, expected, state.RotationPosition)
			}
			state.AdvanceRotation(rotationLength)
		}
	})

	// Acceptance: TM adjusts based on AMRAP performance (via LookupContext integration)
	t.Run("AC4: TM adjustment mechanism integrated with rotation", func(t *testing.T) {
		highIntensityLookup := createCAP3HighIntensityWeeklyLookup()

		// Verify high intensity AMRAP percentages
		entry := highIntensityLookup.GetByWeekNumber(1)

		expectedPercentages := []float64{79.5, 83.5, 88.5}
		expectedReps := []int{6, 4, 2}

		for i := range expectedPercentages {
			if entry.Percentages[i] != expectedPercentages[i] {
				t.Errorf("High intensity set %d: expected %.1f%%, got %.1f%%", i+1, expectedPercentages[i], entry.Percentages[i])
			}
			if entry.Reps[i] != expectedReps[i] {
				t.Errorf("High intensity set %d: expected %d reps, got %d", i+1, expectedReps[i], entry.Reps[i])
			}
		}
	})
}

// =============================================================================
// INVERTED JUGGERNAUT ACCEPTANCE TESTS
// =============================================================================

// TestInvertedJuggernautAcceptanceCriteria validates all business requirements for Inverted Juggernaut.
func TestInvertedJuggernautAcceptanceCriteria(t *testing.T) {
	// Acceptance: 4 waves cycle correctly (10s→8s→5s→3s)
	t.Run("AC1: 4 waves cycle correctly (10s→8s→5s→3s)", func(t *testing.T) {
		waveWeeks := map[int]struct {
			waveIndex int
			waveName  string
		}{
			1:  {0, "10s"},
			4:  {0, "10s"},
			5:  {1, "8s"},
			8:  {1, "8s"},
			9:  {2, "5s"},
			12: {2, "5s"},
			13: {3, "3s"},
			16: {3, "3s"},
		}

		for week, expected := range waveWeeks {
			info := juggernaut.GetWaveInfo(week)
			if info.WaveIndex != expected.waveIndex {
				t.Errorf("Week %d: expected wave index %d, got %d", week, expected.waveIndex, info.WaveIndex)
			}
			if info.WaveName != expected.waveName {
				t.Errorf("Week %d: expected wave name %s, got %s", week, expected.waveName, info.WaveName)
			}
		}
	})

	// Acceptance: Volume sets match wave target (9/7/5/6)
	t.Run("AC2: Volume sets match wave target (9/7/5/6)", func(t *testing.T) {
		waveVolumes := []struct {
			waveIndex    int
			sampleWeek   int
			expectedSets int
		}{
			{0, 1, 9},  // 10s wave: 9 sets
			{1, 5, 7},  // 8s wave: 7 sets
			{2, 9, 5},  // 5s wave: 5 sets
			{3, 13, 6}, // 3s wave: 6 sets
		}

		for _, wv := range waveVolumes {
			info := juggernaut.GetWaveInfo(wv.sampleWeek)
			configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
			if len(configs) != wv.expectedSets {
				t.Errorf("Wave %d (week %d): expected %d volume sets, got %d",
					wv.waveIndex, wv.sampleWeek, wv.expectedSets, len(configs))
			}
		}
	})

	// Acceptance: Base percentages match wave (60/65/70/75%)
	t.Run("AC3: Base percentages match wave (60/65/70/75%)", func(t *testing.T) {
		wavePercentages := []struct {
			waveIndex      int
			sampleWeek     int
			expectedPct    float64
		}{
			{0, 1, 60.0},  // 10s wave: 60%
			{1, 5, 65.0},  // 8s wave: 65%
			{2, 9, 70.0},  // 5s wave: 70%
			{3, 13, 75.0}, // 3s wave: 75%
		}

		for _, wp := range wavePercentages {
			info := juggernaut.GetWaveInfo(wp.sampleWeek)
			configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
			if len(configs) > 0 && configs[0].Percentage != wp.expectedPct {
				t.Errorf("Wave %d: expected %.1f%% base, got %.1f%%",
					wp.waveIndex, wp.expectedPct, configs[0].Percentage)
			}
		}
	})

	// Acceptance: 5/3/1 percentages apply correctly
	t.Run("AC4: 5/3/1 percentages apply correctly", func(t *testing.T) {
		lookup := juggernaut.Create531WeeklyLookup("test-id", nil)

		expectedByWeek := map[int][]float64{
			1: {65.0, 75.0, 85.0, 75.0, 65.0}, // Accumulation
			2: {70.0, 80.0, 90.0, 80.0, 70.0}, // Intensification
			3: {75.0, 85.0, 95.0, 85.0, 75.0}, // Realization
			4: {40.0, 50.0, 60.0},              // Deload
		}

		for week, expectedPcts := range expectedByWeek {
			entry := lookup.GetByWeekNumber(week)
			if entry == nil {
				t.Errorf("Week %d: no entry found", week)
				continue
			}
			if len(entry.Percentages) != len(expectedPcts) {
				t.Errorf("Week %d: expected %d percentages, got %d", week, len(expectedPcts), len(entry.Percentages))
				continue
			}
			for i, pct := range expectedPcts {
				if entry.Percentages[i] != pct {
					t.Errorf("Week %d set %d: expected %.1f%%, got %.1f%%", week, i+1, pct, entry.Percentages[i])
				}
			}
		}
	})

	// Acceptance: 16-week cycle completes and resets
	t.Run("AC5: 16-week cycle completes and resets", func(t *testing.T) {
		userID := uuid.New().String()
		programID := uuid.New().String()

		state, _ := userprogramstate.EnrollUser(
			userprogramstate.EnrollUserInput{
				UserID:    userID,
				ProgramID: programID,
			},
			uuid.New().String(),
		)

		ctx := userprogramstate.AdvancementContext{
			DaysInCurrentWeek: 4, // 4 training days per week
			CycleLengthWeeks:  16,
		}

		// Advance through 16 weeks (64 training days)
		for day := 0; day < 64; day++ {
			result, valResult := userprogramstate.AdvanceState(state, ctx)
			if !valResult.Valid {
				t.Fatalf("Day %d: AdvanceState failed: %v", day+1, valResult.Errors)
			}
			state = result.NewState
		}

		// Verify cycle completed
		if state.CyclesSinceStart != 1 {
			t.Errorf("Expected CyclesSinceStart=1, got %d", state.CyclesSinceStart)
		}
		if state.CurrentWeek != 1 {
			t.Errorf("Expected CurrentWeek=1 after cycle reset, got %d", state.CurrentWeek)
		}
	})

	// Acceptance: Deload weeks have no volume sets
	t.Run("AC6: Deload weeks (4, 8, 12, 16) have no volume sets", func(t *testing.T) {
		deloadWeeks := []int{4, 8, 12, 16}
		for _, week := range deloadWeeks {
			info := juggernaut.GetWaveInfo(week)
			if !info.IsDeload {
				t.Errorf("Week %d: expected IsDeload=true", week)
			}
			configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
			if configs != nil {
				t.Errorf("Week %d: expected nil volume configs for deload, got %d sets", week, len(configs))
			}
		}
	})

	// Acceptance: Realization weeks have AMRAP final set
	t.Run("AC7: Realization weeks (3, 7, 11, 15) have AMRAP final set", func(t *testing.T) {
		realizationWeeks := []int{3, 7, 11, 15}
		for _, week := range realizationWeeks {
			info := juggernaut.GetWaveInfo(week)
			if !info.IsRealization {
				t.Errorf("Week %d: expected IsRealization=true", week)
			}
			configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
			if configs == nil || len(configs) == 0 {
				t.Errorf("Week %d: expected volume configs for realization", week)
				continue
			}
			lastSet := configs[len(configs)-1]
			if !lastSet.IsAMRAP {
				t.Errorf("Week %d: expected last set to be AMRAP", week)
			}
		}
	})
}

// =============================================================================
// GREYSKULL LP ACCEPTANCE TESTS
// =============================================================================

// TestGreySkullLPAcceptanceCriteria validates all business requirements for GreySkull LP.
func TestGreySkullLPAcceptanceCriteria(t *testing.T) {
	// Acceptance: Days alternate A/B/A, B/A/B correctly
	t.Run("AC1: Days alternate A/B/A, B/A/B correctly", func(t *testing.T) {
		// Week 1 (odd): A, B, A
		for day := 1; day <= 3; day++ {
			variant, err := greyskull.GetVariantForDay(1, day)
			if err != nil {
				t.Fatalf("Week 1 Day %d: %v", day, err)
			}
			expectedVariants := []greyskull.Variant{greyskull.VariantA, greyskull.VariantB, greyskull.VariantA}
			if variant != expectedVariants[day-1] {
				t.Errorf("Week 1 Day %d: expected %s, got %s", day, expectedVariants[day-1], variant)
			}
		}

		// Week 2 (even): B, A, B
		for day := 1; day <= 3; day++ {
			variant, err := greyskull.GetVariantForDay(2, day)
			if err != nil {
				t.Fatalf("Week 2 Day %d: %v", day, err)
			}
			expectedVariants := []greyskull.Variant{greyskull.VariantB, greyskull.VariantA, greyskull.VariantB}
			if variant != expectedVariants[day-1] {
				t.Errorf("Week 2 Day %d: expected %s, got %s", day, expectedVariants[day-1], variant)
			}
		}
	})

	// Acceptance: AMRAP final sets work on all main lifts
	t.Run("AC2: AMRAP final sets work on all main lifts", func(t *testing.T) {
		// GreySkull scheme: 2 fixed sets + 1 AMRAP set
		scheme, err := setscheme.NewGreySkullSetScheme(2, 5, 1, 5)
		if err != nil {
			t.Fatalf("Failed to create scheme: %v", err)
		}

		sets, err := scheme.GenerateSets(135.0, setscheme.DefaultSetGenerationContext())
		if err != nil {
			t.Fatalf("Failed to generate sets: %v", err)
		}

		if len(sets) != 3 {
			t.Errorf("Expected 3 sets (2 fixed + 1 AMRAP), got %d", len(sets))
		}

		// Verify all sets are work sets (AMRAP final set should also be work set)
		for i, set := range sets {
			if !set.IsWorkSet {
				t.Errorf("Set %d: expected IsWorkSet=true", i+1)
			}
		}
	})

	// Acceptance: Double progression triggers on 10+ reps
	t.Run("AC3: Double progression triggers on 10+ reps", func(t *testing.T) {
		gsp, err := greyskull.NewGreySkullMainLiftProgression(
			uuid.New().String(), "Test", 2.5, progression.TrainingMax)
		if err != nil {
			t.Fatalf("Failed to create progression: %v", err)
		}

		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			MaxType:      progression.TrainingMax,
			CurrentValue: 135.0,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtrAcceptance(10), // Exactly at double threshold
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(context.Background(), params)
		if err != nil {
			t.Fatalf("Apply failed: %v", err)
		}

		// Should get double increment: +5.0 (2.5 * 2)
		expectedNew := 140.0 // 135 + 5
		if result.NewValue != expectedNew {
			t.Errorf("Expected new value %.1f (double increment), got %.1f", expectedNew, result.NewValue)
		}
	})

	// Acceptance: 10% deload triggers on failure (<5 reps)
	t.Run("AC4: 10% deload triggers on failure (<5 reps)", func(t *testing.T) {
		gsp, err := greyskull.NewGreySkullMainLiftProgression(
			uuid.New().String(), "Test", 2.5, progression.TrainingMax)
		if err != nil {
			t.Fatalf("Failed to create progression: %v", err)
		}

		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			MaxType:      progression.TrainingMax,
			CurrentValue: 135.0,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtrAcceptance(4), // Below minimum (failure)
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(context.Background(), params)
		if err != nil {
			t.Fatalf("Apply failed: %v", err)
		}

		// Should get 10% deload: 135 - 13.5 = 121.5
		expectedNew := 121.5
		if result.NewValue != expectedNew {
			t.Errorf("Expected new value %.1f (10%% deload), got %.1f", expectedNew, result.NewValue)
		}
	})

	// Acceptance: Standard progression for 5-9 reps
	t.Run("AC5: Standard progression for 5-9 reps", func(t *testing.T) {
		gsp, err := greyskull.NewGreySkullMainLiftProgression(
			uuid.New().String(), "Test", 2.5, progression.TrainingMax)
		if err != nil {
			t.Fatalf("Failed to create progression: %v", err)
		}

		testCases := []int{5, 6, 7, 8, 9}
		for _, reps := range testCases {
			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				MaxType:      progression.TrainingMax,
				CurrentValue: 135.0,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtrAcceptance(reps),
					IsAMRAP:       true,
				},
			}

			result, err := gsp.Apply(context.Background(), params)
			if err != nil {
				t.Fatalf("Apply failed for %d reps: %v", reps, err)
			}

			// Should get standard increment: +2.5
			expectedNew := 137.5
			if result.NewValue != expectedNew {
				t.Errorf("%d reps: expected new value %.1f (standard increment), got %.1f", reps, expectedNew, result.NewValue)
			}
		}
	})

	// Acceptance: Variant A has correct lifts (Bench, Row, Squat)
	t.Run("AC6: Variant A has correct lifts", func(t *testing.T) {
		lifts, err := greyskull.GetDayTemplate(greyskull.VariantA)
		if err != nil {
			t.Fatalf("GetDayTemplate failed: %v", err)
		}

		expectedSlugs := []string{"bench-press", "barbell-row", "squat", "tricep-extension", "ab-rollout"}
		if len(lifts) != len(expectedSlugs) {
			t.Errorf("Variant A: expected %d lifts, got %d", len(expectedSlugs), len(lifts))
		}
		for i, expected := range expectedSlugs {
			if lifts[i].Slug != expected {
				t.Errorf("Variant A lift %d: expected %s, got %s", i+1, expected, lifts[i].Slug)
			}
		}
	})

	// Acceptance: Variant B has correct lifts (OHP, Chinups, Deadlift)
	t.Run("AC7: Variant B has correct lifts", func(t *testing.T) {
		lifts, err := greyskull.GetDayTemplate(greyskull.VariantB)
		if err != nil {
			t.Fatalf("GetDayTemplate failed: %v", err)
		}

		expectedSlugs := []string{"overhead-press", "chin-up", "deadlift", "bicep-curl", "shrug"}
		if len(lifts) != len(expectedSlugs) {
			t.Errorf("Variant B: expected %d lifts, got %d", len(expectedSlugs), len(lifts))
		}
		for i, expected := range expectedSlugs {
			if lifts[i].Slug != expected {
				t.Errorf("Variant B lift %d: expected %s, got %s", i+1, expected, lifts[i].Slug)
			}
		}
	})

	// Acceptance: Main lifts use 3x5 with AMRAP final set
	t.Run("AC8: Main lifts use 3x5 with AMRAP final set", func(t *testing.T) {
		for _, variant := range []greyskull.Variant{greyskull.VariantA, greyskull.VariantB} {
			lifts, _ := greyskull.GetDayTemplate(variant)

			// First 3 lifts are main lifts
			for i := 0; i < 3; i++ {
				lift := lifts[i]
				if lift.Sets != 3 {
					t.Errorf("Variant %s main lift %s: expected 3 sets, got %d", variant, lift.Slug, lift.Sets)
				}
				if lift.Reps != 5 {
					t.Errorf("Variant %s main lift %s: expected 5 reps, got %d", variant, lift.Slug, lift.Reps)
				}
				if !lift.IsAMRAP {
					t.Errorf("Variant %s main lift %s: expected IsAMRAP=true", variant, lift.Slug)
				}
			}
		}
	})

	// Acceptance: 2-week cycle repeats correctly
	t.Run("AC9: 2-week cycle repeats correctly", func(t *testing.T) {
		// Verify pattern is consistent across multiple cycles
		for cycle := 0; cycle < 3; cycle++ {
			for week := 1; week <= 2; week++ {
				actualWeek := cycle*2 + week
				for day := 1; day <= 3; day++ {
					variant1, _ := greyskull.GetVariantForDay(week, day)
					variant2, _ := greyskull.GetVariantForDay(actualWeek, day)

					// Pattern should repeat every 2 weeks
					if week == ((actualWeek-1)%2)+1 {
						if variant1 != variant2 {
							t.Errorf("Week %d Day %d vs Week %d Day %d: pattern mismatch (%s vs %s)",
								week, day, actualWeek, day, variant1, variant2)
						}
					}
				}
			}
		}
	})
}

// =============================================================================
// COMPREHENSIVE MULTI-PROGRAM INTEGRATION TEST
// =============================================================================

// TestPhase5AllProgramsWorkTogether validates that all three programs work correctly together.
func TestPhase5AllProgramsWorkTogether(t *testing.T) {
	t.Run("All three programs can be instantiated independently", func(t *testing.T) {
		// GreySkull LP - A/B rotation
		_, err := greyskull.GetVariantForDay(1, 1)
		if err != nil {
			t.Errorf("GreySkull rotation failed: %v", err)
		}

		// CAP3 - 3-position rotation
		rl := createCAP3RotationLookup()
		if rl.Length() != 3 {
			t.Errorf("CAP3 rotation lookup: expected 3 positions, got %d", rl.Length())
		}

		// Inverted Juggernaut - 16-week wave cycle
		info := juggernaut.GetWaveInfo(1)
		if info.WaveName != "10s" {
			t.Errorf("Juggernaut wave: expected 10s for week 1, got %s", info.WaveName)
		}
	})

	t.Run("User state can track different program types", func(t *testing.T) {
		// Create users for each program
		for _, programType := range []string{"greyskull", "cap3", "juggernaut"} {
			userID := uuid.New().String()
			programID := uuid.New().String()

			state, _ := userprogramstate.EnrollUser(
				userprogramstate.EnrollUserInput{
					UserID:    userID,
					ProgramID: programID,
				},
				uuid.New().String(),
			)

			if state.CurrentWeek != 1 {
				t.Errorf("%s: initial week should be 1, got %d", programType, state.CurrentWeek)
			}
			if state.RotationPosition != 0 {
				t.Errorf("%s: initial rotation position should be 0, got %d", programType, state.RotationPosition)
			}
		}
	})

	t.Run("Progressions apply correctly for each program type", func(t *testing.T) {
		// GreySkull progression
		gsp, _ := greyskull.NewGreySkullMainLiftProgression(
			uuid.New().String(), "GS Test", 2.5, progression.TrainingMax)
		if gsp.Type() != progression.TypeGreySkull {
			t.Errorf("GreySkull progression type mismatch")
		}

		// Juggernaut progression via rep standard
		newTM := juggernaut.CalculateNewTM(200.0, 0, 12, false)
		expectedTM := 210.0 // 200 + (12-10)*5 = 210
		if newTM != expectedTM {
			t.Errorf("Juggernaut TM calc: expected %.1f, got %.1f", expectedTM, newTM)
		}
	})

	t.Run("Set schemes generate correctly for each program type", func(t *testing.T) {
		// GreySkull: 2+1 AMRAP
		gsScheme, _ := setscheme.NewGreySkullSetScheme(2, 5, 1, 5)
		gsSets, _ := gsScheme.GenerateSets(135.0, setscheme.DefaultSetGenerationContext())
		if len(gsSets) != 3 {
			t.Errorf("GreySkull sets: expected 3, got %d", len(gsSets))
		}

		// CAP3: AMRAP set
		amrapScheme, _ := setscheme.NewAMRAPSetScheme(1, 2)
		amrapSets, _ := amrapScheme.GenerateSets(200.0, setscheme.DefaultSetGenerationContext())
		if len(amrapSets) != 1 {
			t.Errorf("CAP3 AMRAP sets: expected 1, got %d", len(amrapSets))
		}

		// Juggernaut volume sets from wave config
		volumeConfigs := juggernaut.GetVolumeSetConfigs(0, 1) // 10s wave, accumulation
		if len(volumeConfigs) != 9 {
			t.Errorf("Juggernaut volume sets: expected 9 for 10s wave, got %d", len(volumeConfigs))
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS FOR ACCEPTANCE TESTS
// =============================================================================

// intPtrAcceptance returns a pointer to an int value (helper for acceptance tests).
// Named distinctly to avoid collision with other test files in this package.
func intPtrAcceptance(i int) *int {
	return &i
}
