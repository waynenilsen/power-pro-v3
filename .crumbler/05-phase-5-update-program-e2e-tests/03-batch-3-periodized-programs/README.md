# Batch 3: Periodized Programs (4 tests)

## Overview
Update periodized program E2E tests to use explicit state machine flow.

## Tests to Update

1. **`texas_method_test.go`** - AFTER_WEEK progression on intensity day
   - Remove `Force: true` from progression triggers
   - Test Volume/Recovery/Intensity day pattern
   - Verify AFTER_WEEK progression on intensity day completion
   - Test PR attempts on intensity day

2. **`inverted_juggernaut_test.go`** - Wave progression with AMRAP
   - Remove `advanceUserState()` and `Force: true`
   - Test inverted rep scheme (10s, 8s, 5s, 3s waves)
   - Verify AMRAP-based progression between waves
   - Test realization weeks

3. **`jacked_and_tan_test.go`** - Block periodization
   - Remove `Force: true` from progression triggers
   - Test block transitions (accumulation → intensification → realization)
   - Verify T1/T2/T3 progression patterns
   - Test rep max testing protocols

4. **`gzcl_compendium_test.go`** - VDIP approach
   - Remove `Force: true` from progression triggers
   - Test Volume-Dependent Intensity Progression
   - Verify MRS (max rep sets) handling
   - Test weight increases based on volume targets

## Pattern to Follow

```go
// 1. Enroll and verify initial state
enrollResp := enrollUser(t, ts, userID, programID)

// 2. Work through periodization blocks
for block := range blocks {
    for week := range weeksInBlock {
        for day := range daysInWeek {
            sessionID := startWorkoutSession(t, ts, userID)
            // Log sets appropriate to the block phase
            logSets(t, ts, sessionID, sets)
            finishWorkoutSession(t, ts, sessionID, userID)
        }
    }
    // Verify block transition occurred
}

// 3. Verify progression based on periodization model
workoutResp := getWorkout(t, ts, userID)
```

## Acceptance Criteria
- [ ] No `advanceUserState()` usage
- [ ] No `Force: true` in progression triggers
- [ ] Tests verify phase/block transitions
- [ ] Wave/block periodization works correctly
- [ ] AMRAP-based progressions trigger automatically
