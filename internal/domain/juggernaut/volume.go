package juggernaut

// VolumeSetConfig contains the configuration for a single volume set.
type VolumeSetConfig struct {
	SetNumber  int
	Percentage float64
	TargetReps int
	IsAMRAP    bool
}

// volumeKey is used to look up volume configurations by wave and phase.
type volumeKey struct {
	waveIndex  int
	weekInWave int
}

// volumeConfigs maps (waveIndex, weekInWave) to volume set configurations.
// Deload weeks (weekInWave=4) are not present - they have no volume work.
var volumeConfigs = map[volumeKey][]VolumeSetConfig{
	// Wave 0 (10s wave)
	{0, 1}: accumulation10s(),
	{0, 2}: intensification10s(),
	{0, 3}: realization10s(),

	// Wave 1 (8s wave)
	{1, 1}: accumulation8s(),
	{1, 2}: intensification8s(),
	{1, 3}: realization8s(),

	// Wave 2 (5s wave)
	{2, 1}: accumulation5s(),
	{2, 2}: intensification5s(),
	{2, 3}: realization5s(),

	// Wave 3 (3s wave)
	{3, 1}: accumulation3s(),
	{3, 2}: intensification3s(),
	{3, 3}: realization3s(),
}

// GetVolumeSetConfigs returns the volume set configurations for a given wave and phase.
// Returns nil for deload weeks (no volume work).
func GetVolumeSetConfigs(waveIndex, weekInWave int) []VolumeSetConfig {
	key := volumeKey{waveIndex: waveIndex, weekInWave: weekInWave}
	configs, ok := volumeConfigs[key]
	if !ok {
		return nil
	}
	return configs
}

// Accumulation phase builders - straight sets at fixed percentage, final AMRAP.

func accumulation10s() []VolumeSetConfig {
	return buildAccumulationSets(9, 60.0, 5)
}

func accumulation8s() []VolumeSetConfig {
	return buildAccumulationSets(7, 65.0, 5)
}

func accumulation5s() []VolumeSetConfig {
	return buildAccumulationSets(5, 70.0, 5)
}

func accumulation3s() []VolumeSetConfig {
	return buildAccumulationSets(6, 75.0, 3)
}

func buildAccumulationSets(numSets int, percentage float64, reps int) []VolumeSetConfig {
	configs := make([]VolumeSetConfig, numSets)
	for i := 0; i < numSets; i++ {
		configs[i] = VolumeSetConfig{
			SetNumber:  i + 1,
			Percentage: percentage,
			TargetReps: reps,
			IsAMRAP:    i == numSets-1, // Final set is AMRAP
		}
	}
	return configs
}

// Intensification phase builders - warmup sets then work sets, final AMRAP.

func intensification10s() []VolumeSetConfig {
	// 2 warmup sets (55%, 62.5%), then work sets @ 67.5%, final AMRAP
	return buildIntensificationSets(55.0, 62.5, 67.5, 5, 5)
}

func intensification8s() []VolumeSetConfig {
	// 2 warmup sets (60%, 67.5%), then work sets @ 72.5%, final AMRAP
	return buildIntensificationSets(60.0, 67.5, 72.5, 5, 5)
}

func intensification5s() []VolumeSetConfig {
	// 2 warmup sets (65%, 72.5%), then work sets @ 77.5%, final AMRAP
	return buildIntensificationSets(65.0, 72.5, 77.5, 5, 5)
}

func intensification3s() []VolumeSetConfig {
	// 2 warmup sets (70%, 77.5%), then work sets @ 82.5%, final AMRAP
	return buildIntensificationSets(70.0, 77.5, 82.5, 5, 3)
}

func buildIntensificationSets(warmup1, warmup2, work float64, numWorkSets, reps int) []VolumeSetConfig {
	totalSets := 2 + numWorkSets
	configs := make([]VolumeSetConfig, totalSets)

	// Warmup set 1
	configs[0] = VolumeSetConfig{
		SetNumber:  1,
		Percentage: warmup1,
		TargetReps: reps,
		IsAMRAP:    false,
	}

	// Warmup set 2
	configs[1] = VolumeSetConfig{
		SetNumber:  2,
		Percentage: warmup2,
		TargetReps: reps,
		IsAMRAP:    false,
	}

	// Work sets
	for i := 0; i < numWorkSets; i++ {
		configs[2+i] = VolumeSetConfig{
			SetNumber:  3 + i,
			Percentage: work,
			TargetReps: reps,
			IsAMRAP:    i == numWorkSets-1, // Final set is AMRAP
		}
	}

	return configs
}

// Realization phase builders - ascending pyramid to AMRAP.

func realization10s() []VolumeSetConfig {
	// 50/60/70/75 @ 5/3/1/AMRAP
	return []VolumeSetConfig{
		{SetNumber: 1, Percentage: 50.0, TargetReps: 5, IsAMRAP: false},
		{SetNumber: 2, Percentage: 60.0, TargetReps: 3, IsAMRAP: false},
		{SetNumber: 3, Percentage: 70.0, TargetReps: 1, IsAMRAP: false},
		{SetNumber: 4, Percentage: 75.0, TargetReps: 10, IsAMRAP: true},
	}
}

func realization8s() []VolumeSetConfig {
	// 50/60/70/75/80 @ 5/3/2/1/AMRAP
	return []VolumeSetConfig{
		{SetNumber: 1, Percentage: 50.0, TargetReps: 5, IsAMRAP: false},
		{SetNumber: 2, Percentage: 60.0, TargetReps: 3, IsAMRAP: false},
		{SetNumber: 3, Percentage: 70.0, TargetReps: 2, IsAMRAP: false},
		{SetNumber: 4, Percentage: 75.0, TargetReps: 1, IsAMRAP: false},
		{SetNumber: 5, Percentage: 80.0, TargetReps: 8, IsAMRAP: true},
	}
}

func realization5s() []VolumeSetConfig {
	// 50/60/70/75/80/85 @ 5/3/2/1/1/AMRAP
	return []VolumeSetConfig{
		{SetNumber: 1, Percentage: 50.0, TargetReps: 5, IsAMRAP: false},
		{SetNumber: 2, Percentage: 60.0, TargetReps: 3, IsAMRAP: false},
		{SetNumber: 3, Percentage: 70.0, TargetReps: 2, IsAMRAP: false},
		{SetNumber: 4, Percentage: 75.0, TargetReps: 1, IsAMRAP: false},
		{SetNumber: 5, Percentage: 80.0, TargetReps: 1, IsAMRAP: false},
		{SetNumber: 6, Percentage: 85.0, TargetReps: 5, IsAMRAP: true},
	}
}

func realization3s() []VolumeSetConfig {
	// 50/60/70/75/80/85/90 @ 5/3/2/1/1/1/AMRAP
	return []VolumeSetConfig{
		{SetNumber: 1, Percentage: 50.0, TargetReps: 5, IsAMRAP: false},
		{SetNumber: 2, Percentage: 60.0, TargetReps: 3, IsAMRAP: false},
		{SetNumber: 3, Percentage: 70.0, TargetReps: 2, IsAMRAP: false},
		{SetNumber: 4, Percentage: 75.0, TargetReps: 1, IsAMRAP: false},
		{SetNumber: 5, Percentage: 80.0, TargetReps: 1, IsAMRAP: false},
		{SetNumber: 6, Percentage: 85.0, TargetReps: 1, IsAMRAP: false},
		{SetNumber: 7, Percentage: 90.0, TargetReps: 3, IsAMRAP: true},
	}
}
