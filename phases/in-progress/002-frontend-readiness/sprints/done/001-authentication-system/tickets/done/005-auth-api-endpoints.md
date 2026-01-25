# 005: Auth API Endpoints

## ERD Reference
Implements: REQ-AUTH-001, REQ-AUTH-002, REQ-AUTH-003, REQ-AUTH-004, REQ-AUTH-005

## Description
Implement the HTTP handlers for authentication endpoints: register, login, logout, and current user.

## Context / Background
These endpoints provide the HTTP API layer on top of the auth service. They handle request parsing, response formatting, and HTTP status codes.

## Acceptance Criteria

### POST /auth/register
- [ ] Request body: `{"email": "...", "password": "...", "name": "..."}` (name optional)
- [ ] Returns 201 with user object (id, email, name, createdAt, updatedAt)
- [ ] Returns 400 if email missing or invalid format
- [ ] Returns 400 if password missing or < 8 characters
- [ ] Returns 409 if email already exists
- [ ] User object never includes password_hash

### POST /auth/login
- [ ] Request body: `{"email": "...", "password": "..."}`
- [ ] Returns 200 with `{"token": "...", "expiresAt": "...", "user": {...}}`
- [ ] Returns 400 if email or password missing
- [ ] Returns 401 if credentials invalid (same message for wrong email or password)
- [ ] Creates session in database

### POST /auth/logout
- [ ] Requires valid session (Authorization: Bearer <token>)
- [ ] Returns 204 No Content on success
- [ ] Returns 401 if no valid session
- [ ] Deletes session from database

### GET /auth/me
- [ ] Requires valid session (Authorization: Bearer <token>)
- [ ] Returns 200 with user object
- [ ] Returns 401 if no valid session

### General
- [ ] Register routes with router (e.g., under /auth prefix)
- [ ] Follow existing API response conventions
- [ ] Proper error response format: `{"error": "message"}`
- [ ] Integration tests for all endpoints

## Technical Notes
- Follow existing handler patterns in the codebase
- Use standard library or existing JSON handling patterns
- Consistent timestamp format (ISO8601/RFC3339)
- Consider rate limiting for register/login (future enhancement, not this ticket)

## Dependencies
- Blocks: 006 (E2E tests test these endpoints)
- Blocked by: 003 (needs auth service), 004 (needs middleware for protected endpoints)
