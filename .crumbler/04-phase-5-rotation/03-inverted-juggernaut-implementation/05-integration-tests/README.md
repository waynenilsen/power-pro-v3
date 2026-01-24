# Inverted Juggernaut Integration Tests

Create comprehensive integration tests for the complete Inverted Juggernaut 5/3/1 system.

## Test File

Create `internal/integration/inverted_juggernaut_test.go`

## Test Scenarios

### 1. Week 1 - 10s Wave Accumulation
Verify:
- Wave info: WaveIndex=0, WeekInWave=1, Phase="Accumulation"
- Volume sets: 9 sets x 5 reps @ 60%
- 5/3/1 sets: 65/75/85/75/65 @ 5/5/5+/5/5+

### 2. Week 4 - 10s Wave Deload
Verify:
- Wave info: WaveIndex=0, WeekInWave=4, Phase="Deload"
- Volume sets: None (GetVolumeSetConfigs returns nil)
- 5/3/1 sets: 40/50/60 @ 5/5/5 (only 3 sets)

### 3. Week 5 - 8s Wave Accumulation (Wave Transition)
Verify:
- Wave info: WaveIndex=1, WeekInWave=1, Phase="Accumulation"
- Volume sets: 7 sets x 5 reps @ 65%
- Wave correctly advances from 0 to 1

### 4. Week 11 - 5s Wave Realization
Verify:
- Wave info: WaveIndex=2, WeekInWave=3, Phase="Realization"
- Volume sets: Ascending pyramid (50/60/70/75/80/85)
- 5/3/1 sets: 75/85/95/85/75

### 5. Week 16 - Cycle Completion
Verify:
- Wave info: WaveIndex=3, WeekInWave=4, Phase="Deload"
- After advancing past week 16, CurrentWeek resets to 1
- CyclesSinceStart increments

### 6. TM Progression After Realization
Verify:
- Week 3 AMRAP with 12 reps @ 75% TM
- Rep standard for 10s wave: 10
- Excess: 12 - 10 = 2
- TM increase: 2 × 5 = 10 (lower body)
- New TM correct

### 7. Full Cycle State Tracking
Advance through all 16 weeks (4 days per week = 64 days):
- Track wave transitions at weeks 5, 9, 13
- Verify CyclesSinceStart = 1 after week 16

### 8. Volume Set Count Per Wave
| Wave | Accum Sets | Week |
|------|------------|------|
| 10s | 9 | 1 |
| 8s | 7 | 5 |
| 5s | 5 | 9 |
| 3s | 6 | 13 |

### 9. 5/3/1 Consistency Across Waves
Verify that weeks 1, 5, 9, 13 all use same 5/3/1 percentages (65/75/85/75/65).

### 10. Weight Calculation Integration
With TM = 200:
- Week 1, Set 3 (85%): 200 × 0.85 = 170
- Week 2, Set 3 (90%): 200 × 0.90 = 180
- Week 3, Set 3 (95%): 200 × 0.95 = 190

## Implementation Pattern

Follow the pattern from `cap3_rotation_test.go`:

```go
func TestInvertedJuggernautWaveInfo(t *testing.T) {
    testCases := []struct {
        name           string
        currentWeek    int
        expectedWave   int
        expectedPhase  string
    }{
        {"Week 1", 1, 0, "Accumulation"},
        {"Week 5", 5, 1, "Accumulation"},
        // ...
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            info := juggernaut.GetWaveInfo(tc.currentWeek)
            // assertions
        })
    }
}
```

## Acceptance Criteria

- [ ] All test scenarios pass
- [ ] Wave transitions verified at weeks 5, 9, 13
- [ ] Volume set counts correct per wave
- [ ] 5/3/1 percentages consistent across waves
- [ ] TM progression calculated correctly
- [ ] Full 16-week cycle completes and resets properly
