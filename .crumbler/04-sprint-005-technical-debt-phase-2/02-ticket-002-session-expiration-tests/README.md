# Ticket 002: Session Expiration Tests

## ERD Reference
Implements: REQ-TD2-004

## Description
Verify that session expiration is handled correctly by adding comprehensive tests for expiration scenarios.

## Acceptance Criteria
- [ ] Test verifies expired sessions are rejected with 401 Unauthorized
- [ ] Test verifies session works just before expiration boundary
- [ ] Test verifies session fails just after expiration boundary
- [ ] Test verifies expiration time stored in session is respected
- [ ] Test verifies session cannot be "refreshed" after expiration
- [ ] All expiration tests pass consistently (no timing-dependent flakiness)

## Technical Notes
- Use time mocking or short expiration times for testing
- Consider using a clock interface for deterministic testing
- Test both the service layer and the HTTP layer
- Ensure tests don't depend on wall-clock time

## Dependencies
- Blocked by: 001-auth-security-audit (security review first)
