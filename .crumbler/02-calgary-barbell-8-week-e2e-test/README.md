# Calgary Barbell 8-Week E2E Test

Create E2E test for the Calgary Barbell 8-Week peaking program.

## Reference Files
- Program spec: `programs/018-calgary-barbell-8-week.md`
- Test pattern: `internal/api/e2e/calgary_barbell_16_test.go` (16-week version)

## Program Characteristics
- 8 weeks plus taper week (condensed version of 16-week)
- 4 training days per week
- Two distinct 4-week phases:
  - Phase 1 (Weeks 1-4): Accumulation - percentage-based work, dual-tier heavy/volume
  - Phase 2 (Weeks 5-8): Intensification - RPE top sets + back-off percentages
- Taper week: countdown format (5 days out, 4 days out, etc.)

## Intensity Progression
Phase 1 Heavy: 80%, 82%, 86%, 85%
Phase 1 Volume: 68%, 70%, 72%, 75%
Phase 2 Top Sets: RPE 8 (3 reps) -> RPE 8 (2 reps) -> RPE 8-9 (1 rep) -> RPE 8-9 (1 rep)
Phase 2 Back-Off: 65%, 68%, 72%, 76% of E1RM

## Test Coverage
1. Verify 8-week cycle structure with 4 days/week
2. Verify Phase 1 percentage calculations
3. Verify Phase 2 RPE-based top set structure
4. Verify phase transition at week 5
5. Verify taper week reduced volume
