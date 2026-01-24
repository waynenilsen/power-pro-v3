# DoubleProgression

Implement double progression - a progression type that adds reps until hitting ceiling, then adds weight and resets reps.

## Implementation

Create `/home/wayne/git/power-pro-v3/internal/domain/progression/double.go`:

1. **Add Type Constant** - Add `TypeDouble ProgressionType = "DOUBLE_PROGRESSION"` to `ValidProgressionTypes` map in `progression.go`

2. **Struct Definition**:
   ```go
   type DoubleProgression struct {
       ID               string      `json:"id"`
       Name             string      `json:"name"`
       WeightIncrement  float64     `json:"weightIncrement"`
       MaxTypeValue     MaxType     `json:"maxType"`
       TriggerTypeValue TriggerType `json:"triggerType"` // AFTER_SET
   }
   ```

3. **Methods**:
   - `Type()` - Returns TypeDouble
   - `TriggerType()` - Returns TriggerAfterSet
   - `Validate()` - Checks all fields valid, WeightIncrement > 0
   - `Apply(ctx, params)`:
     - Verify trigger type is AFTER_SET
     - Extract RepsPerformed from TriggerEvent
     - Extract MaxReps from TriggerEvent (the ceiling)
     - If reps >= MaxReps: return Applied=true, Delta=WeightIncrement
     - If reps < MaxReps: return Applied=false with reason (not at ceiling yet)
   - `MarshalJSON()` - Includes type discriminator
   - `UnmarshalDoubleProgression()` - Factory deserialization
   - `RegisterDoubleProgression()` - Factory registration

4. **Extend TriggerEvent** if needed:
   - May need to add `MaxReps *int` field to TriggerEvent for rep ceiling context

5. **Write Unit Tests** in `double_test.go`:
   - Type returns correct constant
   - TriggerType returns AFTER_SET
   - Validate accepts valid progressions
   - Validate rejects invalid progressions
   - Apply returns Applied=true when reps >= MaxReps
   - Apply returns Applied=false when reps < MaxReps
   - Apply handles missing RepsPerformed gracefully
   - JSON marshaling/unmarshaling roundtrip

## Acceptance Criteria

- DoubleProgression exists and implements Progression interface
- Apply correctly detects when user hits rep ceiling
- Apply returns weight increment when ceiling is reached
- All tests pass
