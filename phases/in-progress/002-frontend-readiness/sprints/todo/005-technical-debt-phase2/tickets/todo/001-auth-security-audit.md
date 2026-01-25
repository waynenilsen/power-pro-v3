# 001: Authentication Security Audit

## ERD Reference
Implements: REQ-TD2-001, REQ-TD2-002, REQ-TD2-003

## Description
Conduct a security audit of the authentication implementation to ensure it follows security best practices. This includes reviewing password hashing, session token generation, and authorization middleware.

## Context / Background
Authentication is security-critical code. This audit ensures the Phase 2 auth implementation doesn't contain vulnerabilities that could compromise user accounts or system security.

## Acceptance Criteria

### Password Hashing (REQ-TD2-001)
- [ ] bcrypt or argon2 used with appropriate cost factor (bcrypt >= 10, argon2 default parameters)
- [ ] No plaintext passwords logged anywhere in the codebase
- [ ] Password comparison uses timing-safe comparison
- [ ] Password validation (min length, etc.) applied before hashing

### Session Token Security (REQ-TD2-002)
- [ ] Tokens generated using crypto/rand (not math/rand)
- [ ] Token length provides >= 256 bits of entropy
- [ ] Token validation uses timing-safe comparison
- [ ] Tokens stored as hashes in database, not plaintext

### Authorization Middleware (REQ-TD2-003)
- [ ] All protected endpoints require valid session token
- [ ] User context properly extracted from session and attached to request
- [ ] Expired sessions are rejected with 401 Unauthorized
- [ ] Invalid/malformed tokens rejected with appropriate error
- [ ] Session validation happens before route handlers

## Technical Notes
- Use `go vet` and `staticcheck` for static analysis
- Review OWASP authentication guidelines for reference
- Check for timing attacks in comparison functions
- Verify no sensitive data in logs or error messages

## Dependencies
- Blocks: None
- Blocked by: Sprint 001 (auth system) complete
- Related: 002-session-expiration-tests

## Resources / Links
- ERD: phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/erd.md
- OWASP Auth Cheatsheet: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
