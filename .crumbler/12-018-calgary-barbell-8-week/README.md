# Test: 018 Calgary Barbell 8-Week

## Task
Create an E2E test for the Calgary Barbell 8-Week program in `internal/api/e2e/calgary_barbell_8_test.go`

## Program Characteristics
- **8-week + taper**: Condensed peaking program
- **4 days/week**: Squat/Bench, DL/Bench, etc.
- **2 phases**: Accumulation (1-4), Intensification/Peaking (5-8)
- **Dual-tier work**: Heavy work + volume back-offs
- **RPE-based top sets in Phase 2**: E1RM calculation for back-offs

## Key Features to Test
1. **Phase transition at week 5**: Accumulation -> Intensification
2. **Percentage-based loading Phase 1**: Fixed percentages
3. **RPE-based loading Phase 2**: Top sets based on RPE
4. **Back-off calculations**: Volume sets after heavy work
5. **Fatigue notation**: 1+2F, 1+3R sets
6. **Taper week**: Competition prep in final week

## Test Template
Shorter version of `calgary_barbell_16_test.go`:
- Create WeeklyLookup for 8-week progression
- Create prescriptions for heavy + volume tiers
- Create 8-week cycle with 2 phases
- Test week 4 -> week 5 transition
- Test taper week structure
