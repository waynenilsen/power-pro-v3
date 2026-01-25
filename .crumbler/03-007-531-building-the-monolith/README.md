# Test: 007 5/3/1 Building the Monolith

## Task
Create an E2E test for the 5/3/1 Building the Monolith program in `internal/api/e2e/building_the_monolith_test.go`

## Program Characteristics
- **6-week cycle**: Two 3-week blocks (same percentages, TM increases in block 2)
- **3 days/week**: Monday (Squat+Press), Wednesday (Deadlift+Bench), Friday (Squat+Press volume)
- **WeeklyLookup**: Week 1 (70/80/90%), Week 2 (65/75/85%), Week 3 (75/85/95%)
- **CycleProgression**: After 3 weeks, +5lb upper, +10lb lower

## Key Features to Test
1. **3-week wave pattern**: Test weight calculations for weeks 1, 2, 3
2. **High volume main work**: 8 sets (3 warmup + 5x5 at top percentage)
3. **Friday Widowmaker**: 1x20 at lower percentage (45-55% TM)
4. **Press AMRAP on Monday**: Final set is 5+
5. **Friday Press volume**: 10-15 sets of 5 at moderate percentage
6. **Cycle progression**: TM increase after week 3

## Test Template
Similar to `wendler_531_test.go`:
- Create WeeklyLookup for 3-week percentages
- Create prescriptions for main work and Widowmaker
- Create 3 days with appropriate prescriptions
- Create 6-week (2x3 week) cycle
- Test week 1, 2, 3 weight generation
- Test cycle progression at week 3->4 boundary
