# Composite Set Scheme (Fixed + AMRAP)

Implement a composite set scheme that combines fixed sets with a final AMRAP set for GreySkull LP.

## Overview

GreySkull LP main lifts use a unique set scheme:
- 2 sets of 5 reps (fixed)
- 1 set of 5+ reps (AMRAP)

This is written as: 2x5 + 1x5+

## Current Set Schemes Available

- `FixedSetScheme`: N sets of same weight/reps
- `AMRAPSetScheme`: N AMRAP sets with minimum reps
- `RampSetScheme`: Progressive warmup sets
- `RepRangeSetScheme`: Sets with min/max rep range

## Implementation Options

### Option 1: Composite Set Scheme (Recommended)
Create a new `GreySkullSetScheme` that wraps Fixed + AMRAP:

```go
type GreySkullSetScheme struct {
    Type       string  // "GREYSKULL"
    FixedSets  int     // Number of fixed sets (2)
    FixedReps  int     // Reps per fixed set (5)
    AMRAPSets  int     // Number of AMRAP sets (1)
    MinAMRAPReps int   // Minimum reps for AMRAP (5)
}

func (s *GreySkullSetScheme) GenerateSets(weight float64) []Set {
    sets := make([]Set, 0, s.FixedSets + s.AMRAPSets)

    // Fixed sets
    for i := 0; i < s.FixedSets; i++ {
        sets = append(sets, Set{
            Weight:    weight,
            Reps:      s.FixedReps,
            IsWorkSet: true,
            SetType:   "WORK",
        })
    }

    // AMRAP set
    for i := 0; i < s.AMRAPSets; i++ {
        sets = append(sets, Set{
            Weight:    weight,
            Reps:      s.MinAMRAPReps, // Minimum target
            IsWorkSet: true,
            SetType:   "AMRAP",
        })
    }

    return sets
}
```

### Option 2: Use Existing Schemes in Sequence
Configure prescription to use multiple set schemes in sequence. May require prescription changes.

## File Location
`internal/domain/setscheme/greyskull.go`

## Interface Compliance
Must implement `SetScheme` interface:
- `Type() string`
- `GenerateSets(weight float64) []Set`
- `Validate() error`
- JSON marshaling with type discriminator

## Tasks

1. Create `internal/domain/setscheme/greyskull.go`
2. Define `GreySkullSetScheme` struct
3. Implement `SetScheme` interface
4. Add JSON marshaling with "GREYSKULL" type discriminator
5. Register with set scheme factory
6. Create `internal/domain/setscheme/greyskull_test.go` with tests

## Test Cases

- Standard config (2x5 + 1x5+): Should generate 3 sets
- First 2 sets: Fixed, weight=100, reps=5, type="WORK"
- Last set: AMRAP, weight=100, reps=5, type="AMRAP"
- Accessory config (2x12 + 1x12+): Should work with different reps
- Validate: FixedSets >= 0, FixedReps >= 1, AMRAPSets >= 1, MinAMRAPReps >= 1

## Accessory Set Scheme

Accessories may use standard `RepRangeSetScheme` with:
- Sets: 3
- MinReps: 10
- MaxReps: 15

Or we could add an accessory mode to `GreySkullSetScheme`.

## Acceptance Criteria

- [ ] GreySkullSetScheme struct defined
- [ ] Implements SetScheme interface
- [ ] Generates correct number of sets
- [ ] Fixed sets have correct type and reps
- [ ] AMRAP set has correct type
- [ ] Factory registration works
- [ ] JSON marshaling/unmarshaling works
- [ ] All unit tests pass
