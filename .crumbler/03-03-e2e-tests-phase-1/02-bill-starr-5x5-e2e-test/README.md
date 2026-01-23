# Bill Starr 5x5 E2E Test

## Objective
Create an end-to-end test validating the Bill Starr 5x5 program configuration and execution through the API.

## Program Characteristics
- **Heavy/Light/Medium Days**: Different intensities per day (DailyLookup pattern)
  - Heavy Day: 100% of working weight
  - Light Day: ~80% of working weight
  - Medium Day: ~90% of working weight
- **Ramp Sets**: 5x5 where weight increases each set (ramping to top set)
- **LinearProgression**: AFTER_WEEK trigger with +5lb increment

## Test Requirements

### Setup Phase
1. Create test user
2. Use seeded lifts (Squat, Bench, Deadlift)
3. Set up user's lift maxes (TRAINING_MAX)
4. Create prescriptions for each day intensity:
   - Heavy Day: PERCENT_OF at 100%, RAMP set scheme [50%, 60%, 70%, 80%, 90%]
   - Light Day: PERCENT_OF at 80%, RAMP set scheme
   - Medium Day: PERCENT_OF at 90%, RAMP set scheme
5. Create three days: Heavy, Light, Medium
6. Create 1-week cycle (Mon=Heavy, Wed=Light, Fri=Medium)
7. Create program and link cycle
8. Create LinearProgression (AFTER_WEEK trigger)
9. Link progression to program for each lift

### Execution Phase
1. Enroll user in program
2. Generate Heavy Day workout - verify:
   - RAMP sets with increasing weights
   - Top set at target percentage
3. Generate Light Day workout - verify:
   - Same RAMP pattern but at 80% intensity
4. Generate Medium Day workout - verify:
   - Same RAMP pattern but at 90% intensity
5. Complete all three workouts (full week)
6. Trigger weekly progression - verify:
   - All lifts increase by +5lb
   - Progression only fires once per week
7. Generate next Heavy Day - verify:
   - Weights increased by +5lb across the board

### Validation Phase
- RAMP sets calculate correctly with progressive weights
- Heavy/Light/Medium day intensities are accurate
- Weekly progression fires only after full week
- Weight calculation: top_set = training_max * day_intensity

## Location
`internal/api/e2e/bill_starr_test.go`

## Notes
- The RAMP set scheme should generate 5 sets with ascending weights
- Work set threshold determines which sets are "work sets"
