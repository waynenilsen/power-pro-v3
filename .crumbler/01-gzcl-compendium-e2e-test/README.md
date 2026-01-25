# GZCL Compendium E2E Test

Create E2E test for the GZCL Compendium (VDIP) program.

## Reference Files
- Program spec: `programs/016-gzcl-compendium.md`
- Test pattern: `internal/api/e2e/jacked_and_tan_test.go` (similar GZCL structure)
- Existing tests: `internal/api/e2e/gzclp_t1_test.go`

## Program Characteristics (VDIP)
- 5 days per week, ongoing (no fixed duration)
- Three-tier system: T1, T2, T3
- T1: 3 MRS (Max Rep Sets) at 85% TM
- T2: 3 MRS at 65% TM
- T3: 4 MRS each

## VDIP Progression Rules
- T1: 15+ total reps = +10lb, 10-14 = +5lb, <10 = maintain
- T2: 30+ total reps = +10lb, 25-29 = +5lb, <25 = maintain
- T3: 50+ total reps = +5lb, <50 = maintain

## Test Coverage
1. Verify 5-day week structure with T1/T2/T3 exercises
2. Verify MRS prescription generation
3. Verify VDIP progression logic based on total reps
4. Verify weight calculations at correct percentages
