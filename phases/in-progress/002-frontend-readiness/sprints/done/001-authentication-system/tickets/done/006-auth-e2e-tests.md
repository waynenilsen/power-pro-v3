# 006: Auth E2E Tests

## ERD Reference
Implements: All REQ-* requirements (validation through E2E testing)

## Description
Create comprehensive E2E tests for the authentication system, verifying the full flow from registration through authenticated API access.

## Context / Background
E2E tests ensure the entire authentication system works correctly in integration. They should cover happy paths, error cases, and edge cases. Existing E2E tests should continue to work with X-User-ID headers.

## Acceptance Criteria

### Registration Tests
- [ ] Test successful registration with email, password, and name
- [ ] Test successful registration with email and password only (no name)
- [ ] Test registration fails with missing email (400)
- [ ] Test registration fails with invalid email format (400)
- [ ] Test registration fails with missing password (400)
- [ ] Test registration fails with short password (400)
- [ ] Test registration fails with duplicate email (409)
- [ ] Test registration email is case-insensitive (user@EXAMPLE.com = user@example.com)

### Login Tests
- [ ] Test successful login returns token and user
- [ ] Test login fails with wrong email (401)
- [ ] Test login fails with wrong password (401)
- [ ] Test login fails with missing credentials (400)
- [ ] Test login creates session in database
- [ ] Test returned token works for authenticated requests

### Logout Tests
- [ ] Test successful logout returns 204
- [ ] Test logout invalidates session (subsequent requests fail)
- [ ] Test logout without session returns 401

### Current User Tests
- [ ] Test GET /auth/me returns user with valid session
- [ ] Test GET /auth/me returns 401 without session
- [ ] Test GET /auth/me returns 401 with expired session

### Session Tests
- [ ] Test session token works in Authorization: Bearer header
- [ ] Test invalid token returns 401
- [ ] Test expired session returns 401

### Backwards Compatibility Tests
- [ ] Test existing E2E tests still work with X-User-ID header
- [ ] Test X-Admin header still works for admin operations
- [ ] Test Authorization header takes precedence over X-User-ID

### Integration Tests
- [ ] Test full flow: register -> login -> access protected endpoint -> logout
- [ ] Test authenticated user can access their own resources
- [ ] Test authenticated user cannot access other users' resources

## Technical Notes
- Use existing E2E test infrastructure
- Create test helper for registration/login to reduce boilerplate
- Clean up test users between tests (or use unique emails per test)
- Consider test isolation - each test should be independent
- Verify existing workout/enrollment/etc tests still pass

## Dependencies
- Blocks: None (this is the final validation)
- Blocked by: 001, 002, 003, 004, 005 (needs everything implemented)
