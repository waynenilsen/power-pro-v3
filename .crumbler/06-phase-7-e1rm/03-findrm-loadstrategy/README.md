# FindRM LoadStrategy

Implement a LoadStrategy where the user works up to discover their rep max (no prescribed weight).

## Use Case

Programs like GZCL Jacked & Tan 2.0 have users find their 10RM in week 1, 8RM in week 2, etc. The weight is not prescribed - the user works up until they can only do the target reps.

## Location

Create `internal/domain/loadstrategy/findrm.go`

## Interface

```go
type FindRMLoadStrategy struct {
    TargetReps int  // The rep max to find (e.g., 10 for "find your 10RM")
}

func (s *FindRMLoadStrategy) Type() LoadStrategyType
func (s *FindRMLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error)
func (s *FindRMLoadStrategy) Validate() error
```

## Behavior

- `CalculateLoad` returns 0 (or a sentinel value indicating "user decides")
- The prescription display should show "Find 10RM" not "X lbs × 10"
- After execution, the LoggedSet captures what weight was used
- That LoggedSet can then be used by RelativeTo for back-off sets

## Validation

- TargetReps must be 1-12 (RPE chart range)

## Factory Registration

Register `TypeFindRM` in the strategy factory.

## JSON Format

```json
{
  "type": "FIND_RM",
  "targetReps": 10
}
```

## Tests

- CalculateLoad returns 0 (no prescribed weight)
- Validation rejects invalid rep counts
- Factory unmarshaling works correctly

## Acceptance Criteria

- [ ] FindRM strategy returns 0 for weight
- [ ] TargetReps is validated (1-12)
- [ ] Factory registration complete
- [ ] JSON serialization/deserialization works
- [ ] Comprehensive unit tests pass
