# Auth Middleware

## Ticket Reference
`phases/in-progress/002-frontend-readiness/sprints/in-progress/001-authentication-system/tickets/todo/004-auth-middleware.md`

## Task
Create HTTP middleware that authenticates requests using session tokens, with fallback to X-User-ID for tests.

## Implementation

1. Create middleware that:
   - Extracts token from `Authorization: Bearer <token>`
   - Validates session via auth service
   - Sets user in request context
   - Returns 401 for invalid/expired sessions
2. Test mode fallback (when `POWERPRO_TEST_MODE=true`):
   - If no Authorization header, check X-User-ID header
   - X-Admin header sets admin flag in context
3. Context helpers:
   - UserFromContext(ctx) - extract user
   - IsAdmin(ctx) - check admin flag
4. Optional admin-only middleware for protected endpoints

## Acceptance Criteria
- [ ] Middleware extracts and validates Bearer tokens
- [ ] Valid session sets user in context
- [ ] Invalid/expired session returns 401
- [ ] Test mode falls back to X-User-ID/X-Admin headers
- [ ] Context helper functions created
- [ ] Admin-only middleware created
- [ ] Unit tests with mocked auth service

## When Done
Move ticket from `tickets/todo/` to `tickets/done/` then run `crumbler delete`
