# Batch 4: Frequency/RPE Programs (4 tests)

## Overview
Update high-frequency and RPE-based program E2E tests to use explicit state machine flow.

## Tests to Update

1. **`nuckols_frequency_test.go`** - 3x/week per lift, AMRAP-based
   - Remove `Force: true` from progression triggers
   - Test 3x weekly frequency per main lift
   - Verify AMRAP-based progression
   - Test different intensity/volume days

2. **`nuckols_beginner_test.go`** - Beginner frequency
   - Remove `advanceUserState()` and `Force: true`
   - Test beginner-appropriate frequency
   - Verify faster progression rate
   - Test technique focus prescriptions

3. **`rts_intermediate_test.go`** - RPE-based, fatigue management
   - Remove `Force: true` from progression triggers
   - Test RPE-based autoregulation
   - Verify fatigue management triggers
   - Test load adjustments based on RPE feedback

4. **`nsuns_cap3_test.go`** - CAP3 periodization
   - Remove `advanceUserState()` and `Force: true`
   - Test Cyclical AMRAP Progression
   - Verify 3-week cycles
   - Test AMRAP performance tracking

## Pattern to Follow

```go
// 1. Enroll and verify initial state
enrollResp := enrollUser(t, ts, userID, programID)

// 2. Work through high-frequency week
for day := range daysInWeek {
    sessionID := startWorkoutSession(t, ts, userID)

    // Log sets with RPE/AMRAP data as applicable
    for _, set := range sets {
        logSet(t, ts, sessionID, set)
    }

    finishWorkoutSession(t, ts, sessionID, userID)
}

// 3. Verify progression based on performance data
workoutResp := getWorkout(t, ts, userID)
// Assert weights adjusted based on AMRAP/RPE
```

## Acceptance Criteria
- [ ] No `advanceUserState()` usage
- [ ] No `Force: true` in progression triggers
- [ ] RPE-based tests log RPE values
- [ ] AMRAP tests log rep performance
- [ ] Fatigue management triggers work automatically
