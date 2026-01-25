# Ticket 004: Phase 2 Test Coverage Review

## ERD Reference
Implements: REQ-TD2-006, REQ-TD2-007, REQ-TD2-008

## Description
Review and improve test coverage for all Phase 2 code. Target is >90% coverage.

## Acceptance Criteria

### Auth Service (REQ-TD2-006)
- [ ] Coverage > 90% for auth service
- [ ] Registration success/failure paths tested
- [ ] Login success/failure paths tested
- [ ] Logout functionality tested
- [ ] Session creation and validation tested

### Profile Service (REQ-TD2-007)
- [ ] Coverage > 90% for profile service
- [ ] Profile read/update operations tested
- [ ] Authorization checks tested
- [ ] Invalid data handling tested

### Dashboard Service (REQ-TD2-008)
- [ ] Coverage > 90% for dashboard service
- [ ] Enrollment aggregation tested
- [ ] Next workout calculation tested
- [ ] Empty state handled and tested

## Technical Notes
- Run coverage with `go test -coverprofile=coverage.out ./...`
- Focus on branch coverage, not just line coverage
- Add tests for error paths, not just happy paths

## Dependencies
- Blocked by: 003-session-cleanup-verification
