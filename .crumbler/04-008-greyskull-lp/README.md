# Test: 008 Greyskull LP

## Task
Create an E2E test for the Greyskull LP program in `internal/api/e2e/greyskull_lp_test.go`

## Program Characteristics
- **3-day A/B rotation**: Week 1 (A/B/A), Week 2 (B/A/B)
- **AMRAP final sets**: Every main lift ends with AMRAP (2x5, 1x5+)
- **Autoregulated progression**: Double increment for 2x target reps
- **Built-in deload**: 10% reset on failure (less than target reps)
- **Deadlift once weekly**: Only on Day 2 each week

## Key Features to Test
1. **A/B rotation pattern**: Bench/Row vs OHP/Chinup alternation
2. **AMRAP prescription**: 2 fixed sets + 1 AMRAP
3. **AMRAPProgression with double increment**: 5-9 reps = +2.5lb, 10+ reps = +5lb
4. **DeloadOnFailure**: <5 reps triggers 10% reset
5. **Deadlift frequency**: Verify deadlift only appears on specific days

## Test Template
Combines patterns from `nsuns_531_lp_test.go` (AMRAP) and `starting_strength_test.go` (A/B rotation):
- Create A/B days with appropriate exercises
- Create 2-week cycle with A/B/A and B/A/B pattern
- Create AMRAPProgression with threshold-based double increment
- Test standard progression (5-9 reps)
- Test double progression (10+ reps)
- Test deload on failure (<5 reps)
