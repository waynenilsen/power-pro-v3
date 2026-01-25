# Phase 4: E2E Test Infrastructure

## Overview

Add test helpers and infrastructure for verifying state machine behavior in E2E tests. These helpers will be used by all program E2E tests.

## Tasks

### 1. Add state assertion helpers

Create/update `internal/api/e2e/helpers.go` (or similar):

```go
// assertEnrollmentState verifies all state fields in enrollment response
func assertEnrollmentState(t *testing.T, resp *http.Response, expected EnrollmentState)

// startWorkoutAndVerify starts a workout and verifies state transitions
func startWorkoutAndVerify(t *testing.T, ts *TestServer, userID string,
    expectedCycleStatus, expectedWeekStatus string) string

// finishWorkoutAndVerify finishes a workout and verifies state
func finishWorkoutAndVerify(t *testing.T, ts *TestServer, sessionID string,
    expectedEnrollmentStatus, expectedCycleStatus, expectedWeekStatus string)

// completeFullCycle runs through an entire cycle verifying states
func completeFullCycle(t *testing.T, ts *TestServer, userID string,
    program ProgramConfig) CycleCompletionResult
```

### 2. Add response types with state fields

Update E2E response types to include:
- `enrollment_status`
- `cycle_status`
- `week_status`
- `current_workout_session`

### 3. Create state machine specific E2E tests

Create new test files:

#### `state_machine_enrollment_test.go`
- Test: NONE → ACTIVE (enroll)
- Test: ACTIVE → BETWEEN_CYCLES (cycle completes)
- Test: BETWEEN_CYCLES → ACTIVE (start new cycle)
- Test: ACTIVE → QUIT (quit)
- Test: BETWEEN_CYCLES → QUIT (quit while deciding)
- Test: Can't start workout when BETWEEN_CYCLES
- Test: Can't start new cycle when ACTIVE

#### `state_machine_workout_session_test.go`
- Test: Start workout creates IN_PROGRESS session
- Test: Can't start second workout while one IN_PROGRESS
- Test: Finish workout transitions to COMPLETED
- Test: Abandon workout transitions to ABANDONED
- Test: Can start new workout after COMPLETED
- Test: Can start new workout after ABANDONED
- Test: Logging sets requires active session
- Test: Can't log sets to COMPLETED session

#### `state_machine_week_cycle_test.go`
- Test: First workout of week: PENDING → IN_PROGRESS
- Test: Last workout of week: → COMPLETED, auto-advance
- Test: Last workout of cycle: triggers CYCLE_COMPLETED
- Test: Week status resets to PENDING on new week
- Test: Cycle status resets to PENDING on new cycle

#### `state_machine_progression_events_test.go`
- Test: AFTER_SET progression fires on set log
- Test: AFTER_SESSION progression fires on workout finish
- Test: AFTER_WEEK progression fires on week complete
- Test: AFTER_CYCLE progression fires on cycle complete
- Test: ON_FAILURE progression fires on failed set

## Acceptance Criteria

- [ ] State assertion helpers implemented
- [ ] Response types include state fields
- [ ] `state_machine_enrollment_test.go` passes
- [ ] `state_machine_workout_session_test.go` passes
- [ ] `state_machine_week_cycle_test.go` passes
- [ ] `state_machine_progression_events_test.go` passes
- [ ] All tests serve as documentation for frontend team
