# GreySkull Progression Strategy

Implement the GreySkull-specific AMRAP-based progression algorithm.

## Overview

GreySkull LP uses a unique progression system based on AMRAP performance on the final set of each exercise. The number of reps achieved determines the weight change for the next session.

## Progression Rules

### Main Lifts (Bench, OHP, Squat, Deadlift)
- **< 5 reps (failure)**: 10% deload, reset
- **5-9 reps (standard)**: +2.5 lbs upper / +5 lbs lower
- **10+ reps (double)**: +5 lbs upper / +10 lbs lower (double increment)

### Accessory Lifts (Curls, Extensions, Shrugs, Ab Work)
- **< 10 reps**: 10% deload
- **10-14 reps**: maintain weight
- **15+ reps**: add increment (+2.5 lbs typically)

## Implementation

### File Location
`internal/domain/greyskull/progression.go`

### Struct
```go
type GreySkullProgression struct {
    Type            string  // "GREYSKULL"
    WeightIncrement float64 // Standard increment (2.5 or 5)
    MinReps         int     // Minimum reps to avoid deload (5 for main, 10 for accessory)
    DoubleThreshold int     // Reps to trigger double increment (10 for main, 15 for accessory)
    DeloadPercent   float64 // Deload percentage (0.10 = 10%)
}
```

### Interface Compliance
Must implement `Progression` interface:
- `Type() string`
- `Trigger() TriggerType` (returns `TriggerAfterSet`)
- `Apply(context ProgressionContext) (*ProgressionResult, error)`
- `Validate() error`
- JSON marshaling with type discriminator

### Apply Logic
```go
func (p *GreySkullProgression) Apply(ctx ProgressionContext) (*ProgressionResult, error) {
    reps := ctx.RepsPerformed
    weight := ctx.CurrentWeight

    if reps < p.MinReps {
        // Deload: reduce by DeloadPercent
        newWeight := weight * (1 - p.DeloadPercent)
        return &ProgressionResult{NewWeight: newWeight, Reason: "deload"}, nil
    }

    if reps >= p.DoubleThreshold {
        // Double increment
        newWeight := weight + (p.WeightIncrement * 2)
        return &ProgressionResult{NewWeight: newWeight, Reason: "double"}, nil
    }

    // Standard increment
    newWeight := weight + p.WeightIncrement
    return &ProgressionResult{NewWeight: newWeight, Reason: "standard"}, nil
}
```

## Tasks

1. Create `internal/domain/greyskull/progression.go` with `GreySkullProgression` struct
2. Implement `Progression` interface methods
3. Add JSON marshaling with "GREYSKULL" type discriminator
4. Add `Validate()` method ensuring valid thresholds
5. Register with progression factory
6. Create `internal/domain/greyskull/progression_test.go` with comprehensive tests

## Test Cases

- Reps = 3 (< 5): Should deload 10%
- Reps = 5 (exactly min): Should add standard increment
- Reps = 7 (5-9 range): Should add standard increment
- Reps = 10 (exactly double threshold): Should add double increment
- Reps = 15 (well above threshold): Should add double increment
- Accessory variant: MinReps=10, DoubleThreshold=15, different behavior

## Acceptance Criteria

- [ ] GreySkullProgression struct defined
- [ ] Implements Progression interface
- [ ] Correct deload calculation
- [ ] Correct standard increment
- [ ] Correct double increment
- [ ] Factory registration works
- [ ] JSON marshaling/unmarshaling works
- [ ] All unit tests pass
