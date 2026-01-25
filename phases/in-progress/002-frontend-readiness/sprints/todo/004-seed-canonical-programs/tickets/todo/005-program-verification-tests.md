# 005: Program Seed Verification Tests

## ERD Reference
Implements: REQ-VERIFY-001, REQ-VERIFY-002

## Description
Create comprehensive tests that verify all seeded canonical programs have correct structure and accurate prescriptions. These tests run after migrations to ensure data integrity.

## Context / Background
Seed migrations are critical data - if prescriptions are wrong, users will train incorrectly. Verification tests catch errors before they affect users. Tests should validate both structural correctness (right number of days, weeks) and data accuracy (correct percentages, rep schemes).

## Acceptance Criteria
- [ ] Create test file: `internal/program/canonical_test.go` (or appropriate location)
- [ ] Starting Strength verification:
  - [ ] Program exists with slug `starting-strength`
  - [ ] Has exactly 2 workout days
  - [ ] Day A has Squat (3x5), Bench (3x5), Deadlift (1x5)
  - [ ] Day B has Squat (3x5), Press (3x5), Power Clean (5x3)
  - [ ] Progression increments are correct (+5/+5/+5/+10/+5)
- [ ] Texas Method verification:
  - [ ] Program exists with slug `texas-method`
  - [ ] Has exactly 3 workout days
  - [ ] Volume Day has correct exercises and percentages (90%)
  - [ ] Recovery Day has correct exercises and percentages (80%)
  - [ ] Intensity Day has correct exercises (100%)
  - [ ] Weekly progression model configured
- [ ] Wendler 5/3/1 verification:
  - [ ] Program exists with slug `531`
  - [ ] Has exactly 4 workout days
  - [ ] Has exactly 4 weeks per cycle
  - [ ] Week 1 percentages: 65%, 75%, 85%
  - [ ] Week 2 percentages: 70%, 80%, 90%
  - [ ] Week 3 percentages: 75%, 85%, 95%
  - [ ] Week 4 percentages: 40%, 50%, 60%
  - [ ] AMRAP sets marked correctly
  - [ ] Cycle progression: +5 upper, +10 lower
- [ ] GZCLP verification:
  - [ ] Program exists with slug `gzclp`
  - [ ] Has exactly 4 workout days
  - [ ] T1/T2 pairings are correct per day
  - [ ] T1 default scheme is 5x3+
  - [ ] T2 default scheme is 3x10
  - [ ] Progression increments: +5 lower, +2.5 upper
- [ ] All tests pass with fresh database after migrations
- [ ] Tests are included in CI pipeline

## Technical Notes
- Tests should query database directly, not through API (unit tests)
- Use test helpers to run migrations before tests if needed
- Consider table-driven tests for prescription verification
- Example structure:
```go
func TestCanonicalPrograms(t *testing.T) {
    tests := []struct {
        slug       string
        name       string
        days       int
        // ...
    }{
        {"starting-strength", "Starting Strength", 2},
        {"texas-method", "Texas Method", 3},
        {"531", "Wendler 5/3/1", 4},
        {"gzclp", "GZCLP", 4},
    }
    // ...
}
```
- May need separate tests for structure vs. prescriptions vs. progression

## Dependencies
- Blocks: None (final verification)
- Blocked by: 001, 002, 003, 004 (all program seeds must exist)

## Resources / Links
- Program specifications: `programs/*.md`
- Go testing patterns: existing tests in codebase
