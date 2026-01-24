# Wave State Helpers

Create helper functions to derive wave state from the existing CurrentWeek field in UserProgramState.

## Context

The Inverted Juggernaut runs 16 weeks with 4 waves, each lasting 4 weeks:
- Wave 0 (10s): Weeks 1-4
- Wave 1 (8s): Weeks 5-8
- Wave 2 (5s): Weeks 9-12
- Wave 3 (3s): Weeks 13-16

Instead of adding new fields to UserProgramState, we can derive wave information from the existing `CurrentWeek` field (1-indexed, 1-16).

## Implementation

Create a new file: `internal/domain/juggernaut/wave.go`

```go
package juggernaut

// WaveInfo contains derived wave state information.
type WaveInfo struct {
    WaveIndex     int    // 0-3 (10s, 8s, 5s, 3s)
    WaveName      string // "10s", "8s", "5s", "3s"
    WeekInWave    int    // 1-4
    PhaseName     string // "Accumulation", "Intensification", "Realization", "Deload"
    IsDeload      bool
    IsRealization bool
}

// GetWaveInfo derives wave information from the current week (1-16).
func GetWaveInfo(currentWeek int) WaveInfo
```

### Wave Index Calculation
```
WaveIndex = (CurrentWeek - 1) / 4
```
- Weeks 1-4 → (0-3)/4 = 0
- Weeks 5-8 → (4-7)/4 = 1
- Weeks 9-12 → (8-11)/4 = 2
- Weeks 13-16 → (12-15)/4 = 3

### Week Within Wave Calculation
```
WeekInWave = ((CurrentWeek - 1) % 4) + 1
```
- Week 1, 5, 9, 13 → 1 (Accumulation)
- Week 2, 6, 10, 14 → 2 (Intensification)
- Week 3, 7, 11, 15 → 3 (Realization)
- Week 4, 8, 12, 16 → 4 (Deload)

## Wave-Specific Constants

Also define constants for wave-specific values:

```go
// Volume set configurations per wave (for Accumulation phase)
var WaveVolumeConfig = map[int]struct {
    Sets       int
    Reps       int
    Percentage float64
}{
    0: {9, 5, 60.0},  // 10s wave: 9x5 @ 60%
    1: {7, 5, 65.0},  // 8s wave: 7x5 @ 65%
    2: {5, 5, 70.0},  // 5s wave: 5x5 @ 70%
    3: {6, 3, 75.0},  // 3s wave: 6x3 @ 75%
}

// AMRAP target percentages and rep standards per wave
var WaveAMRAPConfig = map[int]struct {
    Percentage  float64
    RepStandard int
}{
    0: {75.0, 10},  // 10s wave: AMRAP @ 75%, expect 10 reps
    1: {80.0, 8},   // 8s wave: AMRAP @ 80%, expect 8 reps
    2: {85.0, 5},   // 5s wave: AMRAP @ 85%, expect 5 reps
    3: {90.0, 3},   // 3s wave: AMRAP @ 90%, expect 3 reps
}
```

## Tests

Create `internal/domain/juggernaut/wave_test.go` with tests:
1. GetWaveInfo returns correct WaveIndex for weeks 1-16
2. GetWaveInfo returns correct WeekInWave (1-4) for each week
3. GetWaveInfo returns correct phase names
4. IsDeload is true only for weeks 4, 8, 12, 16
5. IsRealization is true only for weeks 3, 7, 11, 15

## Acceptance Criteria

- [ ] WaveInfo struct defined with all needed fields
- [ ] GetWaveInfo correctly calculates wave index from week number
- [ ] GetWaveInfo correctly calculates week within wave
- [ ] Phase names correctly derived
- [ ] Unit tests pass for all 16 weeks
