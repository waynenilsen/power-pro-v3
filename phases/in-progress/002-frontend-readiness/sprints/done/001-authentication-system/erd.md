# ERD 001: Authentication System

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for implementing session-based authentication in PowerPro. This replaces the development-only `X-User-ID` header approach with real authentication suitable for production frontends.

### Scope
- User registration with email/password
- Session-based authentication with SQLite storage
- Auth endpoints (register, login, logout, current user)
- Auth middleware with backwards compatibility for E2E tests
- Does NOT include: email verification, password reset, OAuth, MFA

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Session | Server-side record linking a token to a user with expiration |
| Session Token | Cryptographically secure random string used to authenticate requests |
| Password Hash | One-way hash of user password using bcrypt or argon2 |
| Bearer Token | Authentication scheme where token is sent in Authorization header |
| X-User-ID | Development-only header for testing, to be deprecated |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers (Frontend Teams): Integration and usage

## 2. Business / Stakeholder Needs

### Why
Every protected endpoint requires authentication. The current `X-User-ID` header approach is suitable only for development. A real frontend cannot be built without proper authentication.

### Constraints
- SQLite only - no external session stores (Redis, etc.)
- No external auth services (Auth0, etc.)
- No email infrastructure yet (no verification, no password reset)
- Must maintain backwards compatibility with existing E2E tests

### Success Criteria
- User can create account with email + password
- User can authenticate and receive session token
- Session token works in Authorization header for all protected endpoints
- Existing E2E tests continue to pass with X-User-ID header (test mode)
- Password is never stored or logged in plaintext

## 3. Functional Requirements

### User Entity Extension

#### REQ-USER-001: Email Field
- **Description**: The system shall store an email address for each user
- **Rationale**: Email is the primary identifier for authentication
- **Priority**: Must
- **Acceptance Criteria**:
  - Email column added to users table
  - Email is unique across all users
  - Email is case-insensitive (stored lowercase)
  - Email format is validated (contains @ and .)
  - Max length 255 characters
- **Dependencies**: None

#### REQ-USER-002: Password Hash Field
- **Description**: The system shall store a hashed password for each user
- **Rationale**: Password hashing is required for secure authentication
- **Priority**: Must
- **Acceptance Criteria**:
  - password_hash column added to users table
  - Hash generated using bcrypt (cost 12) or argon2id
  - Plaintext password never stored or logged
  - Hash is nullable for existing test users
- **Dependencies**: REQ-USER-001

#### REQ-USER-003: Name Field
- **Description**: The system shall store an optional display name for each user
- **Rationale**: Users may want to personalize their experience with a name
- **Priority**: Should
- **Acceptance Criteria**:
  - name column added to users table
  - Name is optional (nullable)
  - Max length 100 characters
- **Dependencies**: REQ-USER-001

### Session Entity

#### REQ-SESSION-001: Session Identification
- **Description**: The system shall provide unique identification for each session
- **Rationale**: Enables session lookup and management
- **Priority**: Must
- **Acceptance Criteria**: Each session has a unique UUID identifier
- **Dependencies**: None

#### REQ-SESSION-002: User Association
- **Description**: The system shall associate each session with a specific user
- **Rationale**: Sessions authenticate users; must link to user
- **Priority**: Must
- **Acceptance Criteria**:
  - user_id foreign key to users table
  - Required field, not nullable
  - Cascade delete when user is deleted
- **Dependencies**: REQ-SESSION-001

#### REQ-SESSION-003: Session Token
- **Description**: The system shall store a cryptographically secure token for each session
- **Rationale**: Token is used for authentication in API requests
- **Priority**: Must
- **Acceptance Criteria**:
  - token column with unique constraint
  - Token generated using crypto/rand (32 bytes, base64 encoded)
  - Token indexed for efficient lookup
- **Dependencies**: REQ-SESSION-001

#### REQ-SESSION-004: Session Expiration
- **Description**: The system shall track session expiration
- **Rationale**: Sessions should expire for security
- **Priority**: Must
- **Acceptance Criteria**:
  - expires_at timestamp column
  - Default expiration: 7 days from creation
  - Expired sessions are not valid for authentication
- **Dependencies**: REQ-SESSION-001

#### REQ-SESSION-005: Session Timestamps
- **Description**: The system shall track session creation time
- **Rationale**: Audit trail and debugging
- **Priority**: Must
- **Acceptance Criteria**:
  - created_at timestamp column, required
- **Dependencies**: REQ-SESSION-001

### Registration

#### REQ-AUTH-001: User Registration
- **Description**: The system shall allow users to register with email and password
- **Rationale**: Users need to create accounts to use the system
- **Priority**: Must
- **Acceptance Criteria**:
  - POST /auth/register endpoint
  - Request body: email (required), password (required), name (optional)
  - Password minimum 8 characters
  - Returns 201 with user object (no password hash)
  - Returns 400 if email or password missing
  - Returns 400 if password too short
  - Returns 409 if email already exists
- **Dependencies**: REQ-USER-001, REQ-USER-002, REQ-USER-003

#### REQ-AUTH-002: Registration Validation
- **Description**: The system shall validate registration input
- **Rationale**: Prevent invalid data and provide helpful error messages
- **Priority**: Must
- **Acceptance Criteria**:
  - Email must be valid format (contains @ and valid domain part)
  - Password must be at least 8 characters
  - Name must be <= 100 characters if provided
  - Validation errors return 400 with descriptive message
- **Dependencies**: REQ-AUTH-001

### Login/Logout

#### REQ-AUTH-003: User Login
- **Description**: The system shall authenticate users and return session tokens
- **Rationale**: Users need to authenticate to access protected resources
- **Priority**: Must
- **Acceptance Criteria**:
  - POST /auth/login endpoint
  - Request body: email (required), password (required)
  - On success: Creates session, returns 200 with token and user object
  - Returns 400 if email or password missing
  - Returns 401 if email not found or password incorrect
  - Same error message for wrong email vs wrong password (security)
- **Dependencies**: REQ-USER-001, REQ-USER-002, REQ-SESSION-001 through REQ-SESSION-005

#### REQ-AUTH-004: User Logout
- **Description**: The system shall invalidate sessions on logout
- **Rationale**: Users need to be able to end their sessions
- **Priority**: Must
- **Acceptance Criteria**:
  - POST /auth/logout endpoint
  - Requires valid session token in Authorization header
  - Deletes the session from database
  - Returns 204 No Content on success
  - Returns 401 if no valid session
- **Dependencies**: REQ-SESSION-001 through REQ-SESSION-005

#### REQ-AUTH-005: Current User
- **Description**: The system shall provide the current user from session
- **Rationale**: Frontends need to know who is logged in
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /auth/me endpoint
  - Requires valid session token in Authorization header
  - Returns 200 with user object (id, email, name, created_at, updated_at)
  - Returns 401 if no valid session
- **Dependencies**: REQ-SESSION-001 through REQ-SESSION-005

### Authentication Middleware

#### REQ-MIDDLEWARE-001: Session Token Validation
- **Description**: The system shall validate session tokens in Authorization header
- **Rationale**: Every protected endpoint needs authentication
- **Priority**: Must
- **Acceptance Criteria**:
  - Middleware extracts token from Authorization: Bearer <token> header
  - Looks up session by token
  - Verifies session not expired
  - Sets user context for downstream handlers
  - Returns 401 if token missing, invalid, or expired
- **Dependencies**: REQ-SESSION-001 through REQ-SESSION-005

#### REQ-MIDDLEWARE-002: Test Header Fallback
- **Description**: The system shall support X-User-ID header for E2E tests
- **Rationale**: Existing E2E tests use X-User-ID; need backwards compatibility
- **Priority**: Must
- **Acceptance Criteria**:
  - If Authorization header present, use session auth (primary)
  - If no Authorization header but X-User-ID present, use test auth (fallback)
  - Test auth only active in development/test mode (not production)
  - X-Admin header continues to work for admin testing
- **Dependencies**: REQ-MIDDLEWARE-001

#### REQ-MIDDLEWARE-003: Admin Detection
- **Description**: The system shall detect admin users from session
- **Rationale**: Some endpoints require admin privileges
- **Priority**: Should
- **Acceptance Criteria**:
  - Add is_admin boolean to users table (default false)
  - Admin status set in user context from session
  - X-Admin header fallback for test mode
- **Dependencies**: REQ-MIDDLEWARE-001, REQ-MIDDLEWARE-002

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Auth endpoints shall respond in < 200ms (p95)
- **NFR-002**: Session lookup shall respond in < 50ms (p95)
- **NFR-003**: Password hashing shall use appropriate cost (bcrypt cost 12 or equivalent)

### Security
- **NFR-004**: Passwords shall be hashed with bcrypt (cost 12) or argon2id
- **NFR-005**: Session tokens shall be 32 bytes from crypto/rand, base64 encoded
- **NFR-006**: Session tokens shall be indexed for efficient lookup
- **NFR-007**: Authentication errors shall not reveal whether email exists
- **NFR-008**: Password and tokens shall never appear in logs

### Reliability
- **NFR-009**: Session operations shall be atomic (no partial state)
- **NFR-010**: Expired sessions shall be rejected immediately

### Maintainability
- **NFR-011**: Auth logic shall be isolated in auth service/package
- **NFR-012**: Middleware shall be composable with existing middleware

## 5. External Interfaces

### System Interfaces
- Database: SQLite with goose migrations
- API: REST over HTTP/HTTPS

### Data Formats

#### Registration Request
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"  // optional
}
```

#### Registration Response (201)
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### Login Request
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Login Response (200)
```json
{
  "token": "base64-encoded-session-token",
  "expiresAt": "2024-01-22T10:00:00Z",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

#### Current User Response (200)
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### Error Response (400/401/409)
```json
{
  "error": "Invalid email or password"
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- Go backend
- SQLite database
- goose migrations
- No external services

### Assumptions
- Test mode can be detected via environment variable
- Existing X-User-ID header tests will continue to work
- bcrypt package available for Go

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-USER-001 to REQ-USER-003 | Unit tests, Migration test | User schema extended correctly |
| REQ-SESSION-001 to REQ-SESSION-005 | Unit tests, Migration test | Sessions table created correctly |
| REQ-AUTH-001 to REQ-AUTH-005 | Integration tests, E2E tests | All auth endpoints work correctly |
| REQ-MIDDLEWARE-001 to REQ-MIDDLEWARE-003 | Integration tests, E2E tests | Auth middleware validates correctly |
| NFR-004, NFR-005 | Security review, Unit tests | Proper cryptographic functions used |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- User schema extension (email, password_hash, name)
- Sessions table
- Registration endpoint
- Login endpoint
- Logout endpoint
- Current user endpoint
- Session validation middleware
- X-User-ID fallback for tests

### Should-Have (This ERD)
- Admin flag on users
- Admin detection in middleware

### Won't-Have (Future ERDs)
- Email verification
- Password reset
- OAuth/social login
- MFA
- Session refresh

## 9. Traceability

### Links to PRD
- PRD-001: Authentication System

### Links to Phase Document
- Phase 002: Frontend Readiness - Theme 1 (Authentication System)

### Forward Links (to Tickets)
- 001-user-schema-migration.md (REQ-USER-001, REQ-USER-002, REQ-USER-003)
- 002-sessions-schema-migration.md (REQ-SESSION-001 through REQ-SESSION-005)
- 003-auth-service-domain-logic.md (REQ-AUTH-001 through REQ-AUTH-005)
- 004-auth-middleware.md (REQ-MIDDLEWARE-001 through REQ-MIDDLEWARE-003)
- 005-auth-api-endpoints.md (REQ-AUTH-001 through REQ-AUTH-005)
- 006-auth-e2e-tests.md (All REQs)
