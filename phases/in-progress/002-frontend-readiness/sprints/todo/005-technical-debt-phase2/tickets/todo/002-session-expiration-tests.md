# 002: Session Expiration Tests

## ERD Reference
Implements: REQ-TD2-004

## Description
Verify that session expiration is handled correctly by adding comprehensive tests for expiration scenarios. Sessions must expire to limit the window of exposure if a token is compromised.

## Context / Background
Session tokens have a limited lifetime for security reasons. This ticket ensures the expiration logic is correctly implemented and tested, covering edge cases around expiration timing.

## Acceptance Criteria
- [ ] Test verifies that expired sessions are rejected with 401 Unauthorized
- [ ] Test verifies session works just before expiration boundary
- [ ] Test verifies session fails just after expiration boundary
- [ ] Test verifies expiration time stored in session record is respected
- [ ] Test verifies session cannot be "refreshed" after expiration
- [ ] All expiration tests pass consistently (no timing-dependent flakiness)

## Technical Notes
- Use time mocking or short expiration times for testing
- Consider using a clock interface for deterministic testing
- Test both the service layer and the HTTP layer
- Ensure tests don't depend on wall-clock time

## Dependencies
- Blocks: None
- Blocked by: 001-auth-security-audit (security review first)
- Related: 003-session-cleanup-verification

## Resources / Links
- ERD: phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/erd.md
