# Sprint 001 - Authentication System

## Overview

Implement real session-based authentication to replace the fake `X-User-ID` and `X-Admin` headers used for testing.

## Technical Requirements

- **Session auth, NOT JWT** - Server-side session storage in SQLite
- **Headless API** - Auth token in header (`Authorization: Bearer <session_token>`)
- **Email/password signup** - No OAuth, no magic links
- **Name is optional** - Only email + password required for registration
- **No email verification** - Defer until email infrastructure exists
- **Password hashing** - bcrypt or argon2

## Endpoints Needed

```
POST /auth/register    - Create account (email, password, optional name)
POST /auth/login       - Get session token
POST /auth/logout      - Invalidate session
GET  /auth/me          - Get current user from session
```

## Schema Changes

- Add `email`, `password_hash`, `name` to `users` table
- Create `sessions` table (id, user_id, token, expires_at, created_at)

## Tasks

1. Create sprint directory: `phases/in-progress/002-frontend-readiness/sprints/in-progress/001-authentication-system/`
2. Write `prd.md` - Product requirements for authentication
3. Write `erd.md` - Engineering requirements with detailed specs
4. Create ticket directory structure
5. Create tickets (at least 5):
   - Schema migration for users table (add email, password_hash, name)
   - Schema migration for sessions table
   - Auth service domain logic (register, login, logout, session validation)
   - Auth middleware (replace X-User-ID header handling)
   - Auth API endpoints
   - E2E tests for auth flow

## Acceptance Criteria

- [ ] Sprint directory structure exists
- [ ] PRD document completed
- [ ] ERD document completed with REQ-XXX requirements
- [ ] At least 5 tickets created in `tickets/todo/`
- [ ] Tickets reference ERD requirements
