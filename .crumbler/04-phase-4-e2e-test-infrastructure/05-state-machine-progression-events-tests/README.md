# State Machine Progression Events Tests

## Overview

Create `state_machine_progression_events_test.go` to verify progression events fire at correct times.

NOTE: If the event/progression trigger system isn't fully implemented yet, these tests may need to be simplified or marked as TODO for future implementation.

## Test Cases

### Test: AFTER_SET progression fires on set log
- Set up program with AFTER_SET progression trigger
- Start workout, log a set
- Verify progression was applied (check lift max changed)

### Test: AFTER_SESSION progression fires on workout finish
- Set up program with AFTER_SESSION progression trigger
- Complete a workout session
- Verify progression was applied

### Test: AFTER_WEEK progression fires on week complete
- Set up program with AFTER_WEEK progression trigger
- Complete all workouts in week
- Advance week
- Verify progression was applied

### Test: AFTER_CYCLE progression fires on cycle complete
- Set up program with AFTER_CYCLE progression trigger
- Complete full cycle
- Verify progression was applied

### Test: ON_FAILURE progression fires on failed set
- Set up program with ON_FAILURE progression/deload trigger
- Log a failed set (reps < target)
- Verify deload was applied (lift max decreased)

## File Location

`internal/api/e2e/state_machine_progression_events_test.go`

## Acceptance Criteria

- [ ] Test cases implemented for available progression triggers
- [ ] If triggers not implemented, document what's missing
- [ ] Verify progression application through lift max changes
- [ ] Tests serve as documentation for how progressions integrate with state machine
