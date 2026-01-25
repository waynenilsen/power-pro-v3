# 004: Auth Middleware

## ERD Reference
Implements: REQ-MIDDLEWARE-001, REQ-MIDDLEWARE-002, REQ-MIDDLEWARE-003

## Description
Create HTTP middleware that authenticates requests using session tokens from the Authorization header, with fallback to X-User-ID headers for E2E tests.

## Context / Background
The current system uses X-User-ID and X-Admin headers for all authentication, which is only suitable for testing. The new middleware should:
1. First try Authorization: Bearer <token> (real auth)
2. Fall back to X-User-ID (test auth) if no Authorization header and in test mode
3. Set user context for downstream handlers

This allows gradual migration and keeps existing E2E tests working.

## Acceptance Criteria
- [ ] Create auth middleware that extracts token from `Authorization: Bearer <token>` header
- [ ] Middleware validates session token using auth service
- [ ] On valid session: set user in request context, continue
- [ ] On invalid/expired session: return 401 Unauthorized
- [ ] Implement test mode fallback:
  - If no Authorization header AND test mode enabled
  - Check for X-User-ID header
  - If present, look up user by ID (or create minimal user context)
  - X-Admin header sets admin flag in context
- [ ] Test mode controlled by environment variable (e.g., `POWERPRO_TEST_MODE=true`)
- [ ] Create helper functions to extract user from context
- [ ] Create optional middleware for admin-only endpoints
- [ ] Unit tests for middleware with mocked auth service
- [ ] Integration tests verifying both auth paths

## Technical Notes
- Use standard Go context for user propagation
- Define context keys as unexported type: `type contextKey string`
- Consider middleware composition with existing middleware
- Test mode detection: `os.Getenv("POWERPRO_TEST_MODE") == "true"`
- Authorization header parsing: `strings.TrimPrefix(header, "Bearer ")`
- Return consistent 401 response format matching API conventions

## Dependencies
- Blocks: 005 (endpoints use middleware), 006 (tests verify middleware)
- Blocked by: 003 (needs auth service for validation)
