# Implement MRS (Max Rep Sets) SetScheme

Implement the MRS set scheme for GZCL-style max rep set training. Perform sets at fixed weight until rep target is met or technical failure occurs.

## What MRS Does

Example: Bench Press MRS x 3, target 25 total reps
1. Set 1: 225 lbs x 10 (total: 10, continue)
2. Set 2: 225 lbs x 8 (total: 18, continue)
3. Set 3: 225 lbs x 5 (total: 23, continue)
4. Set 4: 225 lbs x 4 (total: 27, STOP - exceeded 25)

Alternative termination - failure:
1. Set 1: 225 lbs x 10
2. Set 2: 225 lbs x 6
3. Set 3: 225 lbs x 3 (STOP - technical failure, couldn't hit minimum)

## Implementation

### Create `internal/domain/setscheme/mrs.go`

```go
type MRS struct {
    TargetTotalReps int     `json:"target_total_reps"` // Stop when total reps >= this
    MinRepsPerSet   int     `json:"min_reps_per_set"`  // Failure if can't hit this
    MaxSets         int     `json:"max_sets"`          // Safety limit (default 10)
    NumberOfMRS     int     `json:"number_of_mrs"`     // How many MRS blocks (GZCL uses 3)
}
```

### Required Methods

Implements `VariableSetScheme`:
- `Type() string` - returns "mrs"
- `Validate() error` - ensure valid parameters
- `GenerateSets(baseWeight, ctx)` - returns first set only (AMRAP style)
- `IsVariableCount() bool` - returns true
- `GetTerminationCondition()` - composite: TotalReps OR RepFailure
- `GenerateNextSet(ctx, history, termCtx)` - same weight, check termination

### Termination Logic

Stop when ANY of:
1. `TotalReps >= TargetTotalReps`
2. `LastReps < MinRepsPerSet` (failure)
3. `TotalSets >= MaxSets` (safety)

## Dependencies

- Requires Task 1 (TerminationCondition, VariableSetScheme interface)
- Similar pattern to FatigueDrop but fixed weight

## Key Files

- `internal/domain/setscheme/mrs.go` - New implementation
- `internal/domain/setscheme/mrs_test.go` - Unit tests
- `internal/domain/setscheme/factory.go` - Register type

## Acceptance Criteria

- [ ] MRS struct with JSON tags
- [ ] Implements VariableSetScheme interface
- [ ] All sets at same weight (no drops)
- [ ] Terminates when total reps >= target
- [ ] Terminates when reps < min threshold (failure)
- [ ] Respects MaxSets safety limit
- [ ] NumberOfMRS field for GZCL patterns (T1: 3 MRS, T3: 4 MRS)
- [ ] Validation ensures MinRepsPerSet > 0, TargetTotalReps > MinRepsPerSet
- [ ] Registered in factory for JSON marshaling
- [ ] Comprehensive unit tests covering:
  - Normal progression to target reps
  - Failure termination
  - MaxSets limit hit
  - Multiple MRS blocks
