# Auth Service Domain Logic

## Ticket Reference
`phases/in-progress/002-frontend-readiness/sprints/in-progress/001-authentication-system/tickets/todo/003-auth-service-domain-logic.md`

## Task
Implement the core authentication service with domain logic for registration, login, logout, and session management.

## Implementation

1. Create `internal/auth/` package
2. Implement service with functions:
   - **Register(email, password, name)**: Validate email/password, hash with bcrypt (cost 12), create user
   - **Login(email, password)**: Verify credentials, create session with secure random token (32 bytes, base64), 7-day expiry
   - **Logout(token)**: Delete session (idempotent)
   - **ValidateSession(token)**: Check session exists and not expired, return user
   - **GetUserBySession(token)**: Get full user object from session
3. Email validation: contains @, normalize to lowercase
4. Password validation: >= 8 chars
5. Same error for wrong email vs wrong password (security)
6. Use crypto/rand for tokens, golang.org/x/crypto/bcrypt for hashing

## Acceptance Criteria
- [ ] Auth service package created
- [ ] Register validates and creates users with hashed passwords
- [ ] Login verifies credentials and creates sessions
- [ ] Logout deletes sessions idempotently
- [ ] ValidateSession checks token and expiry
- [ ] Unit tests with >90% coverage
- [ ] Passwords and tokens never logged

## When Done
Move ticket from `tickets/todo/` to `tickets/done/` then run `crumbler delete`
