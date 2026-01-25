# Test: 010 Sheiko Beginner

## Task
Create an E2E test for the Sheiko Beginner program in `internal/api/e2e/sheiko_beginner_test.go`

## Program Characteristics
- **4-week Prep + 4-week Comp phases**: Repeatable prep, comp for peaking
- **3 days/week**: Mon/Wed/Fri
- **High frequency**: Squat 2x, Bench 3x, Deadlift 2x per week
- **Submaximal training**: Most work at 70-85% 1RM
- **Fixed sets/reps with percentage-based loading**

## Key Features to Test
1. **Prep phase structure**: 4-week repeatable block
2. **Comp phase structure**: 4-week peaking block with taper
3. **Intensity zones**: 50-59%, 60-69%, 70-79%, 80-89%, 90%+
4. **Multiple sets per exercise**: Variable sets/reps per day
5. **Phase transitions**: Prep -> Comp or Prep -> Prep

## Test Template
Similar to `sheiko_intermediate_test.go` which already tests Sheiko-style programs:
- Create prescriptions for each intensity zone
- Create 8-week cycle (4 prep + 4 comp)
- Test percentage-based weight calculations
- Test phase identification based on week number
