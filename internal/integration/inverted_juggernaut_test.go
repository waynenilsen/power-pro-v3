// Package integration provides integration tests for cross-component behavior.
// This file tests the Inverted Juggernaut 5/3/1 system wave state, volume sets,
// 5/3/1 percentages, and TM progression.
package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/juggernaut"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
)

// TestInvertedJuggernautWaveInfo tests wave state derivation across all 16 weeks.
func TestInvertedJuggernautWaveInfo(t *testing.T) {
	testCases := []struct {
		name          string
		currentWeek   int
		expectedWave  int
		expectedPhase string
		weekInWave    int
	}{
		// Wave 0 (10s) - Weeks 1-4
		{"Week 1 - 10s Accumulation", 1, 0, "Accumulation", 1},
		{"Week 2 - 10s Intensification", 2, 0, "Intensification", 2},
		{"Week 3 - 10s Realization", 3, 0, "Realization", 3},
		{"Week 4 - 10s Deload", 4, 0, "Deload", 4},
		// Wave 1 (8s) - Weeks 5-8
		{"Week 5 - 8s Accumulation", 5, 1, "Accumulation", 1},
		{"Week 6 - 8s Intensification", 6, 1, "Intensification", 2},
		{"Week 7 - 8s Realization", 7, 1, "Realization", 3},
		{"Week 8 - 8s Deload", 8, 1, "Deload", 4},
		// Wave 2 (5s) - Weeks 9-12
		{"Week 9 - 5s Accumulation", 9, 2, "Accumulation", 1},
		{"Week 10 - 5s Intensification", 10, 2, "Intensification", 2},
		{"Week 11 - 5s Realization", 11, 2, "Realization", 3},
		{"Week 12 - 5s Deload", 12, 2, "Deload", 4},
		// Wave 3 (3s) - Weeks 13-16
		{"Week 13 - 3s Accumulation", 13, 3, "Accumulation", 1},
		{"Week 14 - 3s Intensification", 14, 3, "Intensification", 2},
		{"Week 15 - 3s Realization", 15, 3, "Realization", 3},
		{"Week 16 - 3s Deload", 16, 3, "Deload", 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info := juggernaut.GetWaveInfo(tc.currentWeek)

			if info.WaveIndex != tc.expectedWave {
				t.Errorf("WaveIndex: expected %d, got %d", tc.expectedWave, info.WaveIndex)
			}
			if info.PhaseName != tc.expectedPhase {
				t.Errorf("PhaseName: expected %q, got %q", tc.expectedPhase, info.PhaseName)
			}
			if info.WeekInWave != tc.weekInWave {
				t.Errorf("WeekInWave: expected %d, got %d", tc.weekInWave, info.WeekInWave)
			}
		})
	}
}

// TestInvertedJuggernautWeek1_10sWaveAccumulation verifies Week 1 configuration.
func TestInvertedJuggernautWeek1_10sWaveAccumulation(t *testing.T) {
	info := juggernaut.GetWaveInfo(1)

	t.Run("wave info", func(t *testing.T) {
		if info.WaveIndex != 0 {
			t.Errorf("WaveIndex: expected 0, got %d", info.WaveIndex)
		}
		if info.WeekInWave != 1 {
			t.Errorf("WeekInWave: expected 1, got %d", info.WeekInWave)
		}
		if info.PhaseName != "Accumulation" {
			t.Errorf("PhaseName: expected Accumulation, got %q", info.PhaseName)
		}
	})

	t.Run("volume sets - 9 sets x 5 reps @ 60%", func(t *testing.T) {
		configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
		if len(configs) != 9 {
			t.Fatalf("expected 9 volume sets, got %d", len(configs))
		}

		for i, cfg := range configs {
			if cfg.Percentage != 60.0 {
				t.Errorf("set %d: expected 60%%, got %.1f%%", i+1, cfg.Percentage)
			}
			if cfg.TargetReps != 5 {
				t.Errorf("set %d: expected 5 reps, got %d", i+1, cfg.TargetReps)
			}
		}
	})

	t.Run("531 sets - 65/75/85/75/65 @ 5/5/5+/5/5+", func(t *testing.T) {
		lookup := juggernaut.Create531WeeklyLookup("test-id", nil)
		entry := lookup.GetByWeekNumber(1)
		if entry == nil {
			t.Fatal("no entry for week 1")
		}

		expectedPcts := []float64{65.0, 75.0, 85.0, 75.0, 65.0}
		expectedReps := []int{5, 5, -5, 5, -5} // Negative indicates AMRAP

		if len(entry.Percentages) != len(expectedPcts) {
			t.Fatalf("expected %d percentages, got %d", len(expectedPcts), len(entry.Percentages))
		}

		for i, pct := range expectedPcts {
			if entry.Percentages[i] != pct {
				t.Errorf("set %d: expected %.1f%%, got %.1f%%", i+1, pct, entry.Percentages[i])
			}
		}

		for i, rep := range expectedReps {
			if entry.Reps[i] != rep {
				t.Errorf("set %d: expected %d reps, got %d", i+1, rep, entry.Reps[i])
			}
		}
	})
}

// TestInvertedJuggernautWeek4_10sWaveDeload verifies Week 4 deload configuration.
func TestInvertedJuggernautWeek4_10sWaveDeload(t *testing.T) {
	info := juggernaut.GetWaveInfo(4)

	t.Run("wave info", func(t *testing.T) {
		if info.WaveIndex != 0 {
			t.Errorf("WaveIndex: expected 0, got %d", info.WaveIndex)
		}
		if info.WeekInWave != 4 {
			t.Errorf("WeekInWave: expected 4, got %d", info.WeekInWave)
		}
		if info.PhaseName != "Deload" {
			t.Errorf("PhaseName: expected Deload, got %q", info.PhaseName)
		}
		if !info.IsDeload {
			t.Error("IsDeload: expected true")
		}
	})

	t.Run("volume sets - none for deload", func(t *testing.T) {
		configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
		if configs != nil {
			t.Errorf("expected nil volume sets for deload, got %d sets", len(configs))
		}
	})

	t.Run("531 sets - 40/50/60 @ 5/5/5 (only 3 sets)", func(t *testing.T) {
		lookup := juggernaut.Create531WeeklyLookup("test-id", nil)
		entry := lookup.GetByWeekNumber(4)
		if entry == nil {
			t.Fatal("no entry for week 4")
		}

		expectedPcts := []float64{40.0, 50.0, 60.0}
		expectedReps := []int{5, 5, 5}

		if len(entry.Percentages) != 3 {
			t.Fatalf("deload expected 3 sets, got %d", len(entry.Percentages))
		}

		for i, pct := range expectedPcts {
			if entry.Percentages[i] != pct {
				t.Errorf("set %d: expected %.1f%%, got %.1f%%", i+1, pct, entry.Percentages[i])
			}
		}

		for i, rep := range expectedReps {
			if entry.Reps[i] != rep {
				t.Errorf("set %d: expected %d reps, got %d", i+1, rep, entry.Reps[i])
			}
		}
	})
}

// TestInvertedJuggernautWeek5_8sWaveTransition verifies wave transition at week 5.
func TestInvertedJuggernautWeek5_8sWaveTransition(t *testing.T) {
	// Verify wave correctly advances from 0 to 1 between weeks 4 and 5
	info4 := juggernaut.GetWaveInfo(4)
	info5 := juggernaut.GetWaveInfo(5)

	t.Run("wave transition from 10s to 8s", func(t *testing.T) {
		if info4.WaveIndex != 0 {
			t.Errorf("Week 4 WaveIndex: expected 0, got %d", info4.WaveIndex)
		}
		if info5.WaveIndex != 1 {
			t.Errorf("Week 5 WaveIndex: expected 1, got %d", info5.WaveIndex)
		}
		if info5.WaveName != "8s" {
			t.Errorf("Week 5 WaveName: expected 8s, got %q", info5.WaveName)
		}
	})

	t.Run("week 5 is accumulation phase", func(t *testing.T) {
		if info5.WeekInWave != 1 {
			t.Errorf("WeekInWave: expected 1, got %d", info5.WeekInWave)
		}
		if info5.PhaseName != "Accumulation" {
			t.Errorf("PhaseName: expected Accumulation, got %q", info5.PhaseName)
		}
	})

	t.Run("volume sets - 7 sets x 5 reps @ 65%", func(t *testing.T) {
		configs := juggernaut.GetVolumeSetConfigs(info5.WaveIndex, info5.WeekInWave)
		if len(configs) != 7 {
			t.Fatalf("expected 7 volume sets, got %d", len(configs))
		}

		for i, cfg := range configs {
			if cfg.Percentage != 65.0 {
				t.Errorf("set %d: expected 65%%, got %.1f%%", i+1, cfg.Percentage)
			}
			if cfg.TargetReps != 5 {
				t.Errorf("set %d: expected 5 reps, got %d", i+1, cfg.TargetReps)
			}
		}
	})
}

// TestInvertedJuggernautWeek11_5sWaveRealization verifies week 11 realization phase.
func TestInvertedJuggernautWeek11_5sWaveRealization(t *testing.T) {
	info := juggernaut.GetWaveInfo(11)

	t.Run("wave info", func(t *testing.T) {
		if info.WaveIndex != 2 {
			t.Errorf("WaveIndex: expected 2, got %d", info.WaveIndex)
		}
		if info.WeekInWave != 3 {
			t.Errorf("WeekInWave: expected 3, got %d", info.WeekInWave)
		}
		if info.PhaseName != "Realization" {
			t.Errorf("PhaseName: expected Realization, got %q", info.PhaseName)
		}
		if !info.IsRealization {
			t.Error("IsRealization: expected true")
		}
	})

	t.Run("volume sets - ascending pyramid (50/60/70/75/80/85)", func(t *testing.T) {
		configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
		if len(configs) != 6 {
			t.Fatalf("expected 6 volume sets, got %d", len(configs))
		}

		expectedPcts := []float64{50.0, 60.0, 70.0, 75.0, 80.0, 85.0}
		for i, cfg := range configs {
			if cfg.Percentage != expectedPcts[i] {
				t.Errorf("set %d: expected %.1f%%, got %.1f%%", i+1, expectedPcts[i], cfg.Percentage)
			}
		}
	})

	t.Run("531 sets - 75/85/95/85/75", func(t *testing.T) {
		lookup := juggernaut.Create531WeeklyLookup("test-id", nil)
		// Use WeekInWave (3) for the lookup since 5/3/1 uses 4-week pattern
		entry := lookup.GetByWeekNumber(info.WeekInWave)
		if entry == nil {
			t.Fatal("no entry for week 3")
		}

		expectedPcts := []float64{75.0, 85.0, 95.0, 85.0, 75.0}
		for i, pct := range expectedPcts {
			if entry.Percentages[i] != pct {
				t.Errorf("set %d: expected %.1f%%, got %.1f%%", i+1, pct, entry.Percentages[i])
			}
		}
	})
}

// TestInvertedJuggernautWeek16_CycleCompletion verifies cycle completion and reset.
func TestInvertedJuggernautWeek16_CycleCompletion(t *testing.T) {
	info := juggernaut.GetWaveInfo(16)

	t.Run("wave info for week 16", func(t *testing.T) {
		if info.WaveIndex != 3 {
			t.Errorf("WaveIndex: expected 3, got %d", info.WaveIndex)
		}
		if info.WeekInWave != 4 {
			t.Errorf("WeekInWave: expected 4, got %d", info.WeekInWave)
		}
		if info.PhaseName != "Deload" {
			t.Errorf("PhaseName: expected Deload, got %q", info.PhaseName)
		}
		if !info.IsDeload {
			t.Error("IsDeload: expected true")
		}
	})

	t.Run("state resets after week 16", func(t *testing.T) {
		userID := uuid.New().String()
		programID := uuid.New().String()

		state, _ := userprogramstate.EnrollUser(
			userprogramstate.EnrollUserInput{
				UserID:    userID,
				ProgramID: programID,
			},
			uuid.New().String(),
		)

		// Set state to week 16, day 3 (last day before reset)
		state.CurrentWeek = 16
		dayIndex := 3
		state.CurrentDayIndex = &dayIndex

		// Inverted Juggernaut: 16 weeks, 4 training days per week
		ctx := userprogramstate.AdvancementContext{
			DaysInCurrentWeek: 4,
			CycleLengthWeeks:  16,
		}

		// Advance - should complete the cycle
		result, valResult := userprogramstate.AdvanceState(state, ctx)
		if !valResult.Valid {
			t.Fatalf("AdvanceState failed: %v", valResult.Errors)
		}

		if result.NewState.CurrentWeek != 1 {
			t.Errorf("CurrentWeek after cycle: expected 1, got %d", result.NewState.CurrentWeek)
		}
		if result.NewState.CyclesSinceStart != 1 {
			t.Errorf("CyclesSinceStart: expected 1, got %d", result.NewState.CyclesSinceStart)
		}
		if !result.CycleCompleted {
			t.Error("CycleCompleted: expected true")
		}
	})
}

// TestInvertedJuggernautTMProgression verifies TM progression after realization AMRAP.
func TestInvertedJuggernautTMProgression(t *testing.T) {
	t.Run("Week 3 AMRAP with 12 reps @ 75% TM (10s wave)", func(t *testing.T) {
		currentTM := 200.0
		waveIndex := 0    // 10s wave
		amrapReps := 12   // 12 reps achieved
		isUpperBody := false // lower body (squat/deadlift)

		// Rep standard for 10s wave: 10
		// Excess: 12 - 10 = 2
		// TM increase: 2 × 5 (lower body) = 10
		expectedNewTM := 210.0

		newTM := juggernaut.CalculateNewTM(currentTM, waveIndex, amrapReps, isUpperBody)
		if newTM != expectedNewTM {
			t.Errorf("NewTM: expected %.1f, got %.1f", expectedNewTM, newTM)
		}
	})

	t.Run("upper body increment is smaller", func(t *testing.T) {
		currentTM := 150.0
		waveIndex := 0    // 10s wave
		amrapReps := 12   // 12 reps achieved
		isUpperBody := true // upper body (bench/OHP)

		// Excess: 12 - 10 = 2
		// TM increase: 2 × 2.5 (upper body) = 5
		expectedNewTM := 155.0

		newTM := juggernaut.CalculateNewTM(currentTM, waveIndex, amrapReps, isUpperBody)
		if newTM != expectedNewTM {
			t.Errorf("NewTM: expected %.1f, got %.1f", expectedNewTM, newTM)
		}
	})

	t.Run("underperformance decreases TM", func(t *testing.T) {
		currentTM := 200.0
		waveIndex := 0    // 10s wave, rep standard = 10
		amrapReps := 7    // only 7 reps
		isUpperBody := false

		// Excess: 7 - 10 = -3
		// TM decrease: -3 × 5 = -15
		expectedNewTM := 185.0

		newTM := juggernaut.CalculateNewTM(currentTM, waveIndex, amrapReps, isUpperBody)
		if newTM != expectedNewTM {
			t.Errorf("NewTM: expected %.1f, got %.1f", expectedNewTM, newTM)
		}
	})
}

// TestInvertedJuggernautFullCycleStateTracking tests a full 16-week cycle advancement.
func TestInvertedJuggernautFullCycleStateTracking(t *testing.T) {
	userID := uuid.New().String()
	programID := uuid.New().String()

	state, _ := userprogramstate.EnrollUser(
		userprogramstate.EnrollUserInput{
			UserID:    userID,
			ProgramID: programID,
		},
		uuid.New().String(),
	)

	// Inverted Juggernaut: 16 weeks, 4 training days per week
	ctx := userprogramstate.AdvancementContext{
		DaysInCurrentWeek: 4,
		CycleLengthWeeks:  16,
	}

	// Track wave transitions at weeks 5, 9, 13
	waveTransitions := []int{5, 9, 13}
	transitionIdx := 0

	// Advance through all 16 weeks (4 days per week = 64 training days)
	for day := 0; day < 64; day++ {
		result, valResult := userprogramstate.AdvanceState(state, ctx)
		if !valResult.Valid {
			t.Fatalf("Day %d: AdvanceState failed: %v", day+1, valResult.Errors)
		}
		state = result.NewState

		// Check for wave transitions
		if transitionIdx < len(waveTransitions) && state.CurrentWeek == waveTransitions[transitionIdx] {
			info := juggernaut.GetWaveInfo(state.CurrentWeek)
			expectedWave := transitionIdx + 1 // waves 1, 2, 3
			if info.WaveIndex != expectedWave {
				t.Errorf("Week %d transition: expected WaveIndex %d, got %d",
					state.CurrentWeek, expectedWave, info.WaveIndex)
			}
			transitionIdx++
		}
	}

	t.Run("cycle completed after 64 days", func(t *testing.T) {
		if state.CyclesSinceStart != 1 {
			t.Errorf("CyclesSinceStart: expected 1, got %d", state.CyclesSinceStart)
		}
		if state.CurrentWeek != 1 {
			t.Errorf("CurrentWeek: expected 1 (reset), got %d", state.CurrentWeek)
		}
	})
}

// TestInvertedJuggernautVolumeSetCountPerWave verifies set counts for each wave.
func TestInvertedJuggernautVolumeSetCountPerWave(t *testing.T) {
	testCases := []struct {
		wave         int
		waveName     string
		week         int // Sample week for this wave (accumulation)
		expectedSets int
	}{
		{0, "10s", 1, 9},
		{1, "8s", 5, 7},
		{2, "5s", 9, 5},
		{3, "3s", 13, 6},
	}

	for _, tc := range testCases {
		t.Run(tc.waveName+" wave accumulation", func(t *testing.T) {
			info := juggernaut.GetWaveInfo(tc.week)
			if info.WaveIndex != tc.wave {
				t.Fatalf("Week %d: expected wave %d, got %d", tc.week, tc.wave, info.WaveIndex)
			}

			configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
			if len(configs) != tc.expectedSets {
				t.Errorf("expected %d sets, got %d", tc.expectedSets, len(configs))
			}
		})
	}
}

// TestInvertedJuggernaut531ConsistencyAcrossWaves verifies 5/3/1 percentages are consistent.
func TestInvertedJuggernaut531ConsistencyAcrossWaves(t *testing.T) {
	lookup := juggernaut.Create531WeeklyLookup("test-id", nil)

	// Weeks 1, 5, 9, 13 are all "Accumulation" (WeekInWave=1)
	// They should all use the same 5/3/1 percentages (65/75/85/75/65)
	accumulationWeeks := []int{1, 5, 9, 13}
	expectedPcts := []float64{65.0, 75.0, 85.0, 75.0, 65.0}

	for _, week := range accumulationWeeks {
		info := juggernaut.GetWaveInfo(week)
		t.Run("Week "+string(rune('0'+week/10))+string(rune('0'+week%10))+" ("+info.WaveName+" wave)", func(t *testing.T) {
			// 5/3/1 lookup uses WeekInWave (1-4), not absolute week (1-16)
			entry := lookup.GetByWeekNumber(info.WeekInWave)
			if entry == nil {
				t.Fatalf("no entry for WeekInWave %d", info.WeekInWave)
			}

			if len(entry.Percentages) != len(expectedPcts) {
				t.Fatalf("expected %d percentages, got %d", len(expectedPcts), len(entry.Percentages))
			}

			for i, pct := range expectedPcts {
				if entry.Percentages[i] != pct {
					t.Errorf("set %d: expected %.1f%%, got %.1f%%", i+1, pct, entry.Percentages[i])
				}
			}
		})
	}
}

// TestInvertedJuggernautWeightCalculation verifies weight calculations with TM=200.
func TestInvertedJuggernautWeightCalculation(t *testing.T) {
	tm := 200.0
	lookup := juggernaut.Create531WeeklyLookup("test-id", nil)

	testCases := []struct {
		week           int
		weekName       string
		set            int // 1-indexed
		expectedPct    float64
		expectedWeight float64
	}{
		// Week 1 (Accumulation): 85% on set 3
		{1, "Accumulation", 3, 85.0, 170.0},
		// Week 2 (Intensification): 90% on set 3
		{2, "Intensification", 3, 90.0, 180.0},
		// Week 3 (Realization): 95% on set 3
		{3, "Realization", 3, 95.0, 190.0},
	}

	for _, tc := range testCases {
		t.Run(tc.weekName+" Set "+string(rune('0'+tc.set)), func(t *testing.T) {
			entry := lookup.GetByWeekNumber(tc.week)
			if entry == nil {
				t.Fatalf("no entry for week %d", tc.week)
			}

			// Set 3 is index 2
			pct := entry.Percentages[tc.set-1]
			if pct != tc.expectedPct {
				t.Errorf("percentage: expected %.1f, got %.1f", tc.expectedPct, pct)
			}

			weight := tm * (pct / 100.0)
			if weight != tc.expectedWeight {
				t.Errorf("weight: expected %.1f, got %.1f", tc.expectedWeight, weight)
			}
		})
	}
}

// TestInvertedJuggernautWaveTransitionsAtCorrectWeeks verifies wave boundaries.
func TestInvertedJuggernautWaveTransitionsAtCorrectWeeks(t *testing.T) {
	transitions := []struct {
		fromWeek int
		toWeek   int
		fromWave int
		toWave   int
	}{
		{4, 5, 0, 1},   // 10s -> 8s
		{8, 9, 1, 2},   // 8s -> 5s
		{12, 13, 2, 3}, // 5s -> 3s
	}

	for _, tr := range transitions {
		t.Run("Week "+string(rune('0'+tr.fromWeek))+" to "+string(rune('0'+tr.toWeek/10))+string(rune('0'+tr.toWeek%10)), func(t *testing.T) {
			fromInfo := juggernaut.GetWaveInfo(tr.fromWeek)
			toInfo := juggernaut.GetWaveInfo(tr.toWeek)

			if fromInfo.WaveIndex != tr.fromWave {
				t.Errorf("from week %d: expected wave %d, got %d", tr.fromWeek, tr.fromWave, fromInfo.WaveIndex)
			}
			if toInfo.WaveIndex != tr.toWave {
				t.Errorf("to week %d: expected wave %d, got %d", tr.toWeek, tr.toWave, toInfo.WaveIndex)
			}
		})
	}
}

// TestInvertedJuggernautDeloadWeeksHaveNoVolumeSets verifies all deload weeks.
func TestInvertedJuggernautDeloadWeeksHaveNoVolumeSets(t *testing.T) {
	deloadWeeks := []int{4, 8, 12, 16}

	for _, week := range deloadWeeks {
		t.Run("Week "+string(rune('0'+week/10))+string(rune('0'+week%10)), func(t *testing.T) {
			info := juggernaut.GetWaveInfo(week)

			if !info.IsDeload {
				t.Error("expected IsDeload=true")
			}

			configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
			if configs != nil {
				t.Errorf("deload week should have nil volume configs, got %d sets", len(configs))
			}
		})
	}
}

// TestInvertedJuggernautRealizationWeeksHaveAMRAP verifies realization AMRAP sets.
func TestInvertedJuggernautRealizationWeeksHaveAMRAP(t *testing.T) {
	realizationWeeks := []int{3, 7, 11, 15}

	for _, week := range realizationWeeks {
		t.Run("Week "+string(rune('0'+week/10))+string(rune('0'+week%10)), func(t *testing.T) {
			info := juggernaut.GetWaveInfo(week)

			if !info.IsRealization {
				t.Error("expected IsRealization=true")
			}

			configs := juggernaut.GetVolumeSetConfigs(info.WaveIndex, info.WeekInWave)
			if configs == nil {
				t.Fatal("expected non-nil volume configs for realization")
			}

			// Last set should be AMRAP
			lastSet := configs[len(configs)-1]
			if !lastSet.IsAMRAP {
				t.Error("expected last set to be AMRAP")
			}
		})
	}
}

// TestInvertedJuggernautRepStandardsByWave verifies rep standards for TM progression.
func TestInvertedJuggernautRepStandardsByWave(t *testing.T) {
	testCases := []struct {
		waveIndex   int
		waveName    string
		repStandard int
	}{
		{0, "10s", 10},
		{1, "8s", 8},
		{2, "5s", 5},
		{3, "3s", 3},
	}

	for _, tc := range testCases {
		t.Run(tc.waveName+" wave", func(t *testing.T) {
			if juggernaut.RepStandards[tc.waveIndex] != tc.repStandard {
				t.Errorf("expected rep standard %d, got %d",
					tc.repStandard, juggernaut.RepStandards[tc.waveIndex])
			}
		})
	}
}

// TestInvertedJuggernautMultipleCycleProgression tests TM progression across cycles.
func TestInvertedJuggernautMultipleCycleProgression(t *testing.T) {
	// Start with TM = 200 for squat (lower body)
	initialTM := 200.0

	// Simulate performing each wave's realization AMRAP exactly at standard
	waveResults := []struct {
		waveIndex int
		reps      int
	}{
		{0, 10}, // 10s wave: exactly 10 reps
		{1, 8},  // 8s wave: exactly 8 reps
		{2, 5},  // 5s wave: exactly 5 reps
		{3, 3},  // 3s wave: exactly 3 reps
	}

	tm := initialTM
	for _, result := range waveResults {
		newTM := juggernaut.CalculateNewTM(tm, result.waveIndex, result.reps, false)
		// Meeting exactly the rep standard should result in no change
		if newTM != tm {
			t.Errorf("Wave %d: TM should stay at %.1f when meeting rep standard, got %.1f",
				result.waveIndex, tm, newTM)
		}
	}

	// Now test with exceeding rep standards
	tm = initialTM
	for _, result := range waveResults {
		// Exceed by 2 reps each wave
		newTM := juggernaut.CalculateNewTM(tm, result.waveIndex, result.reps+2, false)
		// Lower body: 2 excess × 5 = 10
		expectedTM := tm + 10.0
		if newTM != expectedTM {
			t.Errorf("Wave %d: expected TM %.1f after +2 reps, got %.1f",
				result.waveIndex, expectedTM, newTM)
		}
		tm = newTM
	}
}
