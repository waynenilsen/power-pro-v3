# Batch 5: Peaking/Other Programs (6 tests)

## Overview
Update peaking and miscellaneous program E2E tests to use explicit state machine flow.

## Tests to Update

1. **`calgary_barbell_8_test.go`** - Peaking program, meet_date handling
   - Remove `Force: true` from progression triggers
   - Test 8-week peaking phase
   - Verify meet_date countdown
   - Test competition-specific progressions

2. **`calgary_barbell_16_test.go`** - Longer peaking, phase transitions
   - Remove `advanceUserState()` and `Force: true`
   - Test 16-week periodization
   - Verify phase transitions (hypertrophy → strength → peaking)
   - Test volume/intensity waves

3. **`sheiko_beginner_test.go`** - Percentage-based, no AMRAP
   - Remove `Force: true` from progression triggers
   - Test strict percentage-based prescriptions
   - Verify no AMRAP elements
   - Test block transitions

4. **`sheiko_intermediate_test.go`** - Intermediate periodization
   - Remove `advanceUserState()` and `Force: true`
   - Test intermediate-level periodization
   - Verify competition prep cycles
   - Test multiple preparation blocks

5. **`reddit_ppl_test.go`** - 6-day rotation
   - Remove `Force: true` from progression triggers
   - Test Push/Pull/Legs rotation
   - Verify 6-day weekly pattern
   - Test linear progression on compounds

6. **`phase5_rotation_e2e_test.go`** - Rotation-based programs
   - Remove `advanceUserState()` and `Force: true`
   - Test A/B or other rotation patterns
   - Verify day rotation logic
   - Test multi-week rotations

## Pattern to Follow

```go
// 1. Enroll and verify initial state
enrollResp := enrollUser(t, ts, userID, programID)

// 2. For peaking programs, work through prep phases
for phase := range phases {
    for week := range weeksInPhase {
        for day := range daysInWeek {
            sessionID := startWorkoutSession(t, ts, userID)
            logSets(t, ts, sessionID, sets)
            finishWorkoutSession(t, ts, sessionID, userID)
        }
    }
}

// 3. Verify phase transitions and meet prep
enrollmentState := getEnrollment(t, ts, userID)
// Assert progression toward meet date
```

## Acceptance Criteria
- [ ] No `advanceUserState()` usage
- [ ] No `Force: true` in progression triggers
- [ ] Peaking programs track meet_date correctly
- [ ] Phase transitions work automatically
- [ ] Rotation-based programs cycle correctly
