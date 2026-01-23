# DeloadOnFailure Progression

Implement the `DeloadOnFailure` progression type that reduces weight after consecutive failures.

## Context

Many programs (GZCLP T2/T3, Texas Method reset, etc.) reduce weight after failing multiple times. This progression type allows configurable deload behavior.

## Implementation Requirements

### 1. DeloadOnFailure Progression Type

Create `internal/domain/progression/deload_on_failure.go`:

```go
type DeloadOnFailure struct {
    FailureThreshold int     // Number of consecutive failures before deload (e.g., 3)
    DeloadPercent    float64 // Percentage to reduce (e.g., 0.10 for 10% deload)
    DeloadAmount     float64 // Or fixed amount to reduce (e.g., 10 lbs)
    DeloadType       string  // "percent" or "fixed"
    ResetOnDeload    bool    // Reset failure counter after deload
}
```

### 2. Trigger Type

- Uses `TriggerOnFailure` (implemented in crumb 01)
- Only fires when failure is detected

### 3. Apply Logic

When `Apply()` is called:
1. Check if `ConsecutiveFailures >= FailureThreshold`
2. If threshold met:
   - Calculate new weight based on DeloadType
   - Return ProgressionResult with negative adjustment
   - Optionally reset failure counter
3. If threshold not met:
   - Return no-op result (weight unchanged)

### 4. Validation

The `Validate()` method must ensure:
- `FailureThreshold > 0`
- `DeloadPercent` is between 0 and 1 (if percent type)
- `DeloadAmount > 0` (if fixed type)
- `DeloadType` is either "percent" or "fixed"

### 5. Factory Registration

Register `DeloadOnFailure` in the progression factory:
- Type discriminator: `"deload_on_failure"`
- Implement JSON marshaling/unmarshaling

## Example Configurations

### GZCLP T2 Style (Reset on 3 failures)
```json
{
  "type": "deload_on_failure",
  "failure_threshold": 1,
  "deload_percent": 0.15,
  "deload_type": "percent",
  "reset_on_deload": true
}
```

### Texas Method Style (5lb deload on stall)
```json
{
  "type": "deload_on_failure",
  "failure_threshold": 2,
  "deload_amount": 5,
  "deload_type": "fixed",
  "reset_on_deload": true
}
```

## Files to Create/Modify

- `internal/domain/progression/deload_on_failure.go` - New progression type
- `internal/domain/progression/deload_on_failure_test.go` - Unit tests
- `internal/domain/progression/progression.go` - Register in factory

## Acceptance Criteria

- [ ] DeloadOnFailure type implements Progression interface
- [ ] Both percent and fixed deload types work
- [ ] Threshold is respected (no deload until threshold met)
- [ ] Deload calculation is correct
- [ ] Validation catches invalid configurations
- [ ] Factory can deserialize from JSON
- [ ] Unit tests cover all deload scenarios
- [ ] Integration test shows deload after N failures

## Dependencies

- Crumb 01 (Failure Tracking and OnFailure Trigger) must be complete
- Failure counter must be accessible in trigger context

## Programs This Enables

| Program | Configuration |
|---------|---------------|
| GZCLP T2 | 1 failure -> move to next stage (StageProgression handles) |
| GZCLP T3 | Keep weight on failure (no deload, just no progress) |
| Texas Method | 2 consecutive stalls -> reduce 5-10lbs |
| Greg Nuckols Beginner | AMRAP <= target -> no weight increase (implicit via AMRAP) |
