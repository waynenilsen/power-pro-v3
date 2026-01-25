# 003: Auth Service Domain Logic

## ERD Reference
Implements: REQ-AUTH-001, REQ-AUTH-002, REQ-AUTH-003, REQ-AUTH-004, REQ-AUTH-005

## Description
Implement the core authentication service with domain logic for registration, login, logout, and session management. This is the business logic layer independent of HTTP concerns.

## Context / Background
The auth service encapsulates all authentication logic: password hashing, session creation, session validation, etc. It should be usable from HTTP handlers and potentially other contexts (CLI, background jobs).

## Acceptance Criteria
- [ ] Create auth service package/module (e.g., `internal/auth/service.go`)
- [ ] Implement Register function:
  - Validates email format (contains @, valid domain part)
  - Validates password length (>= 8 chars)
  - Normalizes email to lowercase
  - Hashes password with bcrypt (cost 12)
  - Creates user with email, password_hash, optional name
  - Returns user object (no hash) or validation error
- [ ] Implement Login function:
  - Looks up user by email (case-insensitive)
  - Compares password with bcrypt
  - Creates new session with secure random token
  - Sets expiration to 7 days from now
  - Returns session token and user, or auth error
  - Same error for wrong email vs wrong password
- [ ] Implement Logout function:
  - Deletes session by token
  - Returns success even if session not found (idempotent)
- [ ] Implement ValidateSession function:
  - Looks up session by token
  - Verifies not expired
  - Returns user or error
- [ ] Implement GetUserBySession function:
  - Gets full user object from session token
  - Returns user or error
- [ ] Unit tests for all functions with >90% coverage
- [ ] Password and tokens never logged

## Technical Notes
- Use `golang.org/x/crypto/bcrypt` for password hashing
- Use `crypto/rand` for token generation, not math/rand
- Token: 32 bytes from crypto/rand, base64.URLEncoding.EncodeToString
- Consider extracting password hashing to separate module for testability
- Session expiration: time.Now().Add(7 * 24 * time.Hour)
- Email validation: regexp or net/mail package

## Dependencies
- Blocks: 005 (API endpoints call service)
- Blocked by: 001 (needs user schema), 002 (needs sessions schema)
