# StageProgression Implementation

Implement `StageProgression` - a progression type that changes set/rep schemes on failure.

## Context

GZCLP is the canonical example of stage-based progression:
- **T1 Default**: 5x3+ -> 6x2+ -> 10x1+ (on failure, move to next stage)
- **T2 Default**: 3x10 -> 3x8 -> 3x6 (on failure, move to next stage)

When the lifter fails to hit the prescribed volume, they move to the next stage which uses a different set/rep scheme at the SAME weight. After exhausting all stages, they reset (typically with DeloadOnFailure or manual intervention).

## Implementation Requirements

### 1. Stage Definition

```go
type Stage struct {
    Name       string    // e.g., "5x3+", "6x2+", "10x1+"
    Sets       int       // Number of sets
    Reps       int       // Reps per set (target)
    IsAMRAP    bool      // If true, last set is AMRAP
    MinVolume  int       // Minimum total reps to pass (e.g., 15 for 5x3)
}
```

### 2. StageProgression Type

Create `internal/domain/progression/stage.go`:

```go
type StageProgression struct {
    Stages            []Stage // Ordered list of stages
    CurrentStage      int     // Current stage index (0-based)
    ResetOnExhaustion bool    // If true, reset to stage 0 after last stage fails
    DeloadOnReset     bool    // If true, apply deload when resetting
    DeloadPercent     float64 // Deload amount when resetting (if DeloadOnReset)
}
```

### 3. Trigger Type

- Uses `TriggerOnFailure` (from crumb 01)
- Fires when failure is detected

### 4. Apply Logic

When `Apply()` is called with OnFailure trigger:
1. Check if lifter failed current stage (volume < MinVolume)
2. If failed:
   - If not at last stage: advance `CurrentStage += 1`
   - If at last stage:
     - If `ResetOnExhaustion`: reset to stage 0, optionally deload
     - Otherwise: return error/special result for manual intervention
3. Return result indicating:
   - The new stage
   - Weight adjustment (0 unless deload on reset)
   - Whether SetScheme should change (always yes on stage change)

### 5. SetScheme Integration

StageProgression must communicate the new SetScheme to use:
- Option A: Return stage info in ProgressionResult, service updates prescription
- Option B: StageProgression manages a reference to the prescription's SetScheme

Recommend Option A - keep Progressions focused on calculating changes.

### 6. Progression Result Enhancement

May need to extend `ProgressionResult`:
```go
type ProgressionResult struct {
    // Existing fields...
    NewStage       *int    // If set, update to this stage
    NewSetScheme   *SetScheme // If set, replace prescription's SetScheme
}
```

### 7. Tracking Current Stage

The current stage must be persisted:
- Option A: Add `stage` field to prescription
- Option B: Use a separate `stage_tracker` table
- Option C: Store in progression_logs and derive from history

Recommend Option A for simplicity - prescriptions already exist.

## GZCLP Stage Configurations

### T1 Default Scheme
```json
{
  "type": "stage_progression",
  "stages": [
    {"name": "5x3+", "sets": 5, "reps": 3, "is_amrap": true, "min_volume": 15},
    {"name": "6x2+", "sets": 6, "reps": 2, "is_amrap": true, "min_volume": 12},
    {"name": "10x1+", "sets": 10, "reps": 1, "is_amrap": true, "min_volume": 10}
  ],
  "current_stage": 0,
  "reset_on_exhaustion": true,
  "deload_on_reset": true,
  "deload_percent": 0.15
}
```

### T1 Modified Scheme
```json
{
  "type": "stage_progression",
  "stages": [
    {"name": "3x5+", "sets": 3, "reps": 5, "is_amrap": true, "min_volume": 15},
    {"name": "4x3+", "sets": 4, "reps": 3, "is_amrap": true, "min_volume": 12},
    {"name": "5x2+", "sets": 5, "reps": 2, "is_amrap": true, "min_volume": 10}
  ],
  "current_stage": 0,
  "reset_on_exhaustion": true,
  "deload_on_reset": true,
  "deload_percent": 0.15
}
```

### T2 Default Scheme (No AMRAP)
```json
{
  "type": "stage_progression",
  "stages": [
    {"name": "3x10", "sets": 3, "reps": 10, "is_amrap": false, "min_volume": 30},
    {"name": "3x8", "sets": 3, "reps": 8, "is_amrap": false, "min_volume": 24},
    {"name": "3x6", "sets": 3, "reps": 6, "is_amrap": false, "min_volume": 18}
  ],
  "current_stage": 0,
  "reset_on_exhaustion": true,
  "deload_on_reset": false
}
```

## Files to Create/Modify

- `internal/domain/progression/stage.go` - Stage and StageProgression types
- `internal/domain/progression/stage_test.go` - Unit tests
- `internal/domain/progression/progression.go` - Register in factory
- `internal/domain/progression/progression_result.go` - Extend result if needed
- Possibly: `migrations/00017_add_stage_to_prescriptions.sql` - If tracking stage in prescription

## Acceptance Criteria

- [ ] Stage struct captures all needed info (sets, reps, AMRAP, min volume)
- [ ] StageProgression implements Progression interface
- [ ] Stages advance correctly on failure
- [ ] Reset to stage 0 works when exhausted (if configured)
- [ ] Deload on reset works (if configured)
- [ ] Current stage persists between sessions
- [ ] SetScheme changes are communicated correctly
- [ ] Unit tests cover all stage transitions
- [ ] Integration test shows full GZCLP T1 cycle: 5x3+ fail -> 6x2+ fail -> 10x1+ fail -> reset

## Dependencies

- Crumb 01 (Failure Tracking and OnFailure Trigger) must be complete
- SetScheme domain (Fixed, AMRAP) already exists

## Notes

- StageProgression is the most complex progression type so far
- Consider whether stages should be immutable after creation
- The "success" path (weight increase) is handled by LinearProgression or AMRAPProgression, not StageProgression
- StageProgression only handles the FAILURE path (stage changes)
