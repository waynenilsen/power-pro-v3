# Auth API Endpoints

## Ticket Reference
`phases/in-progress/002-frontend-readiness/sprints/in-progress/001-authentication-system/tickets/todo/005-auth-api-endpoints.md`

## Task
Implement HTTP handlers for auth endpoints: register, login, logout, current user.

## Endpoints

### POST /auth/register
- Body: `{"email": "...", "password": "...", "name": "..."}`
- 201: User created (id, email, name, createdAt, updatedAt)
- 400: Invalid email/password
- 409: Email exists

### POST /auth/login
- Body: `{"email": "...", "password": "..."}`
- 200: `{"token": "...", "expiresAt": "...", "user": {...}}`
- 400: Missing credentials
- 401: Invalid credentials

### POST /auth/logout
- Requires valid session
- 204: Success
- 401: No valid session

### GET /auth/me
- Requires valid session
- 200: User object
- 401: No valid session

## Acceptance Criteria
- [ ] All 4 endpoints implemented
- [ ] Routes registered under /auth prefix
- [ ] Follow existing response conventions
- [ ] Password_hash never in responses
- [ ] Integration tests for all endpoints

## When Done
Move ticket from `tickets/todo/` to `tickets/done/` then run `crumbler delete`
