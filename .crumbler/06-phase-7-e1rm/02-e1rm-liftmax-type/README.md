# E1RM LiftMax Type

Add E1RM as a new LiftMax type that stores estimated 1RM values calculated from performed sets.

## Changes to LiftMax Domain

In `internal/domain/liftmax/liftmax.go`:

1. Add new MaxType constant:
```go
const (
    OneRM       MaxType = "ONE_RM"
    TrainingMax MaxType = "TRAINING_MAX"
    E1RM        MaxType = "E1RM"  // NEW
)
```

2. Update `ValidateType()` to accept E1RM

3. Add to MaxCalculator if needed for conversions

## Changes to LoadStrategy

In `internal/domain/loadstrategy/loadstrategy.go`:

1. Add new ReferenceType constant:
```go
const (
    ReferenceOneRM       ReferenceType = "ONE_RM"
    ReferenceTrainingMax ReferenceType = "TRAINING_MAX"
    ReferenceE1RM        ReferenceType = "E1RM"  // NEW
)
```

2. Ensure PercentOfLoadStrategy can reference E1RM type

## Database

No migration needed - E1RM is just a new value for the existing `max_type` column.

## Tests

- Validate E1RM is accepted as MaxType
- Validate E1RM can be used in PercentOf strategy
- Integration: Create E1RM LiftMax, use it in prescription calculation

## Acceptance Criteria

- [ ] E1RM is a valid MaxType constant
- [ ] E1RM can be stored and retrieved as LiftMax
- [ ] PercentOf strategy can reference E1RM
- [ ] Unit tests for E1RM validation
- [ ] Integration test for E1RM → PercentOf flow
