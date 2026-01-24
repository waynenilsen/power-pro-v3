// Package juggernaut provides domain logic for the Inverted Juggernaut program.
// This package contains wave state derivation and wave-specific configuration
// for the 16-week Inverted Juggernaut cycle.
package juggernaut

// WaveInfo contains derived wave state information.
// Wave information is derived from the CurrentWeek field (1-16) without
// requiring additional state fields.
type WaveInfo struct {
	WaveIndex     int    // 0-3 (10s, 8s, 5s, 3s waves)
	WaveName      string // "10s", "8s", "5s", "3s"
	WeekInWave    int    // 1-4 within the current wave
	PhaseName     string // "Accumulation", "Intensification", "Realization", "Deload"
	IsDeload      bool   // true for weeks 4, 8, 12, 16
	IsRealization bool   // true for weeks 3, 7, 11, 15
}

// Wave names indexed by wave index (0-3).
var waveNames = []string{"10s", "8s", "5s", "3s"}

// Phase names indexed by week within wave (0-3 for weeks 1-4).
var phaseNames = []string{"Accumulation", "Intensification", "Realization", "Deload"}

// GetWaveInfo derives wave information from the current week (1-16).
// Returns a WaveInfo struct with all derived state.
func GetWaveInfo(currentWeek int) WaveInfo {
	// Calculate wave index (0-3)
	// Weeks 1-4 → 0, Weeks 5-8 → 1, Weeks 9-12 → 2, Weeks 13-16 → 3
	waveIndex := (currentWeek - 1) / 4

	// Calculate week within wave (1-4)
	// Week 1,5,9,13 → 1, Week 2,6,10,14 → 2, etc.
	weekInWave := ((currentWeek - 1) % 4) + 1

	// Phase index is 0-based for internal use
	phaseIndex := weekInWave - 1

	return WaveInfo{
		WaveIndex:     waveIndex,
		WaveName:      waveNames[waveIndex],
		WeekInWave:    weekInWave,
		PhaseName:     phaseNames[phaseIndex],
		IsDeload:      weekInWave == 4,
		IsRealization: weekInWave == 3,
	}
}

// VolumeConfig contains the volume set configuration for accumulation phase.
type VolumeConfig struct {
	Sets       int
	Reps       int
	Percentage float64
}

// WaveVolumeConfig maps wave index to volume set configuration for accumulation phase.
// These are the volume sets performed before the AMRAP set.
var WaveVolumeConfig = map[int]VolumeConfig{
	0: {Sets: 9, Reps: 5, Percentage: 60.0}, // 10s wave: 9x5 @ 60%
	1: {Sets: 7, Reps: 5, Percentage: 65.0}, // 8s wave: 7x5 @ 65%
	2: {Sets: 5, Reps: 5, Percentage: 70.0}, // 5s wave: 5x5 @ 70%
	3: {Sets: 6, Reps: 3, Percentage: 75.0}, // 3s wave: 6x3 @ 75%
}

// AMRAPConfig contains the AMRAP set configuration for a wave.
type AMRAPConfig struct {
	Percentage  float64
	RepStandard int
}

// WaveAMRAPConfig maps wave index to AMRAP target percentages and rep standards.
// RepStandard is the expected/target number of reps at the given percentage.
var WaveAMRAPConfig = map[int]AMRAPConfig{
	0: {Percentage: 75.0, RepStandard: 10}, // 10s wave: AMRAP @ 75%, expect 10 reps
	1: {Percentage: 80.0, RepStandard: 8},  // 8s wave: AMRAP @ 80%, expect 8 reps
	2: {Percentage: 85.0, RepStandard: 5},  // 5s wave: AMRAP @ 85%, expect 5 reps
	3: {Percentage: 90.0, RepStandard: 3},  // 3s wave: AMRAP @ 90%, expect 3 reps
}
