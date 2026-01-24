package juggernaut

import (
	"testing"
)

func TestGetVolumeSetConfigs_AccumulationSetCounts(t *testing.T) {
	tests := []struct {
		waveIndex    int
		waveName     string
		expectedSets int
	}{
		{0, "10s", 9},
		{1, "8s", 7},
		{2, "5s", 5},
		{3, "3s", 6},
	}

	for _, tt := range tests {
		configs := GetVolumeSetConfigs(tt.waveIndex, 1) // weekInWave=1 is Accumulation
		if len(configs) != tt.expectedSets {
			t.Errorf("%s wave accumulation: expected %d sets, got %d", tt.waveName, tt.expectedSets, len(configs))
		}
	}
}

func TestGetVolumeSetConfigs_AccumulationPercentages(t *testing.T) {
	tests := []struct {
		waveIndex   int
		waveName    string
		expectedPct float64
	}{
		{0, "10s", 60.0},
		{1, "8s", 65.0},
		{2, "5s", 70.0},
		{3, "3s", 75.0},
	}

	for _, tt := range tests {
		configs := GetVolumeSetConfigs(tt.waveIndex, 1)
		for i, cfg := range configs {
			if cfg.Percentage != tt.expectedPct {
				t.Errorf("%s wave accumulation set %d: expected %.1f%%, got %.1f%%",
					tt.waveName, i+1, tt.expectedPct, cfg.Percentage)
			}
		}
	}
}

func TestGetVolumeSetConfigs_IntensificationPattern(t *testing.T) {
	tests := []struct {
		waveIndex     int
		waveName      string
		warmup1Pct    float64
		warmup2Pct    float64
		workPct       float64
		expectedTotal int
	}{
		{0, "10s", 55.0, 62.5, 67.5, 7},  // 2 warmup + 5 work
		{1, "8s", 60.0, 67.5, 72.5, 7},   // 2 warmup + 5 work
		{2, "5s", 65.0, 72.5, 77.5, 7},   // 2 warmup + 5 work
		{3, "3s", 70.0, 77.5, 82.5, 7},   // 2 warmup + 5 work
	}

	for _, tt := range tests {
		configs := GetVolumeSetConfigs(tt.waveIndex, 2) // weekInWave=2 is Intensification
		if len(configs) != tt.expectedTotal {
			t.Errorf("%s wave intensification: expected %d sets, got %d",
				tt.waveName, tt.expectedTotal, len(configs))
			continue
		}

		// Check warmup sets
		if configs[0].Percentage != tt.warmup1Pct {
			t.Errorf("%s wave intensification warmup1: expected %.1f%%, got %.1f%%",
				tt.waveName, tt.warmup1Pct, configs[0].Percentage)
		}
		if configs[1].Percentage != tt.warmup2Pct {
			t.Errorf("%s wave intensification warmup2: expected %.1f%%, got %.1f%%",
				tt.waveName, tt.warmup2Pct, configs[1].Percentage)
		}

		// Check work sets (sets 3-7)
		for i := 2; i < len(configs); i++ {
			if configs[i].Percentage != tt.workPct {
				t.Errorf("%s wave intensification work set %d: expected %.1f%%, got %.1f%%",
					tt.waveName, i+1, tt.workPct, configs[i].Percentage)
			}
		}
	}
}

func TestGetVolumeSetConfigs_RealizationPyramids(t *testing.T) {
	tests := []struct {
		waveIndex       int
		waveName        string
		expectedPcts    []float64
		expectedReps    []int
		expectedSetCount int
	}{
		{0, "10s", []float64{50.0, 60.0, 70.0, 75.0}, []int{5, 3, 1, 10}, 4},
		{1, "8s", []float64{50.0, 60.0, 70.0, 75.0, 80.0}, []int{5, 3, 2, 1, 8}, 5},
		{2, "5s", []float64{50.0, 60.0, 70.0, 75.0, 80.0, 85.0}, []int{5, 3, 2, 1, 1, 5}, 6},
		{3, "3s", []float64{50.0, 60.0, 70.0, 75.0, 80.0, 85.0, 90.0}, []int{5, 3, 2, 1, 1, 1, 3}, 7},
	}

	for _, tt := range tests {
		configs := GetVolumeSetConfigs(tt.waveIndex, 3) // weekInWave=3 is Realization
		if len(configs) != tt.expectedSetCount {
			t.Errorf("%s wave realization: expected %d sets, got %d",
				tt.waveName, tt.expectedSetCount, len(configs))
			continue
		}

		for i, cfg := range configs {
			if cfg.Percentage != tt.expectedPcts[i] {
				t.Errorf("%s wave realization set %d: expected %.1f%%, got %.1f%%",
					tt.waveName, i+1, tt.expectedPcts[i], cfg.Percentage)
			}
			if cfg.TargetReps != tt.expectedReps[i] {
				t.Errorf("%s wave realization set %d: expected %d reps, got %d",
					tt.waveName, i+1, tt.expectedReps[i], cfg.TargetReps)
			}
		}
	}
}

func TestGetVolumeSetConfigs_DeloadReturnsNil(t *testing.T) {
	for waveIndex := 0; waveIndex <= 3; waveIndex++ {
		configs := GetVolumeSetConfigs(waveIndex, 4) // weekInWave=4 is Deload
		if configs != nil {
			t.Errorf("Wave %d deload: expected nil, got %d configs", waveIndex, len(configs))
		}
	}
}

func TestGetVolumeSetConfigs_AMRAPFlagsCorrect(t *testing.T) {
	// Test all non-deload weeks
	for waveIndex := 0; waveIndex <= 3; waveIndex++ {
		for weekInWave := 1; weekInWave <= 3; weekInWave++ {
			configs := GetVolumeSetConfigs(waveIndex, weekInWave)
			if configs == nil {
				t.Errorf("Wave %d, week %d: expected configs, got nil", waveIndex, weekInWave)
				continue
			}

			// Check that only the last set is AMRAP
			for i, cfg := range configs {
				isLast := i == len(configs)-1
				if cfg.IsAMRAP != isLast {
					t.Errorf("Wave %d, week %d, set %d: IsAMRAP=%v, expected %v (last=%v)",
						waveIndex, weekInWave, i+1, cfg.IsAMRAP, isLast, isLast)
				}
			}
		}
	}
}

func TestGetVolumeSetConfigs_SetNumbersSequential(t *testing.T) {
	// Verify set numbers are 1-indexed and sequential
	for waveIndex := 0; waveIndex <= 3; waveIndex++ {
		for weekInWave := 1; weekInWave <= 3; weekInWave++ {
			configs := GetVolumeSetConfigs(waveIndex, weekInWave)
			for i, cfg := range configs {
				expected := i + 1
				if cfg.SetNumber != expected {
					t.Errorf("Wave %d, week %d, set %d: SetNumber=%d, expected %d",
						waveIndex, weekInWave, i+1, cfg.SetNumber, expected)
				}
			}
		}
	}
}

func TestGetVolumeSetConfigs_InvalidInputsReturnNil(t *testing.T) {
	invalidCases := []struct {
		waveIndex  int
		weekInWave int
	}{
		{-1, 1},
		{4, 1},
		{0, 0},
		{0, 5},
		{-1, -1},
		{100, 100},
	}

	for _, tt := range invalidCases {
		configs := GetVolumeSetConfigs(tt.waveIndex, tt.weekInWave)
		if configs != nil {
			t.Errorf("GetVolumeSetConfigs(%d, %d): expected nil for invalid input, got %d configs",
				tt.waveIndex, tt.weekInWave, len(configs))
		}
	}
}

func TestGetVolumeSetConfigs_AccumulationReps(t *testing.T) {
	tests := []struct {
		waveIndex    int
		waveName     string
		expectedReps int
	}{
		{0, "10s", 5},
		{1, "8s", 5},
		{2, "5s", 5},
		{3, "3s", 3}, // 3s wave uses 3 reps, not 5
	}

	for _, tt := range tests {
		configs := GetVolumeSetConfigs(tt.waveIndex, 1)
		for i, cfg := range configs {
			if cfg.TargetReps != tt.expectedReps {
				t.Errorf("%s wave accumulation set %d: expected %d reps, got %d",
					tt.waveName, i+1, tt.expectedReps, cfg.TargetReps)
			}
		}
	}
}
