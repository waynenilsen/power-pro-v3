# Auth E2E Tests

## Ticket Reference
`phases/in-progress/002-frontend-readiness/sprints/in-progress/001-authentication-system/tickets/todo/006-auth-e2e-tests.md`

## Task
Create comprehensive E2E tests for the authentication system.

## Test Coverage

### Registration
- Successful with email, password, name
- Successful with email, password only
- Fails: missing/invalid email (400)
- Fails: missing/short password (400)
- Fails: duplicate email (409)
- Email case-insensitivity

### Login
- Successful login returns token and user
- Fails: wrong email/password (401)
- Fails: missing credentials (400)
- Token works for authenticated requests

### Logout
- Returns 204
- Invalidates session
- Without session returns 401

### Current User
- Returns user with valid session
- Returns 401 without/expired session

### Backwards Compatibility
- X-User-ID header still works
- X-Admin header still works
- Authorization takes precedence over X-User-ID

### Full Flow
- register -> login -> access protected endpoint -> logout
- User can access own resources
- User cannot access others' resources

## Acceptance Criteria
- [ ] All registration tests pass
- [ ] All login tests pass
- [ ] All logout tests pass
- [ ] All current user tests pass
- [ ] Backwards compatibility verified
- [ ] Full flow integration tests pass
- [ ] Existing E2E tests still pass

## When Done
Move ticket from `tickets/todo/` to `tickets/done/`, move sprint from `in-progress/` to `done/`, then run `crumbler delete`
