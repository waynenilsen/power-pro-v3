# Sprint 005: Technical Debt - Phase 2 Cleanup

## Overview

Mandatory 5th sprint technical debt paydown. Review security of authentication implementation, verify session management, ensure test coverage, and update API documentation.

## Reference Documents

- PRD: `phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/prd.md`
- ERD: `phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/erd.md`
- Tickets: `phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/tickets/todo/`

## Requirements Summary

### Security Audit
1. **Password hashing review** - Verify bcrypt/argon2 with proper cost, no plaintext logging, timing-safe comparison
2. **Session token review** - Verify cryptographic randomness, sufficient entropy (>=256 bits), timing-safe validation, hashed storage
3. **Authorization middleware review** - All protected endpoints require valid session, proper user context, expired sessions rejected

### Session Management
- Test session expiration handling
- Verify session cleanup on logout
- Verify cascade delete on user deletion

### Test Coverage
- Auth service: >90% coverage
- Profile service: >90% coverage
- Dashboard service: >90% coverage

### Documentation
- Update API docs for auth endpoints (register, login, logout)
- Update API docs for profile endpoints
- Update API docs for dashboard endpoint

## Tickets (in order)

1. `001-auth-security-audit.md` - Password, token, and middleware security review
2. `002-session-expiration-tests.md` - Session expiration verification tests
3. `003-session-cleanup-verification.md` - Session cleanup and cascade delete
4. `004-phase2-test-coverage-review.md` - Test coverage for Phase 2 code
5. `005-api-documentation-update.md` - API documentation sync

## Dependencies

- Sprints 001-004 must be complete

## Constraints

- No API contract changes
- No schema changes
- No new features

## Success Criteria

- Security audit items reviewed and addressed
- Session edge cases verified
- Test coverage meets targets
- API docs complete for new endpoints
