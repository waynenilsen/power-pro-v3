# Modify Logged Sets Handler

## Overview

Update the logged sets endpoint to require active workout session and emit events.

## Endpoint to Update

```
POST /sessions/{id}/sets
```

## Implementation Details

### Current Behavior
- Logs sets against a session
- Tracks failure counters

### Required Changes

1. **Require Active Workout Session**
   - Validate session is IN_PROGRESS
   - Return error if session is COMPLETED or ABANDONED

2. **Emit SET_LOGGED Event**
   - After successfully logging each set
   - Include in payload:
     - `loggedSetId`
     - `liftId`
     - `repsPerformed`
     - `targetReps`
     - `weight`
     - `isAMRAP`
     - `isFailure`

3. **Failure Detection for ON_FAILURE Trigger**
   - Detect if set is a failure (reps < targetReps)
   - Include `isFailure` flag in event
   - This enables ON_FAILURE progression triggers

## Files to Modify

- `internal/api/logged_set_handler.go` - Update handler
- `internal/api/logged_set_handler_test.go` - Update tests

## Acceptance Criteria

- [ ] Cannot log sets to non-IN_PROGRESS session
- [ ] SET_LOGGED event emitted for each logged set
- [ ] Event payload includes all required fields
- [ ] Failure detection works correctly
- [ ] Unit tests pass
