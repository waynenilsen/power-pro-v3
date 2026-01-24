# Update Workout Session Handling

Modify the workout and session handling to support variable set counts. Currently all sets are pre-generated - this task adds the ability to generate sets dynamically during a workout session.

## What Needs to Change

### 1. Workout Set Generation

Update `internal/domain/workout/workout.go`:
- When prescription uses VariableSetScheme, generate only first set as "provisional"
- Add method to request next set based on logged performance

```go
// GenerateNextSet creates the next set for variable schemes
func (w *Workout) GenerateNextSet(prescriptionID uuid.UUID, lastLogged LoggedSetData) (*GeneratedSet, bool, error)

type LoggedSetData struct {
    Reps int
    RPE  *float64
    // Whatever termination conditions need
}
```

### 2. Prescription Awareness

The workout needs to know which prescriptions support variable counts:
- Check if `SetScheme` implements `VariableSetScheme`
- Track which prescriptions are "in progress" vs "complete"

### 3. Session State Tracking

May need new concept or extension:
- Track per-prescription: sets completed, total reps, last RPE
- Know when a variable prescription is "done"
- Handle multiple prescriptions in same workout

### 4. API/Handler Updates

If there's workout logging endpoints, they need to:
- Accept logged set data
- Return next set (if any) or completion status
- Handle the "what's my next set?" query

## Key Files to Modify

- `internal/domain/workout/workout.go` - Add GenerateNextSet method
- `internal/domain/workout/generated_workout.go` - Track provisional sets
- Integration with logging layer (investigate current structure)

## Acceptance Criteria

- [ ] Workouts with variable schemes generate first set only initially
- [ ] GenerateNextSet method implemented
- [ ] Termination conditions evaluated correctly
- [ ] Session can track multiple prescriptions independently
- [ ] API can return "next set" or "complete" status
- [ ] Works alongside fixed-count schemes (no regression)
- [ ] Unit tests for new workout methods
- [ ] Integration tests showing full flow

## Notes

This is the "plumbing" task - it connects the SetScheme changes to the actual workout execution. Requires understanding the current logging/session architecture.
