# 006: E2E Tests for Program Discovery

## ERD Reference
Implements: All requirements (verification)

## Description
Create comprehensive end-to-end tests for all program discovery features. Tests should verify filtering, search, combined queries, and detail enhancement against a seeded test database.

## Context / Background
E2E tests ensure the full stack works correctly from HTTP request to database query to response. With multiple filter parameters and combinations, thorough testing is essential to prevent regressions and ensure all edge cases are handled.

## Acceptance Criteria
- [ ] Create E2E test file: `e2e/program_discovery_test.go`
- [ ] Test filtering by difficulty:
  - Filter beginner returns correct programs
  - Filter intermediate returns correct programs
  - Filter advanced returns empty (no advanced programs seeded)
  - Invalid difficulty returns 400
- [ ] Test filtering by days_per_week:
  - Filter 3 days returns Starting Strength, Texas Method
  - Filter 4 days returns 5/3/1, GZCLP
  - Filter 5 days returns empty
  - Invalid value (0, 8, -1, "abc") returns 400
- [ ] Test filtering by focus:
  - Filter strength returns all 4 canonical programs
  - Filter hypertrophy returns empty
  - Invalid focus returns 400
- [ ] Test filtering by has_amrap:
  - Filter true returns 5/3/1, GZCLP
  - Filter false returns Starting Strength, Texas Method
  - Invalid boolean returns 400
- [ ] Test combined filters:
  - difficulty=beginner&days_per_week=3 returns Starting Strength only
  - difficulty=intermediate&has_amrap=true returns 5/3/1 only
  - No matches returns 200 with empty array
- [ ] Test search:
  - search=strength returns Starting Strength
  - search=531 returns 5/3/1
  - search=gz returns GZCLP
  - search=method returns Texas Method
  - Case-insensitive: search=STRENGTH returns Starting Strength
- [ ] Test search + filters:
  - search=strength&difficulty=beginner returns Starting Strength
  - search=5&difficulty=intermediate returns 5/3/1
- [ ] Test program detail enhancements:
  - Detail response includes sampleWeek array
  - Detail response includes liftRequirements array
  - Detail response includes estimatedSessionMinutes
  - liftRequirements is sorted alphabetically
  - sampleWeek has correct structure for each canonical program
- [ ] Test metadata in list response:
  - List response includes difficulty, daysPerWeek, focus, hasAmrap
  - Values match backfilled data

## Technical Notes
- Use existing E2E test patterns from `e2e/` directory
- Seed test data or rely on canonical programs from Sprint 004
- Use table-driven tests for filter validation cases
- Verify response JSON structure matches DTOs
- Consider test helpers for common assertions
- Ensure tests clean up or use isolated database state

## Dependencies
- Blocks: None (final ticket)
- Blocked by: 001, 002, 003, 004, 005 (all features must be implemented)

## Resources / Links
- Existing E2E tests: `e2e/` directory
- Test patterns: look for existing program tests
- ERD requirements: all sections
