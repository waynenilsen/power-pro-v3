# State Assertion Helpers

## Overview

Create state assertion helpers and update E2E response types to support state machine testing.

## Tasks

### 1. Update EnrollmentData response type

Add state fields to the existing `EnrollmentData` struct in `starting_strength_test.go`:

```go
type EnrollmentData struct {
    ID               string                        `json:"id"`
    UserID           string                        `json:"userId"`
    Program          EnrollmentProgramData         `json:"program"`
    State            EnrollmentStateData           `json:"state"`
    EnrollmentStatus string                        `json:"enrollmentStatus"`
    CycleStatus      string                        `json:"cycleStatus"`
    WeekStatus       string                        `json:"weekStatus"`
    CurrentWorkoutSession *CurrentWorkoutSessionData `json:"currentWorkoutSession"`
}

type EnrollmentProgramData struct {
    ID               string  `json:"id"`
    Name             string  `json:"name"`
    Slug             string  `json:"slug"`
    CycleLengthWeeks int     `json:"cycleLengthWeeks"`
}

type EnrollmentStateData struct {
    CurrentWeek           int  `json:"currentWeek"`
    CurrentCycleIteration int  `json:"currentCycleIteration"`
    CurrentDayIndex       *int `json:"currentDayIndex,omitempty"`
}

type CurrentWorkoutSessionData struct {
    ID         string `json:"id"`
    WeekNumber int    `json:"weekNumber"`
    DayIndex   int    `json:"dayIndex"`
    Status     string `json:"status"`
}
```

### 2. Create state assertion helpers

Add these helpers to `starting_strength_test.go`:

```go
// ExpectedEnrollmentState defines expected state values for assertions.
type ExpectedEnrollmentState struct {
    EnrollmentStatus string
    CycleStatus      string
    WeekStatus       string
    CurrentWeek      int
    CycleIteration   int
    HasActiveSession bool
    SessionStatus    string // Optional, only checked if HasActiveSession
}

// assertEnrollmentState verifies all state fields in an enrollment response.
func assertEnrollmentState(t *testing.T, enrollment EnrollmentData, expected ExpectedEnrollmentState)

// getEnrollment fetches the current enrollment state for a user.
func getEnrollment(t *testing.T, ts *testutil.TestServer, userID string) EnrollmentData

// enrollUser enrolls a user in a program and returns the enrollment data.
func enrollUser(t *testing.T, ts *testutil.TestServer, userID, programID string) EnrollmentData

// unenrollUser removes enrollment for a user.
func unenrollUser(t *testing.T, ts *testutil.TestServer, userID string)

// startWorkoutAndVerify starts a workout and verifies state transitions.
func startWorkoutAndVerify(t *testing.T, ts *testutil.TestServer, userID string,
    expectedCycleStatus, expectedWeekStatus string) string

// finishWorkoutAndVerify finishes a workout and verifies state transitions.
func finishWorkoutAndVerify(t *testing.T, ts *testutil.TestServer, sessionID, userID string,
    expectedEnrollmentStatus, expectedCycleStatus, expectedWeekStatus string)

// advanceWeek advances to the next week and returns updated enrollment.
func advanceWeek(t *testing.T, ts *testutil.TestServer, userID string) EnrollmentData

// startNextCycle starts a new cycle when in BETWEEN_CYCLES state.
func startNextCycle(t *testing.T, ts *testutil.TestServer, userID string) EnrollmentData
```

## Acceptance Criteria

- [ ] EnrollmentData updated with all state fields
- [ ] Supporting types (EnrollmentProgramData, etc.) added
- [ ] assertEnrollmentState helper implemented
- [ ] getEnrollment helper implemented
- [ ] enrollUser helper implemented
- [ ] unenrollUser helper implemented
- [ ] startWorkoutAndVerify helper implemented
- [ ] finishWorkoutAndVerify helper implemented
- [ ] advanceWeek helper implemented
- [ ] startNextCycle helper implemented
- [ ] Existing E2E tests still pass
