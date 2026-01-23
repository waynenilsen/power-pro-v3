# Sheiko Beginner E2E Test

## Objective
Create an end-to-end test validating the Sheiko Beginner program configuration and execution through the API.

## Program Characteristics
- **High Volume, Moderate Intensity**: Many sets at sub-maximal percentages
- **Fixed Sets at Various Percentages**: No ramping, just multiple straight sets
- **No Autoregulation**: Prescribed percentages only
- **Multiple Exercises Per Day**: Several variations and accessories

## Test Requirements

### Setup Phase
1. Create test user
2. Use seeded lifts (Squat, Bench, Deadlift)
3. Set up user's ONE_RM values
4. Create prescriptions representing typical Sheiko day structure:
   - Squat: 5x3 at 70%, 4x2 at 80%
   - Bench: 5x4 at 65%, 3x3 at 75%
   - Deadlift: 4x3 at 70%
   - Each as separate FIXED set schemes at PERCENT_OF ONE_RM
5. Create training days with multiple prescription blocks
6. Create multi-week cycle (Sheiko typically 4-week blocks)
7. Create program without complex progression (fixed percentages)

### Execution Phase
1. Enroll user in program
2. Generate workout for Day 1 - verify:
   - Multiple prescription blocks per lift
   - Correct percentages (e.g., 70%, 80%)
   - Correct set/rep schemes (e.g., 5x3, 4x2)
3. Verify prescription resolution accuracy:
   - 70% of 400lb 1RM = 280lb
   - 80% of 400lb 1RM = 320lb
4. Generate several days of workouts
5. Verify high volume accumulates correctly:
   - Total sets per session
   - Total reps per session

### Validation Phase
- All percentage calculations are accurate
- Multiple prescriptions per lift resolve independently
- No progression interference (Sheiko manages progression externally)
- Volume metrics (total sets/reps) match expectations

## Location
`internal/api/e2e/sheiko_beginner_test.go`

## Notes
- Sheiko programs are known for high volume with submaximal weights
- This test validates prescription resolution accuracy more than progression
- Multiple prescriptions for the same lift in one day is a key feature
- Rounding should work correctly (e.g., to nearest 5lb plate)
