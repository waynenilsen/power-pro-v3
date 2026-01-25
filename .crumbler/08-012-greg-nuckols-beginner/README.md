# Test: 012 Greg Nuckols Beginner

## Task
Create an E2E test for the Greg Nuckols Beginner program in `internal/api/e2e/nuckols_beginner_test.go`

## Program Characteristics
- **3 days/week**: Different lift combinations each day
- **Frequency optimization**: Bench 3x/week, Squat 2x/week, Deadlift 2x/week
- **AMAP (As Many As Possible) sets**: 2 fixed sets + 1 AMAP
- **Daily undulation**: Day 1 = 8 reps, Day 2 = 6 reps, Day 3 = 4 reps
- **4-week cycles for squat/deadlift**: Intensity progressions

## Key Features to Test
1. **Daily rep variation**: 70% x8, 75% x6, 80% x4 for bench
2. **AMAP final sets**: 2 fixed + 1 AMAP structure
3. **AMRAPProgression**: Progress based on AMAP performance
4. **Multi-lift days**: Squat+Bench+DL on Day 1
5. **4-week periodization for squat/deadlift**

## Test Template
Similar to `greyskull_lp_test.go` (AMAP sets) and `nuckols_frequency_test.go`:
- Create DailyLookup for rep/intensity variation
- Create prescriptions with 2 fixed + 1 AMAP
- Create 3-day week with multi-lift days
- Test AMAP logging and progression
- Test 4-week cycle for squat/deadlift
