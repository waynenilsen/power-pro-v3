# RTS Intermediate E2E Test

Create E2E test for the RTS (Reactive Training Systems) Generalized Intermediate program.

## Reference Files
- Program spec: `programs/019-rts-intermediate.md`
- Test pattern: `internal/api/e2e/sheiko_intermediate_test.go` (similar intermediate level)

## Program Characteristics
- 9 weeks with mesocycle structure
- RPE-based load autoregulation
- Fatigue percentage system for back-off work
- Volume tracking by movement category (Squat, Press, Pull)

## RPE-to-Percentage Chart
Used to convert RPE + reps to % of 1RM:
- 1 rep @ RPE 10 = 100%
- 5 reps @ RPE 9 = 80%
- etc.

## Fatigue Methods
1. Load Drop: Work up to set, reduce weight by fatigue %, continue until target RPE
2. Repeat Sets: Same weight until RPE increases
3. Rep Drop: Same weight, fewer reps

## Weekly Structure
- Week 1: Baseline (RPE 9-10, no fatigue method)
- Weeks 2-4: Development (Load Drop, 0-7.5% fatigue)
- Weeks 5-8: Intensification (Load Drop + Repeat, 5-7% fatigue)
- Week 9: Peaking/Testing (singles at RPE 7-10, 0% fatigue)

## Test Coverage
1. Verify 9-week cycle structure
2. Verify RPE-based weight calculation using lookup table
3. Verify fatigue percentage back-off calculations
4. Verify 3-weight warm-up progression (90%, 95%, 100%)
5. Verify phase transitions
