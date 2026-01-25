# State Machine Enrollment Tests

## Overview

Create `state_machine_enrollment_test.go` to verify enrollment state transitions.

## Test Cases

### Test: NONE → ACTIVE (enroll)
- Enroll new user in program
- Verify enrollmentStatus = "ACTIVE"
- Verify cycleStatus = "PENDING"
- Verify weekStatus = "PENDING"
- Verify currentWeek = 1, cycleIteration = 1

### Test: ACTIVE → BETWEEN_CYCLES (cycle completes)
- Enroll user, complete all weeks in cycle
- Use advanceWeek to reach final week
- Advance past final week
- Verify enrollmentStatus = "BETWEEN_CYCLES"
- Verify cycleStatus = "COMPLETED"

### Test: BETWEEN_CYCLES → ACTIVE (start new cycle)
- Get user to BETWEEN_CYCLES state
- Call startNextCycle
- Verify enrollmentStatus = "ACTIVE"
- Verify cycleStatus = "PENDING"
- Verify weekStatus = "PENDING"
- Verify cycleIteration incremented

### Test: ACTIVE → QUIT (quit)
- Enroll user
- Unenroll user (DELETE /users/{userId}/program)
- Verify user is no longer enrolled

### Test: BETWEEN_CYCLES → QUIT (quit while deciding)
- Get user to BETWEEN_CYCLES state
- Unenroll user
- Verify user is no longer enrolled

### Test: Can't start workout when BETWEEN_CYCLES
- Get user to BETWEEN_CYCLES state
- Attempt to start workout
- Verify error response with appropriate message

### Test: Can't start new cycle when ACTIVE
- Enroll user (starts ACTIVE)
- Attempt to call startNextCycle
- Verify error response (invalid state transition)

## File Location

`internal/api/e2e/state_machine_enrollment_test.go`

## Acceptance Criteria

- [ ] All 7 test cases implemented
- [ ] Tests use shared helpers from 01-state-assertion-helpers
- [ ] Tests pass consistently
- [ ] Error responses validated for invalid transitions
