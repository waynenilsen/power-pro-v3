package juggernaut

import (
	"testing"
)

func TestGetWaveInfo_WaveIndex(t *testing.T) {
	tests := []struct {
		week          int
		expectedIndex int
		expectedName  string
	}{
		// Wave 0 (10s): Weeks 1-4
		{1, 0, "10s"},
		{2, 0, "10s"},
		{3, 0, "10s"},
		{4, 0, "10s"},
		// Wave 1 (8s): Weeks 5-8
		{5, 1, "8s"},
		{6, 1, "8s"},
		{7, 1, "8s"},
		{8, 1, "8s"},
		// Wave 2 (5s): Weeks 9-12
		{9, 2, "5s"},
		{10, 2, "5s"},
		{11, 2, "5s"},
		{12, 2, "5s"},
		// Wave 3 (3s): Weeks 13-16
		{13, 3, "3s"},
		{14, 3, "3s"},
		{15, 3, "3s"},
		{16, 3, "3s"},
	}

	for _, tt := range tests {
		info := GetWaveInfo(tt.week)
		if info.WaveIndex != tt.expectedIndex {
			t.Errorf("Week %d: expected WaveIndex %d, got %d", tt.week, tt.expectedIndex, info.WaveIndex)
		}
		if info.WaveName != tt.expectedName {
			t.Errorf("Week %d: expected WaveName %q, got %q", tt.week, tt.expectedName, info.WaveName)
		}
	}
}

func TestGetWaveInfo_WeekInWave(t *testing.T) {
	tests := []struct {
		week             int
		expectedWeekIn   int
	}{
		// Week 1 of each wave
		{1, 1},
		{5, 1},
		{9, 1},
		{13, 1},
		// Week 2 of each wave
		{2, 2},
		{6, 2},
		{10, 2},
		{14, 2},
		// Week 3 of each wave
		{3, 3},
		{7, 3},
		{11, 3},
		{15, 3},
		// Week 4 of each wave
		{4, 4},
		{8, 4},
		{12, 4},
		{16, 4},
	}

	for _, tt := range tests {
		info := GetWaveInfo(tt.week)
		if info.WeekInWave != tt.expectedWeekIn {
			t.Errorf("Week %d: expected WeekInWave %d, got %d", tt.week, tt.expectedWeekIn, info.WeekInWave)
		}
	}
}

func TestGetWaveInfo_PhaseNames(t *testing.T) {
	tests := []struct {
		week          int
		expectedPhase string
	}{
		// Accumulation (week 1 of each wave)
		{1, "Accumulation"},
		{5, "Accumulation"},
		{9, "Accumulation"},
		{13, "Accumulation"},
		// Intensification (week 2 of each wave)
		{2, "Intensification"},
		{6, "Intensification"},
		{10, "Intensification"},
		{14, "Intensification"},
		// Realization (week 3 of each wave)
		{3, "Realization"},
		{7, "Realization"},
		{11, "Realization"},
		{15, "Realization"},
		// Deload (week 4 of each wave)
		{4, "Deload"},
		{8, "Deload"},
		{12, "Deload"},
		{16, "Deload"},
	}

	for _, tt := range tests {
		info := GetWaveInfo(tt.week)
		if info.PhaseName != tt.expectedPhase {
			t.Errorf("Week %d: expected PhaseName %q, got %q", tt.week, tt.expectedPhase, info.PhaseName)
		}
	}
}

func TestGetWaveInfo_IsDeload(t *testing.T) {
	deloadWeeks := map[int]bool{4: true, 8: true, 12: true, 16: true}

	for week := 1; week <= 16; week++ {
		info := GetWaveInfo(week)
		expected := deloadWeeks[week]
		if info.IsDeload != expected {
			t.Errorf("Week %d: expected IsDeload %v, got %v", week, expected, info.IsDeload)
		}
	}
}

func TestGetWaveInfo_IsRealization(t *testing.T) {
	realizationWeeks := map[int]bool{3: true, 7: true, 11: true, 15: true}

	for week := 1; week <= 16; week++ {
		info := GetWaveInfo(week)
		expected := realizationWeeks[week]
		if info.IsRealization != expected {
			t.Errorf("Week %d: expected IsRealization %v, got %v", week, expected, info.IsRealization)
		}
	}
}

func TestWaveVolumeConfig(t *testing.T) {
	tests := []struct {
		waveIndex      int
		expectedSets   int
		expectedReps   int
		expectedPct    float64
	}{
		{0, 9, 5, 60.0},
		{1, 7, 5, 65.0},
		{2, 5, 5, 70.0},
		{3, 6, 3, 75.0},
	}

	for _, tt := range tests {
		cfg := WaveVolumeConfig[tt.waveIndex]
		if cfg.Sets != tt.expectedSets {
			t.Errorf("Wave %d: expected Sets %d, got %d", tt.waveIndex, tt.expectedSets, cfg.Sets)
		}
		if cfg.Reps != tt.expectedReps {
			t.Errorf("Wave %d: expected Reps %d, got %d", tt.waveIndex, tt.expectedReps, cfg.Reps)
		}
		if cfg.Percentage != tt.expectedPct {
			t.Errorf("Wave %d: expected Percentage %.1f, got %.1f", tt.waveIndex, tt.expectedPct, cfg.Percentage)
		}
	}
}

func TestWaveAMRAPConfig(t *testing.T) {
	tests := []struct {
		waveIndex       int
		expectedPct     float64
		expectedStd     int
	}{
		{0, 75.0, 10},
		{1, 80.0, 8},
		{2, 85.0, 5},
		{3, 90.0, 3},
	}

	for _, tt := range tests {
		cfg := WaveAMRAPConfig[tt.waveIndex]
		if cfg.Percentage != tt.expectedPct {
			t.Errorf("Wave %d: expected Percentage %.1f, got %.1f", tt.waveIndex, tt.expectedPct, cfg.Percentage)
		}
		if cfg.RepStandard != tt.expectedStd {
			t.Errorf("Wave %d: expected RepStandard %d, got %d", tt.waveIndex, tt.expectedStd, cfg.RepStandard)
		}
	}
}

func TestGetWaveInfo_AllWeeksValid(t *testing.T) {
	// Verify all 16 weeks produce valid WaveInfo
	for week := 1; week <= 16; week++ {
		info := GetWaveInfo(week)

		// WaveIndex should be 0-3
		if info.WaveIndex < 0 || info.WaveIndex > 3 {
			t.Errorf("Week %d: WaveIndex %d out of range [0,3]", week, info.WaveIndex)
		}

		// WaveName should be non-empty
		if info.WaveName == "" {
			t.Errorf("Week %d: WaveName is empty", week)
		}

		// WeekInWave should be 1-4
		if info.WeekInWave < 1 || info.WeekInWave > 4 {
			t.Errorf("Week %d: WeekInWave %d out of range [1,4]", week, info.WeekInWave)
		}

		// PhaseName should be non-empty
		if info.PhaseName == "" {
			t.Errorf("Week %d: PhaseName is empty", week)
		}

		// Verify configs exist for this wave
		if _, ok := WaveVolumeConfig[info.WaveIndex]; !ok {
			t.Errorf("Week %d: missing WaveVolumeConfig for wave %d", week, info.WaveIndex)
		}
		if _, ok := WaveAMRAPConfig[info.WaveIndex]; !ok {
			t.Errorf("Week %d: missing WaveAMRAPConfig for wave %d", week, info.WaveIndex)
		}
	}
}
