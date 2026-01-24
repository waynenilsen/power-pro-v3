# Volume Set Generation

Implement the Juggernaut volume work that precedes the 5/3/1 main sets.

## Context

Before the 5/3/1 work each session, Inverted Juggernaut has wave-specific volume sets. These vary by:
1. **Wave** (10s/8s/5s/3s)
2. **Phase within wave** (Accumulation/Intensification/Realization/Deload)

## Volume Set Specifications

### Accumulation Phase (Week 1 of each wave)
| Wave | Sets | Reps | Percentage |
|------|------|------|------------|
| 10s | 9 | 5 | 60% TM |
| 8s | 7 | 5 | 65% TM |
| 5s | 5 | 5 | 70% TM |
| 3s | 6 | 3 | 75% TM |

Final set is AMRAP for the target reps.

### Intensification Phase (Week 2 of each wave)
Complex ascending pattern - simplified for implementation:

| Wave | Pattern |
|------|---------|
| 10s | 2 warmup sets (55%, 62.5%), then sets @ 67.5%, final AMRAP |
| 8s | 2 warmup sets (60%, 67.5%), then sets @ 72.5%, final AMRAP |
| 5s | 2 warmup sets (65%, 72.5%), then sets @ 77.5%, final AMRAP |
| 3s | 2 warmup sets (70%, 77.5%), then sets @ 82.5%, final AMRAP |

### Realization Phase (Week 3 of each wave)
Ascending pyramid to AMRAP:

| Wave | Percentages | Reps |
|------|-------------|------|
| 10s | 50/60/70/75 | 5/3/1/AMRAP |
| 8s | 50/60/70/75/80 | 5/3/2/1/AMRAP |
| 5s | 50/60/70/75/80/85 | 5/3/2/1/1/AMRAP |
| 3s | 50/60/70/75/80/85/90 | 5/3/2/1/1/1/AMRAP |

### Deload Phase (Week 4)
No volume work - just the 5/3/1 deload sets.

## Implementation

Create `internal/domain/juggernaut/volume.go`:

```go
package juggernaut

// VolumeSetConfig contains the configuration for a single volume set.
type VolumeSetConfig struct {
    SetNumber   int
    Percentage  float64
    TargetReps  int
    IsAMRAP     bool
}

// GetVolumeSetConfigs returns the volume set configurations for a given wave and phase.
// Returns nil for deload weeks (no volume work).
func GetVolumeSetConfigs(waveIndex, weekInWave int) []VolumeSetConfig
```

This returns the volume set configurations. The actual weight calculation happens at resolution time using the training max.

## Implementation Approach

Use a lookup table keyed by `(waveIndex, weekInWave)` that returns the slice of VolumeSetConfig.

For Accumulation (weekInWave=1):
```go
// 10s wave accumulation
[]VolumeSetConfig{
    {1, 60.0, 5, false},
    {2, 60.0, 5, false},
    // ... 8 more sets
    {9, 60.0, 5, true}, // Final AMRAP
}
```

## Tests

Create `internal/domain/juggernaut/volume_test.go`:

1. Accumulation phase returns correct number of sets per wave (9/7/5/6)
2. Accumulation phase returns correct percentages per wave (60/65/70/75)
3. Intensification phase returns ascending warmup + work sets
4. Realization phase returns correct pyramid
5. Deload phase returns nil (no volume work)
6. Final set of each non-deload week is marked AMRAP

## Acceptance Criteria

- [ ] GetVolumeSetConfigs returns correct set count per wave/phase
- [ ] Accumulation percentages correct (60/65/70/75)
- [ ] Intensification patterns correct
- [ ] Realization pyramids correct
- [ ] Deload returns nil
- [ ] AMRAP flags set correctly
- [ ] Unit tests pass
