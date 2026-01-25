# Test: 019 RTS Intermediate

## Task
Create an E2E test for the RTS Generalized Intermediate program in `internal/api/e2e/rts_intermediate_test.go`

## Program Characteristics
- **RPE-based autoregulation**: Primary load selection mechanism
- **Fatigue percentage calculations**: Manage training stress
- **RPE-to-Percentage lookup**: Converts RPE + reps to % of 1RM
- **Volume-load targets**: Weekly volume goals per movement class
- **Exercise variation support**: Different 1RMs for variants (squat, squat pause, etc.)

## Key Features to Test
1. **RPE prescription**: Target RPE values for sets
2. **RPE-to-percentage lookup**: RPE 9 x 5 reps = 80% 1RM
3. **Fatigue percentage drops**: Back-off sets from top sets
4. **Volume-load tracking**: Sets x reps x weight
5. **Movement classification**: Squat/Press/Pull categories
6. **Variant 1RM support**: Different maxes for variations

## Test Template
Unique RPE-based program - may need new helpers:
- Create RPE lookup table (if API supports)
- Create prescriptions with RPE targets
- Test RPE-to-weight conversion
- Test fatigue percentage calculations
- Test volume-load tracking
