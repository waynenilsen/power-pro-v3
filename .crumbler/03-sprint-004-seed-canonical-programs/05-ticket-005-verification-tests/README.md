# Ticket 005: Program Verification Tests

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/tickets/todo/005-program-verification-tests.md`
- ERD: REQ-VERIFY-001, REQ-VERIFY-002

## Task
Create comprehensive tests that verify all seeded canonical programs have correct structure and prescriptions.

## Test Coverage

### Starting Strength
- Program exists with slug `starting-strength`
- Has 2 workout days (A/B)
- Day A: Squat 3x5, Bench 3x5, Deadlift 1x5
- Day B: Squat 3x5, Press 3x5, Power Clean 5x3
- Progression: +5/+5/+5/+10/+5

### Texas Method
- Program exists with slug `texas-method`
- Has 3 workout days
- Volume Day: correct exercises @ 90%
- Recovery Day: correct exercises @ 80%
- Intensity Day: correct exercises @ 100%

### Wendler 5/3/1
- Program exists with slug `531`
- Has 4 workout days
- Has 4 weeks per cycle
- Week 1-4 percentages correct
- AMRAP sets marked

### GZCLP
- Program exists with slug `gzclp`
- Has 4 workout days
- T1/T2 pairings correct
- T1 = 5x3+, T2 = 3x10

## Acceptance Criteria
- [ ] Create test file `internal/program/canonical_test.go`
- [ ] Test program existence by slug
- [ ] Test correct number of days
- [ ] Test prescription accuracy
- [ ] Test progression rules
- [ ] All tests pass with fresh database

## Technical Notes
- Query database directly (unit tests)
- Table-driven tests for multiple programs
- Run after migrations
