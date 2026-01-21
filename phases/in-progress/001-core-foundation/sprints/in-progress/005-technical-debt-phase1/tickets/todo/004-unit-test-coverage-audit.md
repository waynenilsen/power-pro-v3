# 004: Unit Test Coverage Audit and Gap Filling

## ERD Reference
Implements: REQ-DEBT-004

## Description
Audit unit test coverage for domain logic and fill gaps to achieve >90% coverage. High test coverage prevents regressions during future development.

## Context / Background
This is a MUST-priority requirement. The domain packages contain critical business logic that must be thoroughly tested to support future development with confidence.

## Acceptance Criteria
- [ ] Generate coverage report for all domain packages
- [ ] Identify coverage gaps in domain logic
- [ ] Write unit tests to fill coverage gaps
- [ ] Coverage report shows >90% for domain packages
- [ ] Critical paths are tested (prescription resolution, workout generation, progression)
- [ ] All tests pass consistently

## Technical Notes
- Use `go test -coverprofile` to generate coverage reports
- Focus on internal/domain packages
- Test both happy paths and edge cases
- Use table-driven tests where appropriate
- Mock dependencies for isolated unit testing
- Priority areas: prescription resolvers, progression evaluators, workout generators

## Dependencies
- Blocks: 005-integration-test-review (should have unit tests first)
- Blocked by: None
- Related: 006-flaky-test-resolution

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/005-technical-debt-phase1/erd.md
