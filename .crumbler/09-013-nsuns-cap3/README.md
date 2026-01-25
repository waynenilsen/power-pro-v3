# Test: 013 nSuns CAP3

## Task
Create an E2E test for the nSuns CAP3 program in `internal/api/e2e/nsuns_cap3_test.go`

## Program Characteristics
- **3-week rotating cycle**: Each lift gets AMRAP test once per 3 weeks
- **Cyclical AMRAP rotation**: Week 1 = DL AMRAP, Week 2 = Squat AMRAP, Week 3 = Bench AMRAP
- **Training Max based**: TM = 90% of estimated 1RM
- **6 days/week**: Different lifts each day
- **EMOM protocol**: Secondary lifts use EMOM sets

## Key Features to Test
1. **Cyclical rotation**: Verify which lift gets high-intensity AMRAP each week
2. **Training Max calculations**: Weight = TM * percentage
3. **Volume vs Intensity phases**: Medium volume when not AMRAP testing
4. **3-week WeeklyLookup**: Different programming per week
5. **Dual progression**: Major (AMRAP test) + secondary (regular AMRAP)

## Test Template
Complex rotation similar to `nuckols_frequency_test.go`:
- Create WeeklyLookup for 3-week cycle
- Create prescriptions for AMRAP test days vs volume days
- Create 3-week cycle with 6 days per week
- Test rotation: Week 1 DL test, Week 2 Squat test, Week 3 Bench test
- Test progression after AMRAP performance
