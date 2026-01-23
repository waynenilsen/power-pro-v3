# Greg Nuckols High Frequency E2E Test

## Objective
Create an end-to-end test validating the Greg Nuckols High Frequency program configuration and execution through the API.

## Program Characteristics
- **High Frequency**: Same lifts multiple times per week
- **DailyLookup + WeeklyLookup**: Different intensities by day AND by week
- **CycleProgression**: Weight increases at end of cycle
- **Varied Set/Rep Schemes**: Different prescriptions per day

## Test Requirements

### Setup Phase
1. Create test user
2. Use seeded lifts (Squat, Bench, Deadlift)
3. Set up user's TRAINING_MAX values
4. Create prescriptions for various day types:
   - Heavy Day: 3x3 at 85%
   - Volume Day: 5x5 at 75%
   - Light Day: 3x8 at 65%
5. Create multiple training days per week (e.g., 4-5 days)
6. Configure DailyLookup for day-of-week variation:
   - Monday: Heavy Squat, Volume Bench
   - Tuesday: Light Squat, Heavy Deadlift
   - Thursday: Volume Squat, Heavy Bench
   - Friday: Light Bench, Volume Deadlift
7. Create multi-week cycle (e.g., 3-4 weeks)
8. Configure WeeklyLookup for weekly wave (if applicable):
   - Week 1: Standard intensities
   - Week 2: +5% across the board
   - Week 3: Peak week
9. Create CycleProgression
10. Link progression to program for each lift

### Execution Phase
1. Enroll user in program
2. Generate Monday workout - verify:
   - Heavy Squat prescription (3x3 at 85%)
   - Volume Bench prescription (5x5 at 75%)
3. Generate Tuesday workout - verify:
   - Light Squat (3x8 at 65%)
   - Heavy Deadlift (3x3 at 85%)
4. Complete full week of training
5. Move to Week 2 - verify:
   - Intensities adjusted per WeeklyLookup
6. Complete full cycle
7. Trigger cycle progression - verify:
   - All lifts increase appropriately
8. Start next cycle - verify new training maxes apply

### Validation Phase
- DailyLookup correctly varies prescriptions by day
- WeeklyLookup correctly varies prescriptions by week
- Same lift can appear multiple times per week with different prescriptions
- CycleProgression applies correctly at cycle end
- Prescription resolution handles combined lookups

## Location
`internal/api/e2e/greg_nuckols_hf_test.go`

## Notes
- High frequency = same lift 3-4x per week
- Key challenge is handling multiple prescriptions for same lift on different days
- The lookup system allows variation without creating duplicate prescriptions
- This tests the most complex configuration scenario
