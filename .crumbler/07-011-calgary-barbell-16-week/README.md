# Test: 011 Calgary Barbell 16-Week

## Task
Create an E2E test for the Calgary Barbell 16-Week program in `internal/api/e2e/calgary_barbell_16_test.go`

## Program Characteristics
- **16-week meet prep**: 15 training + 1 taper week
- **4 days/week**: Squat/Bench, DL/Bench/Squat, Squat/Bench/DL, DL/Bench
- **Phase structure**: Hypertrophy (1-4), Strength (5-8), Peaking (9-11), Intensification (12-15), Taper (16)
- **Mixed loading**: Percentage-based and RPE-based
- **Fatigue/Repeat sets**: 1+2F, 1+3R notation

## Key Features to Test
1. **5-phase periodization**: Verify phase transitions at weeks 5, 9, 12, 16
2. **Percentage progression**: Week 1 = 67%, Week 2 = 70%, etc.
3. **Rep scheme progression**: Week 1 = 7 reps, Week 2 = 6 reps, etc.
4. **WeeklyLookup for intensity/reps**: Different values each week
5. **Taper week**: Reduced volume in week 16

## Test Template
Complex multi-phase program similar to `sheiko_intermediate_test.go`:
- Create WeeklyLookup for 16-week progression
- Create prescriptions for each phase style
- Create 16-week cycle with phase-specific days
- Test intensity/rep progression across weeks
- Test taper week volume reduction
