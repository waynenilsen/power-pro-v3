# E2E Tests for AMRAP Programs

Create end-to-end tests demonstrating the three AMRAP-enabled programs from the acceptance criteria.

## Programs to Test

### 1. Wendler 5/3/1 BBB (`internal/api/e2e/wendler_531_test.go`)

**Program characteristics:**
- 4-week cycle (weeks 1-3 working, week 4 deload)
- AMRAP on final set of weeks 1-3 (5+, 3+, 1+)
- WeeklyLookup for the 4-week wave (65/75/85%, 70/80/90%, 75/85/95%, 40/50/60%)
- CycleProgression: +5lb upper, +10lb lower at cycle end

**Test flow:**
1. Create program with cycle, weeks, days, prescriptions
2. Create user and set Training Maxes
3. Enroll user in program
4. Generate Week 1 workout - verify AMRAP set (85% x 5+)
5. Log AMRAP set with 8 reps
6. Advance through weeks 1-4
7. Verify cycle progression triggers at week 4→1 transition
8. Generate new workout - verify weights increased

### 2. Greg Nuckols High Frequency (`internal/api/e2e/nuckols_frequency_test.go`)

**Program characteristics:**
- 3-week cycle
- AMAP (AMRAP) sets in Week 3 at 85% (Thursday/Friday)
- DailyLookup for day-specific intensities (Mon=75%, Tue=80%, etc.)
- WeeklyLookup for volume progression
- CycleProgression based on Week 3 AMAP performance

**Test flow:**
1. Create program structure with lookups
2. Enroll user, set maxes
3. Progress through weeks 1-2 with standard sets
4. Week 3: Log AMAP set with 7 reps
5. Verify progression applies at cycle end
6. Verify Week 1 of new cycle has updated weights

### 3. nSuns 5/3/1 LP 5-Day (`internal/api/e2e/nsuns_531_lp_test.go`)

**Program characteristics:**
- 1-week cycle with weekly linear progression
- Multiple AMRAP sets per day (1+ sets on primary lifts)
- AMRAPProgression with threshold-based increments
- AfterSet trigger for immediate progression

**Test flow:**
1. Create 5-day program structure
2. Set up AMRAPProgression with thresholds (2-3 reps=+5lb, 4-5=+10lb, 6+=+15lb)
3. Enroll user, set Training Maxes
4. Generate Day 1 workout (Bench with 1+ set)
5. Log AMRAP with 5 reps
6. Verify progression triggered immediately (+10lb)
7. Generate Day 2 workout
8. Log Squat AMRAP with 3 reps
9. Verify Squat TM increased (+5lb)
10. Complete week, verify all lifts progressed appropriately

## Test Patterns to Follow

Reference existing E2E tests:
- `internal/api/e2e/starting_strength_test.go`
- `internal/api/e2e/bill_starr_test.go`

**Each test should:**
1. Set up complete program structure via API
2. Create user and establish maxes
3. Enroll in program
4. Generate workouts and verify structure
5. Log sets (including AMRAP)
6. Verify progression triggers
7. Verify updated maxes affect future workouts

## Files to Create

- `internal/api/e2e/wendler_531_test.go`
- `internal/api/e2e/nuckols_frequency_test.go`
- `internal/api/e2e/nsuns_531_lp_test.go`

## Verification

- `go test ./internal/api/e2e/...` passes
- All three programs demonstrate:
  - AMRAP set generation
  - Logged set persistence
  - Progression based on AMRAP performance
