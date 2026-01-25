# Batch 1: LP Programs (5 tests)

## Overview
Update Linear Progression (LP) program E2E tests to use explicit state machine flow.

## Tests to Update

1. **`starting_strength_test.go`** - SESSION progression (every workout)
   - Remove `advanceUserState()` calls
   - Use explicit workout start/finish flow
   - Verify AFTER_SESSION progression triggers automatically

2. **`greyskull_lp_test.go`** - AMRAP-driven progression
   - Remove `Force: true` from progression triggers
   - Log AMRAP sets with performance data
   - Verify progression based on AMRAP performance

3. **`bill_starr_test.go`** - Classic 5x5 progression
   - Remove manual state advancement
   - Replace `Force: true` with event-driven progression
   - Verify weekly linear progression flow

4. **`gzclp_t1_test.go`** - T1 stage progression
   - Remove `advanceUserState()` and `Force: true`
   - Test T1 progression through stages (5x3 → 6x2 → 10x1)
   - Verify stage transitions on failure

5. **`gzclp_t2_test.go`** - T2 stage progression
   - Remove `advanceUserState()` and `Force: true`
   - Test T2 progression through stages (3x10 → 3x8 → 3x6)
   - Verify stage transitions and weight increases

## Pattern to Follow

```go
// 1. Enroll and verify initial state
enrollResp := enrollUser(t, ts, userID, programID)

// 2. Start workout - triggers auto-transitions
sessionID := startWorkoutSession(t, ts, userID)

// 3. Log sets with performance data
logSets(t, ts, sessionID, sets)

// 4. Finish workout - verify auto-progression
finishWorkoutSession(t, ts, sessionID, userID)

// 5. Get next workout - verify progression applied
workoutResp := getWorkout(t, ts, userID)
// Assert increased weights
```

## Acceptance Criteria
- [ ] No `advanceUserState()` usage in any test
- [ ] No `Force: true` in progression triggers
- [ ] Each test verifies state at each step
- [ ] All existing assertions still pass
- [ ] Tests demonstrate full workout lifecycle
