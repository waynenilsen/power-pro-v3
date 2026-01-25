# 003: Session Cleanup Verification

## ERD Reference
Implements: REQ-TD2-005

## Description
Verify that session cleanup mechanisms work correctly. This includes logout invalidation, expired session purging, and cascade deletion when users are deleted.

## Context / Background
Sessions accumulate in the database over time. Proper cleanup prevents database bloat and ensures invalidated sessions cannot be reused.

## Acceptance Criteria
- [ ] Logout endpoint invalidates the session (deleted or marked invalid)
- [ ] Logged-out session token cannot be used for subsequent requests
- [ ] Expired sessions can be purged via a cleanup mechanism
- [ ] User deletion cascades to delete all user sessions
- [ ] Test verifies concurrent sessions are handled correctly
- [ ] Test verifies one session logout doesn't affect other sessions for same user

## Technical Notes
- Consider adding a session cleanup job/function if not already present
- Verify foreign key constraints handle cascade delete
- Test concurrent session scenarios (user logged in on multiple devices)
- Document the cleanup strategy (immediate vs. batch cleanup)

## Dependencies
- Blocks: None
- Blocked by: 002-session-expiration-tests
- Related: 001-auth-security-audit

## Resources / Links
- ERD: phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/erd.md
