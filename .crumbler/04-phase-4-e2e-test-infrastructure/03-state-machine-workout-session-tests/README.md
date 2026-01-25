# State Machine Workout Session Tests

## Overview

Create `state_machine_workout_session_test.go` to verify workout session state transitions.

## Test Cases

### Test: Start workout creates IN_PROGRESS session
- Enroll user
- Start workout session
- Verify session status = "IN_PROGRESS"
- Verify enrollment has currentWorkoutSession populated
- Verify cycleStatus transitions to "IN_PROGRESS" if was "PENDING"
- Verify weekStatus transitions to "IN_PROGRESS" if was "PENDING"

### Test: Can't start second workout while one IN_PROGRESS
- Enroll user
- Start first workout session
- Attempt to start second workout session
- Verify error response

### Test: Finish workout transitions to COMPLETED
- Start workout session
- Log at least one set
- Finish workout session
- Verify session status = "COMPLETED"
- Verify currentWorkoutSession is nil in enrollment

### Test: Abandon workout transitions to ABANDONED
- Start workout session
- Abandon workout (if endpoint exists, otherwise skip)
- Verify session status = "ABANDONED"

### Test: Can start new workout after COMPLETED
- Complete one workout
- Start new workout session
- Verify new session created successfully

### Test: Can start new workout after ABANDONED
- Abandon workout (if possible)
- Start new workout session
- Verify new session created successfully

### Test: Logging sets requires active session
- Enroll user (no session started)
- Attempt to log sets
- Verify error response

### Test: Can't log sets to COMPLETED session
- Complete a workout session
- Attempt to log sets to that session ID
- Verify error response

## File Location

`internal/api/e2e/state_machine_workout_session_test.go`

## Acceptance Criteria

- [ ] All applicable test cases implemented
- [ ] Tests use shared helpers
- [ ] Session state transitions verified
- [ ] Error cases validated
