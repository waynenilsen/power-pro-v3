# Wave Advancement Logic

Implement TM progression after each wave's Realization phase.

## Context

After Week 3 of each wave (Realization), the Training Max is recalculated based on AMRAP performance. The existing AMRAP_PROGRESSION type can be used for this.

## Progression Rules

### AMRAP Performance Targets
| Wave | AMRAP % TM | Rep Standard |
|------|-----------|--------------|
| 10s | 75% | 10 reps |
| 8s | 80% | 8 reps |
| 5s | 85% | 5 reps |
| 3s | 90% | 3 reps |

### TM Calculation Formula
```
New TM = Previous TM + ((AMRAP Reps - Rep Standard) × Weight Increment)
```

Weight Increments:
- Bench/OHP: 2.5 units per rep over standard
- Squat/Deadlift: 5 units per rep over standard

### Example
- Current Squat TM: 200
- Wave: 10s (Rep Standard: 10)
- AMRAP performed: 13 reps @ 75%
- Excess reps: 13 - 10 = 3
- Increase: 3 × 5 = 15
- New TM: 200 + 15 = 215

### Underperformance
If AMRAP reps < rep standard, TM can decrease (negative adjustment) or stay same. The formula handles both cases naturally.

## Implementation

Create `internal/domain/juggernaut/progression.go`:

```go
package juggernaut

// CalculateNewTM calculates the new Training Max based on AMRAP performance.
//
// Parameters:
//   - currentTM: Current training max
//   - waveIndex: 0-3 (determines rep standard)
//   - amrapReps: Actual reps performed on realization AMRAP
//   - isUpperBody: true for bench/OHP (2.5 increment), false for squat/deadlift (5 increment)
//
// Returns the new training max.
func CalculateNewTM(currentTM float64, waveIndex int, amrapReps int, isUpperBody bool) float64
```

### Integration with Existing Progression System

The existing `progression` package has `AMRAP_PROGRESSION` type. We may either:
1. Use it directly if it supports custom rep standards and increments
2. Create a Juggernaut-specific wrapper that uses the formula above

Check `internal/domain/progression/` for the existing AMRAP progression implementation.

## Cycle-to-Cycle Progression

After completing a full 16-week cycle, base TM increases:
- Squat/Deadlift: +10 units
- Bench/OHP: +5 units

This happens when `CyclesSinceStart` increments (detected by `AdvanceState` returning `CycleCompleted = true`).

## Tests

Create `internal/domain/juggernaut/progression_test.go`:

1. Exact rep standard returns same TM
2. Exceeding rep standard increases TM correctly
3. Upper body uses 2.5 increment
4. Lower body uses 5 increment
5. Underperformance decreases TM
6. Each wave uses correct rep standard (10/8/5/3)

## Acceptance Criteria

- [ ] CalculateNewTM returns correct TM for over-performance
- [ ] CalculateNewTM returns correct TM for under-performance
- [ ] CalculateNewTM returns same TM for exact rep standard
- [ ] Increment values correct (2.5 upper, 5 lower)
- [ ] Rep standards correct per wave (10/8/5/3)
- [ ] Unit tests pass
