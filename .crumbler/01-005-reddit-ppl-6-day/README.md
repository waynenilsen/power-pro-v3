# Test: 005 Reddit PPL 6-Day

## Task
Create an E2E test for the Reddit PPL 6-Day program in `internal/api/e2e/reddit_ppl_test.go`

## Program Characteristics
- **6-day split**: Push/Pull/Legs x 2 per week + 1 rest day
- **Alternating Primary Lifts**: Different primary compounds on A vs B days
- **Linear Progression**: Main compounds with AMRAP sets (1x5+, 4x5+1x5+, 2x5+1x5+)
- **Double Progression**: Accessories with 8-12 rep ranges

## Key Features to Test
1. **6-day weekly cycle structure**: Pull A -> Push A -> Legs A -> Pull B -> Push B -> Legs B
2. **AMRAP prescriptions**: Deadlift 1x5+, Bench 4x5+1x5+, Squat 2x5+1x5+, Rows 4x5+1x5+, OHP 4x5+1x5+
3. **Linear Progression**: +2.5lb upper body, +5lb lower body per session
4. **Day rotation**: Verify correct workout generation for each day

## Test Template
Follow the pattern in `wendler_531_test.go` and `starting_strength_test.go`:
- Create test server
- Create lift maxes (squat, bench, deadlift, OHP, rows)
- Create prescriptions for each exercise type
- Create 6 days with appropriate prescriptions
- Create 1-week cycle with 6 days
- Create program and link progressions
- Enroll user and verify workout generation
- Test AMRAP logging and progression triggers
