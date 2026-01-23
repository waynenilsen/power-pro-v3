# Wendler 5/3/1 BBB E2E Test

## Objective
Create an end-to-end test validating the Wendler 5/3/1 Boring But Big program configuration and execution through the API.

## Program Characteristics
- **4-Week Cycle**: Different rep schemes each week (WeeklyLookup pattern)
  - Week 1: 5/5/5+ at 65%/75%/85%
  - Week 2: 3/3/3+ at 70%/80%/90%
  - Week 3: 5/3/1+ at 75%/85%/95%
  - Week 4: Deload at 40%/50%/60%
- **BBB Accessory**: 5x10 at 50% after main work
- **CycleProgression**: +10lb lower body, +5lb upper body at cycle end

## Test Requirements

### Setup Phase
1. Create test user
2. Use seeded lifts (Squat, Bench, Deadlift) + create OHP
3. Set up user's TRAINING_MAX (typically 90% of 1RM)
4. Create prescriptions for each week's main work:
   - Week 1: PERCENT_OF with 85% top set, FIXED 3x5
   - Week 2: PERCENT_OF with 90% top set, FIXED 3x3
   - Week 3: PERCENT_OF with 95% top set, custom rep scheme
   - Week 4: PERCENT_OF with 60% deload
5. Create BBB accessory prescriptions (5x10 at 50%)
6. Create days for each lift focus (Squat Day, Bench Day, Deadlift Day, OHP Day)
7. Create 4-week cycle with appropriate week configurations
8. Create program and link cycle
9. Create CycleProgression with base increment
10. Link progression with override increments:
    - Squat/Deadlift: +10lb
    - Bench/OHP: +5lb

### Execution Phase
1. Enroll user in program
2. Execute Week 1 workouts - verify:
   - Main work at 65%/75%/85%
   - BBB work at 50%
3. Execute Week 2 workouts - verify:
   - Main work at 70%/80%/90%
4. Execute Week 3 workouts - verify:
   - Main work at 75%/85%/95%
5. Execute Week 4 (deload) workouts - verify:
   - Reduced intensity
6. Complete full cycle
7. Trigger cycle progression - verify:
   - Squat +10lb
   - Deadlift +10lb
   - Bench +5lb
   - OHP +5lb
8. Start next cycle - verify new training maxes apply

### Validation Phase
- Weekly prescriptions change according to wave pattern
- BBB accessory maintains consistent 5x10 at 50%
- CycleProgression applies different increments per lift
- Deload week is properly reduced intensity

## Location
`internal/api/e2e/wendler_531_test.go`

## Notes
- Training max is 90% of true 1RM (user sets this up)
- The + in 5/5/5+ indicates AMRAP (as many reps as possible) - for testing, treat as minimum reps
