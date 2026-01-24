# Implement TotalReps SetScheme

Implement a new variable set scheme that tracks total reps across sets.

## What is TotalReps?

TotalReps is for exercises where the user has a target number of total reps to complete, but can distribute them across as many sets as needed. Examples:
- "100 chin-ups" from 5/3/1 Building the Monolith
- "100-200 dips" accessory work

Unlike MRS which targets reps-per-set AND total-reps, TotalReps is purely about cumulative volume - the user can do sets of any rep count.

## Implementation Tasks

1. **Add `TypeTotalReps` constant** in `internal/domain/setscheme/setscheme.go`

2. **Create `totalreps.go`** in `internal/domain/setscheme/`:
   - Struct with fields:
     - `TargetTotalReps` (int) - the rep goal
     - `SuggestedRepsPerSet` (int, optional) - hint for initial set size
     - `MaxSets` (int) - safety limit, default 20 (higher than MRS since reps may be small)
   - Implement `SetScheme` interface
   - Implement `VariableSetScheme` interface
   - Termination: when TotalReps >= TargetTotalReps OR MaxSets reached
   - JSON marshaling with type discriminator

3. **Create `totalreps_test.go`** with comprehensive tests:
   - Unit tests for GenerateSets (first provisional set)
   - Unit tests for GenerateNextSet (subsequent sets)
   - Termination tests (target reached, max sets reached)
   - JSON round-trip tests
   - Validation tests

4. **Register the scheme** in `internal/server/server.go`

5. **Verify SessionService compatibility** - should work without changes since it uses the VariableSetScheme interface

## Key Differences from MRS

| Aspect | MRS | TotalReps |
|--------|-----|-----------|
| Primary use | Strength work (T1/T3) | Accessory volume |
| Rep consistency | Expects similar reps/set | Any reps/set okay |
| Weight | Same weight all sets | Could vary (optional) |
| Failure condition | Reps < MinRepsPerSet | None - just keep going |
| Typical sets | 3-10 | Could be many (20+) |

## JSON Example

```json
{
  "type": "TotalReps",
  "target_total_reps": 100,
  "suggested_reps_per_set": 10,
  "max_sets": 20
}
```

## Acceptance Criteria

- [ ] TotalReps scheme generates provisional sets
- [ ] Tracks cumulative reps across sets
- [ ] Terminates when target reached
- [ ] Has max sets safety limit
- [ ] All tests pass
- [ ] Registered in server startup
