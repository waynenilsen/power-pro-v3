# Batch 2: 531 Variants (3 tests)

## Overview
Update 5/3/1 program variant E2E tests to use explicit state machine flow.

## Tests to Update

1. **`wendler_531_test.go`** - 4-week cycle, AFTER_CYCLE progression
   - Remove `Force: true` from progression triggers
   - Test full 4-week cycle flow (5s, 3s, 531, deload)
   - Verify AFTER_CYCLE progression triggers at cycle completion
   - Test deload week behavior

2. **`nsuns_531_lp_test.go`** - High frequency, multiple sessions per week
   - Remove `advanceUserState()` calls
   - Test multiple workouts per day/week
   - Verify AMRAP-based progression
   - Test T1/T2 lift relationships

3. **`building_the_monolith_test.go`** - 6-week cycle, specific deload
   - Remove `Force: true` from progression triggers
   - Test 6-week periodization
   - Verify specific deload protocol
   - Test accessory work prescriptions

## Pattern to Follow

```go
// 1. Enroll and verify initial state
enrollResp := enrollUser(t, ts, userID, programID)

// 2. Work through cycle weeks
for week := 1; week <= totalWeeks; week++ {
    for day := range daysInWeek {
        sessionID := startWorkoutSession(t, ts, userID)
        logSets(t, ts, sessionID, sets)
        finishWorkoutSession(t, ts, sessionID, userID)
    }
}

// 3. Verify cycle completion triggers progression
enrollmentState := getEnrollment(t, ts, userID)
// Assert cycle status is BETWEEN_CYCLES or new cycle started

// 4. Verify progression applied
workoutResp := getWorkout(t, ts, userID)
// Assert training maxes increased
```

## Acceptance Criteria
- [ ] No `advanceUserState()` usage
- [ ] No `Force: true` in progression triggers
- [ ] Tests verify week-by-week state transitions
- [ ] AFTER_CYCLE progression works automatically
- [ ] Deload weeks handled correctly
