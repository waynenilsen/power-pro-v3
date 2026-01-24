# Implement FatigueDrop SetScheme

Implement the FatigueDrop set scheme for RTS-style load drop training. Sets continue at progressively lower weights until target RPE is reached.

## What FatigueDrop Does

Example: Squat @ 3 reps, start at RPE 8, drop 5%, stop at RPE 10
1. Set 1: 315 lbs x 3 @ RPE 8.0 (target achieved)
2. Set 2: 299 lbs x 3 @ RPE 8.5 (dropped 5%, continue)
3. Set 3: 284 lbs x 3 @ RPE 9.0 (continue)
4. Set 4: 270 lbs x 3 @ RPE 9.5 (continue)
5. Set 5: 256 lbs x 3 @ RPE 10 (STOP - target RPE reached)

## Implementation

### Create `internal/domain/setscheme/fatigue_drop.go`

```go
type FatigueDrop struct {
    TargetReps      int     `json:"target_reps"`
    StartRPE        float64 `json:"start_rpe"`       // Initial target RPE
    StopRPE         float64 `json:"stop_rpe"`        // Stop when RPE reaches this
    DropPercent     float64 `json:"drop_percent"`    // e.g., 0.05 for 5%
    MaxSets         int     `json:"max_sets"`        // Safety limit (default 10)
}
```

### Required Methods

Implements `VariableSetScheme`:
- `Type() string` - returns "fatigue_drop"
- `Validate() error` - ensure valid parameters
- `GenerateSets(baseWeight, ctx)` - returns first set only (provisional)
- `IsVariableCount() bool` - returns true
- `GetTerminationCondition()` - returns RPEThreshold condition
- `GenerateNextSet(ctx, history, termCtx)` - calculate dropped weight, check termination

### Weight Calculation

Each subsequent set:
```go
nextWeight = previousWeight * (1 - DropPercent)
// Apply rounding from context
```

## Dependencies

- Requires Task 1 (TerminationCondition, VariableSetScheme interface)
- Uses existing rounding utilities

## Key Files

- `internal/domain/setscheme/fatigue_drop.go` - New implementation
- `internal/domain/setscheme/fatigue_drop_test.go` - Unit tests
- `internal/domain/setscheme/factory.go` - Register type

## Acceptance Criteria

- [ ] FatigueDrop struct with JSON tags
- [ ] Implements VariableSetScheme interface
- [ ] First set generated at starting weight
- [ ] Subsequent sets drop by configured percentage
- [ ] Terminates when reported RPE >= StopRPE
- [ ] Respects MaxSets safety limit
- [ ] Validation ensures StopRPE > StartRPE
- [ ] Registered in factory for JSON marshaling
- [ ] Comprehensive unit tests covering:
  - Normal progression to termination
  - MaxSets limit hit
  - Edge cases (0% drop, immediate termination)
