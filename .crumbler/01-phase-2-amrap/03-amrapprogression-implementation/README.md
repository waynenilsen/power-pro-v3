# AMRAPProgression Implementation

Implement AMRAP-based progression rules that adjust training maxes based on actual AMRAP performance.

## What to Implement

### 1. Add New Progression Type

Update `internal/domain/progression/progression.go`:
- Add `TypeAMRAP ProgressionType = "AMRAP_PROGRESSION"` constant
- Add to `ValidProgressionTypes` map

### 2. Add New Trigger Type

Update `internal/domain/progression/progression.go`:
- Add `TriggerAfterSet TriggerType = "AFTER_SET"` constant
- Add to `ValidTriggerTypes` map

### 3. Create `internal/domain/progression/amrap.go`

```go
type AMRAPProgression struct {
    ID               string      `json:"id"`
    Name             string      `json:"name"`
    MaxTypeValue     MaxType     `json:"maxType"`
    TriggerTypeValue TriggerType `json:"triggerType"` // AFTER_SET
    Thresholds       []RepsThreshold `json:"thresholds"`
}

type RepsThreshold struct {
    MinReps   int     `json:"minReps"`   // If reps >= this...
    Increment float64 `json:"increment"` // ...add this much
}
```

**Example thresholds (nSuns style):**
```json
{
  "thresholds": [
    {"minReps": 2, "increment": 5.0},
    {"minReps": 4, "increment": 10.0},
    {"minReps": 6, "increment": 15.0}
  ]
}
```

**Apply logic:**
1. Trigger must be `AFTER_SET`
2. TriggerEvent must include `RepsPerformed` and `IsAMRAP: true`
3. Find highest threshold where `repsPerformed >= minReps`
4. Apply that increment (or 0 if no threshold met)

### 4. Extend TriggerEvent

Update `internal/domain/progression/progression.go` `TriggerEvent`:

```go
type TriggerEvent struct {
    // ... existing fields ...

    // AMRAP-specific fields (for AFTER_SET trigger)
    RepsPerformed *int  `json:"repsPerformed,omitempty"`
    IsAMRAP       bool  `json:"isAMRAP,omitempty"`
    SetWeight     *float64 `json:"setWeight,omitempty"`
}
```

### 5. Create `internal/domain/progression/amrap_test.go`

Unit tests covering:
- Validation (thresholds sorted, positive increments)
- Apply with various rep counts hitting different thresholds
- Apply with non-AMRAP set (should not apply)
- Apply with wrong trigger type (should not apply)
- JSON round-trip

### 6. Register in Factory

In `internal/server/server.go`, add:
```go
progression.RegisterAMRAPProgression(progressionFactory)
```

## Design Notes

- Thresholds must be validated: sorted by minReps ascending, positive increments
- The `AFTER_SET` trigger is fired when a user logs an AMRAP set
- Only the lift that was just performed should progress
- This progression type is commonly used with `AFTER_SET` but could theoretically be used with other triggers

## Integration Point

The API endpoint that logs AMRAP sets (`POST /sessions/{sessionId}/sets`) should trigger progression evaluation when `is_amrap: true`.

## Files to Create/Modify

- `internal/domain/progression/progression.go` (add type constants)
- `internal/domain/progression/amrap.go` (new)
- `internal/domain/progression/amrap_test.go` (new)
- `internal/server/server.go` (register progression)

## Verification

- `go test ./internal/domain/progression/...` passes
- `go build ./...` succeeds
