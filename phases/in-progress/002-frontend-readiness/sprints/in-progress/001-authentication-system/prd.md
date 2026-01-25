# PRD 001: Authentication System

## Product Vision

PowerPro requires real authentication to transition from a development API to a frontend-ready system. This PRD establishes session-based authentication that replaces the fake `X-User-ID` header system used during development.

## Strategic Objectives

1. **Enable Production Frontend**: Provide secure authentication that a real frontend can use
2. **Replace Test Headers**: Transition from fake `X-User-ID` headers to real session-based auth
3. **Simple Auth Flow**: Email/password authentication without external dependencies
4. **Secure by Design**: Proper password hashing and session management from the start

## Themes & Initiatives

### Theme 1: User Registration
- **Strategic Objective**: Enable Production Frontend
- **Rationale**: Users need to create accounts before they can use the API. Registration is the entry point.
- **Initiatives**:
  - Initiative A: Email/password registration endpoint
  - Initiative B: Password hashing with bcrypt or argon2
  - Initiative C: Optional name field during registration
  - Initiative D: Email uniqueness enforcement

### Theme 2: Session Management
- **Strategic Objective**: Replace Test Headers, Secure by Design
- **Rationale**: Sessions stored in SQLite provide server-side session management. Session tokens in headers enable stateless API calls.
- **Initiatives**:
  - Initiative A: Sessions table with user association and expiration
  - Initiative B: Secure random token generation
  - Initiative C: Session validation middleware
  - Initiative D: Token in Authorization header (Bearer scheme)

### Theme 3: Authentication Endpoints
- **Strategic Objective**: Enable Production Frontend
- **Rationale**: Standard auth endpoints (login, logout, current user) are required by every frontend.
- **Initiatives**:
  - Initiative A: Login endpoint returning session token
  - Initiative B: Logout endpoint invalidating session
  - Initiative C: Current user endpoint from session

### Theme 4: Middleware Transition
- **Strategic Objective**: Replace Test Headers
- **Rationale**: Existing endpoints use `X-User-ID` headers. New middleware must support both real auth and test headers for backwards compatibility.
- **Initiatives**:
  - Initiative A: Auth middleware that validates session tokens
  - Initiative B: Fallback to X-User-ID for E2E tests (configurable)
  - Initiative C: Admin role detection from session

## Success Metrics

| Metric | Target |
|--------|--------|
| User can register with email/password | Complete |
| User can login and receive session token | Complete |
| Session token works in Authorization header | Complete |
| Logout invalidates session | Complete |
| GET /auth/me returns current user | Complete |
| Existing E2E tests continue to pass | Complete |
| Password is never stored in plaintext | Complete |
| Session tokens are cryptographically secure | Complete |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Schema migrations (users, sessions) |
| Now | Auth service domain logic |
| Now | Auth endpoints |
| Now | Auth middleware with backwards compatibility |

## Dependencies

- None - this is the foundation of Phase 002

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking existing E2E tests | Medium | High | Maintain backwards compatibility with X-User-ID headers for testing |
| Session security vulnerabilities | Low | High | Use established patterns; crypto/rand for tokens; proper expiration |
| Password storage vulnerabilities | Low | Critical | Use bcrypt or argon2; never store plaintext |

## Out of Scope

- Email verification - deferred until email infrastructure exists
- Password reset flow - deferred until email infrastructure exists
- OAuth/social login - explicitly not planned (no external services)
- Magic links - explicitly not planned (no email infrastructure)
- JWT tokens - using server-side sessions instead
- Refresh tokens - sessions can be extended, not refreshed
- Multi-factor authentication - future enhancement
