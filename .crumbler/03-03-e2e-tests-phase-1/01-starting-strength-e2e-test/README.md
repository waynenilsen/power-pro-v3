# Starting Strength E2E Test

## Objective
Create an end-to-end test validating the Starting Strength program configuration and execution through the API.

## Program Characteristics
- **A/B Rotation**: Alternating workouts (A: Squat/Bench/Deadlift, B: Squat/Press/Power Clean)
- **Fixed 3x5**: All main lifts use FIXED set scheme with 3 sets of 5 reps
- **LinearProgression**: AFTER_SESSION trigger with +5lb increment for upper, +10lb for lower

## Test Requirements

### Setup Phase
1. Create test user
2. Create lifts (use seeded IDs for Squat/Bench/Deadlift, create Press/Power Clean)
3. Set up user's lift maxes (TRAINING_MAX for each lift)
4. Create prescriptions for each exercise:
   - LoadStrategy: PERCENT_OF with TRAINING_MAX reference at 100%
   - SetScheme: FIXED with sets=3, reps=5
5. Create Day A and Day B
6. Create 1-week cycle with A/B/A pattern (Mon/Wed/Fri)
7. Create program and link cycle
8. Create LinearProgression (AFTER_SESSION trigger)
9. Link progression to program for each lift

### Execution Phase
1. Enroll user in program
2. Generate workout for Day A - verify:
   - Correct exercises (Squat, Bench, Deadlift)
   - 3x5 sets at correct weights
3. Complete Day A workout
4. Trigger progression - verify:
   - Only performed lifts increase
   - Increment is +5lb for upper, +10lb for lower
5. Generate workout for Day B - verify:
   - Different exercises (Squat, Press, Power Clean)
   - Squat weight increased by +10lb
6. Complete Day B workout
7. Trigger progression - verify correct increments applied

### Validation Phase
- All prescriptions resolve correctly
- Progression fires only for lifts actually performed
- Weight increments are accurate
- A/B rotation works correctly

## Location
`internal/api/e2e/starting_strength_test.go`

## Test File Setup
- Use `testutil.NewTestServer()` for isolated database
- Use helper functions for API calls (authPost, adminPost, etc.)
- Follow existing integration_test.go patterns
