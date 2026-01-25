# 004: Phase 2 Test Coverage Review

## ERD Reference
Implements: REQ-TD2-006, REQ-TD2-007, REQ-TD2-008

## Description
Review and improve test coverage for all Phase 2 code including auth service, profile service, and dashboard service. Target is >90% coverage for all new domain logic.

## Context / Background
Phase 2 introduced significant new functionality. Comprehensive test coverage ensures regressions are caught early and provides confidence for future development.

## Acceptance Criteria

### Auth Service (REQ-TD2-006)
- [ ] Coverage report shows > 90% for auth service
- [ ] Registration success path tested
- [ ] Registration failure paths tested (duplicate email, invalid input)
- [ ] Login success path tested
- [ ] Login failure paths tested (wrong password, unknown user)
- [ ] Logout functionality tested
- [ ] Session creation and validation tested

### Profile Service (REQ-TD2-007)
- [ ] Coverage report shows > 90% for profile service
- [ ] Profile read operation tested
- [ ] Profile update operation tested
- [ ] Authorization checks tested (can only access own profile)
- [ ] Invalid profile data handling tested

### Dashboard Service (REQ-TD2-008)
- [ ] Coverage report shows > 90% for dashboard service
- [ ] Enrollment aggregation tested
- [ ] Next workout calculation tested
- [ ] Current session query tested
- [ ] Recent workouts query tested
- [ ] Current maxes query tested
- [ ] Empty state (no enrollments) handled and tested

## Technical Notes
- Run coverage with `go test -coverprofile=coverage.out ./...`
- View coverage report with `go tool cover -html=coverage.out`
- Focus on branch coverage, not just line coverage
- Add tests for error paths, not just happy paths

## Dependencies
- Blocks: None
- Blocked by: Auth, Profile, Dashboard services complete
- Related: 001-auth-security-audit

## Resources / Links
- ERD: phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/erd.md
