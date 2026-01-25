# Test: 009 GZCL Jacked and Tan 2.0

## Task
Create an E2E test for the GZCL Jacked and Tan 2.0 program in `internal/api/e2e/jacked_and_tan_test.go`

## Program Characteristics
- **12-week cycle**: Four 3-week mesocycles (A, B, C, D)
- **4-5 days/week**: Squat, Bench, Deadlift, OHP focus days
- **Tiered system**: T1 (competition lifts), T2a (primary accessories), T2b/T2c, T3
- **RM Finding**: Progressive RM testing (10RM -> 8RM -> 6RM -> etc.)
- **WeeklyLookup**: Different RM targets per week

## Key Features to Test
1. **RM progression**: Week 1 = 10RM, Week 2 = 8RM, etc.
2. **Back-off sets**: Sets after RM finding at reduced intensity
3. **Multi-tier prescription**: T1, T2a, T2b, T3 exercises per day
4. **Block periodization**: Verify mesocycle transitions

## Test Template
Similar to `nuckols_frequency_test.go` (multi-week cycle):
- Create WeeklyLookup for RM targets per week
- Create tiered prescriptions for T1, T2a, T2b, T3
- Create 12-week cycle with 4 mesocycles
- Test RM target changes across weeks
- Test progression at mesocycle boundaries
