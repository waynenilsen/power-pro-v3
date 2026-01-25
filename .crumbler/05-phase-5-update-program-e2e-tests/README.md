# Phase 5: Update Program E2E Tests

## Overview

Update all 22 program E2E tests to use explicit state machine flow. Remove `advanceUserState()` and `Force: true` patterns. Each test should demonstrate the full workout lifecycle.

## Key Changes for All Tests

### Remove
- `advanceUserState()` helper usage - replace with explicit workout start/finish
- `Force: true` in progression triggers - test event-driven behavior
- Direct DB manipulation - everything through API

### Add
- Explicit workout session lifecycle (start → log sets → finish)
- State verification at each step
- Auto-progression verification (no manual triggers)
- Between-cycles handling where applicable

## Test Batches

### Batch 1: LP Programs (5 tests)
- `starting_strength_test.go` - SESSION progression (every workout)
- `greyskull_lp_test.go` - AMRAP-driven progression
- `bill_starr_test.go` - Classic 5x5 progression
- `gzclp_t1_test.go` - T1 stage progression
- `gzclp_t2_test.go` - T2 stage progression

### Batch 2: 531 Variants (3 tests)
- `wendler_531_test.go` - 4-week cycle, AFTER_CYCLE progression
- `nsuns_531_lp_test.go` - High frequency, multiple sessions per week
- `building_the_monolith_test.go` - 6-week cycle, specific deload

### Batch 3: Periodized Programs (4 tests)
- `texas_method_test.go` - AFTER_WEEK progression on intensity day
- `inverted_juggernaut_test.go` - Wave progression with AMRAP
- `jacked_and_tan_test.go` - Block periodization
- `gzcl_compendium_test.go` - VDIP approach

### Batch 4: Frequency/RPE Programs (4 tests)
- `nuckols_frequency_test.go` - 3x/week per lift, AMRAP-based
- `nuckols_beginner_test.go` - Beginner frequency
- `rts_intermediate_test.go` - RPE-based, fatigue management
- `nsuns_cap3_test.go` - CAP3 periodization

### Batch 5: Peaking/Other Programs (6 tests)
- `calgary_barbell_8_test.go` - Peaking program, meet_date handling
- `calgary_barbell_16_test.go` - Longer peaking, phase transitions
- `sheiko_beginner_test.go` - Percentage-based, no AMRAP
- `sheiko_intermediate_test.go` - Intermediate periodization
- `reddit_ppl_test.go` - 6-day rotation
- `phase5_rotation_e2e_test.go` - Rotation-based programs

## Test Pattern Template

Each updated test should follow this pattern:

```go
// 1. Enroll and verify initial state
enrollResp := createEnrollment(t, ts, userID, programID)
assertState(t, enrollResp, State{
    EnrollmentStatus: "ACTIVE",
    CycleStatus:      "PENDING",
    WeekStatus:       "PENDING",
})

// 2. Start workout - triggers auto-transitions
startResp := startWorkout(t, ts, userID)
assertState(t, startResp, State{
    CycleStatus: "IN_PROGRESS",
    WeekStatus:  "IN_PROGRESS",
})
sessionID := startResp.SessionID

// 3. Log sets
logSets(t, ts, sessionID, sets)

// 4. Finish workout - check for auto-progression
finishResp := finishWorkout(t, ts, sessionID)
// Verify state transitions based on position in program

// 5. Continue through cycle...
// 6. Verify BETWEEN_CYCLES and progression applied
```

## Acceptance Criteria

- [ ] All 22 program E2E tests updated
- [ ] No usage of `advanceUserState()` helper
- [ ] No usage of `Force: true` in progressions
- [ ] All tests verify state at each step
- [ ] All tests demonstrate workout session lifecycle
- [ ] All existing test assertions still pass
- [ ] Tests serve as documentation for frontend team
