# ERD 002: User Profile

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for implementing user profile viewing and updating in PowerPro. This builds on the authentication system to allow users to manage their profile data and preferences.

### Scope
- User profile retrieval (GET endpoint)
- User profile updates (PUT endpoint)
- Weight unit preference storage
- Authorization rules for profile access
- Does NOT include: email changes, password changes, profile pictures, account deletion

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Profile | User data including name, email, and preferences |
| Weight Unit | User's preferred unit for displaying weights (lb or kg) |
| Owner | The user whose profile is being accessed |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers (Frontend Teams): Integration and usage

## 2. Business / Stakeholder Needs

### Why
Users need to view their account information and customize their experience. The weight unit preference is critical for powerlifting programs where weights are calculated and displayed.

### Constraints
- SQLite only - no external services
- Email is immutable (changes require re-authentication - deferred)
- Build on existing authentication from Sprint 001

### Success Criteria
- User can retrieve their profile via GET endpoint
- User can update their profile via PUT endpoint
- Weight unit preference is stored and returned
- Unauthorized access returns 403

## 3. Functional Requirements

### Profile Entity Extension

#### REQ-PROFILE-001: Weight Unit Field
- **Description**: The system shall store a weight unit preference for each user
- **Rationale**: Powerlifters use different weight systems (lb/kg); users should choose their preferred unit
- **Priority**: Must
- **Acceptance Criteria**:
  - weight_unit column added to users table
  - Valid values: "lb" or "kg"
  - Default value: "lb"
  - NOT NULL constraint
- **Dependencies**: User table from Sprint 001

### Profile Retrieval

#### REQ-PROFILE-002: Get Profile Endpoint
- **Description**: The system shall provide an endpoint to retrieve a user's profile
- **Rationale**: Users need to view their account information
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /users/{id}/profile endpoint
  - Requires valid session token in Authorization header
  - Returns 200 with profile object (id, email, name, weight_unit, created_at, updated_at)
  - Returns 401 if no valid session
  - Returns 403 if user is not owner and not admin
  - Returns 404 if user does not exist
- **Dependencies**: REQ-PROFILE-001, Authentication from Sprint 001

#### REQ-PROFILE-003: Profile Response Format
- **Description**: The system shall return profile data in a consistent format
- **Rationale**: Frontend needs predictable data structure
- **Priority**: Must
- **Acceptance Criteria**:
  - Response includes: id, email, name (nullable), weightUnit, createdAt, updatedAt
  - Email is always present (required for authenticated users)
  - Name may be null
  - weightUnit is always present with default "lb"
  - Timestamps in ISO 8601 format
- **Dependencies**: REQ-PROFILE-002

### Profile Updates

#### REQ-PROFILE-004: Update Profile Endpoint
- **Description**: The system shall provide an endpoint to update a user's profile
- **Rationale**: Users need to modify their account information
- **Priority**: Must
- **Acceptance Criteria**:
  - PUT /users/{id}/profile endpoint
  - Requires valid session token in Authorization header
  - Request body: name (optional), weightUnit (optional)
  - Returns 200 with updated profile object
  - Returns 400 if invalid input (bad weight unit, name too long)
  - Returns 401 if no valid session
  - Returns 403 if user is not owner
  - Returns 404 if user does not exist
- **Dependencies**: REQ-PROFILE-001, REQ-PROFILE-002

#### REQ-PROFILE-005: Update Validation
- **Description**: The system shall validate profile update input
- **Rationale**: Prevent invalid data and provide helpful error messages
- **Priority**: Must
- **Acceptance Criteria**:
  - Name must be <= 100 characters if provided
  - weightUnit must be "lb" or "kg" if provided
  - Empty string for name clears the name (sets to null)
  - Validation errors return 400 with descriptive message
- **Dependencies**: REQ-PROFILE-004

#### REQ-PROFILE-006: Partial Updates
- **Description**: The system shall support partial profile updates
- **Rationale**: Users should be able to update only specific fields
- **Priority**: Should
- **Acceptance Criteria**:
  - Only provided fields are updated
  - Omitted fields retain their current values
  - Empty request body is valid (no-op, returns current profile)
- **Dependencies**: REQ-PROFILE-004

### Authorization

#### REQ-AUTH-006: Owner Access
- **Description**: The system shall allow users to access only their own profile
- **Rationale**: Users should not be able to view or modify other users' profiles
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /users/{id}/profile returns 403 if requester is not the owner (unless admin)
  - PUT /users/{id}/profile returns 403 if requester is not the owner
  - Admins cannot update other users' profiles (only view)
- **Dependencies**: REQ-PROFILE-002, REQ-PROFILE-004

#### REQ-AUTH-007: Admin Read Access
- **Description**: The system shall allow admins to view any user's profile
- **Rationale**: Admins need to view user data for support and moderation
- **Priority**: Should
- **Acceptance Criteria**:
  - GET /users/{id}/profile returns 200 for admin users regardless of owner
  - Admin status determined from session (is_admin field)
- **Dependencies**: REQ-AUTH-006, is_admin from Sprint 001

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Profile endpoints shall respond in < 100ms (p95)
- **NFR-002**: Profile update shall be atomic (no partial state)

### Security
- **NFR-003**: Profile endpoints require authentication
- **NFR-004**: Authorization checked before any data returned
- **NFR-005**: Error messages shall not leak information about other users

### Maintainability
- **NFR-006**: Profile logic shall be isolated in profile service/package
- **NFR-007**: Authorization logic shall be reusable for future endpoints

## 5. External Interfaces

### System Interfaces
- Database: SQLite with goose migrations
- API: REST over HTTP/HTTPS
- Authentication: Session tokens from Sprint 001

### Data Formats

#### Get Profile Response (200)
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "weightUnit": "lb",
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### Update Profile Request
```json
{
  "name": "John Doe",
  "weightUnit": "kg"
}
```

#### Update Profile Response (200)
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "weightUnit": "kg",
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-20T14:30:00Z"
}
```

#### Error Response (400/401/403/404)
```json
{
  "error": "Invalid weight unit. Must be 'lb' or 'kg'"
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- Go backend
- SQLite database
- goose migrations
- No external services

### Assumptions
- Authentication from Sprint 001 is implemented and working
- Users table has email, name, is_admin columns from Sprint 001
- Session middleware sets user context with user ID and admin status

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-PROFILE-001 | Unit tests, Migration test | weight_unit column added correctly |
| REQ-PROFILE-002, REQ-PROFILE-003 | Integration tests, E2E tests | GET profile returns correct data |
| REQ-PROFILE-004, REQ-PROFILE-005, REQ-PROFILE-006 | Integration tests, E2E tests | PUT profile updates correctly |
| REQ-AUTH-006, REQ-AUTH-007 | Integration tests, E2E tests | Authorization works correctly |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Weight unit column migration
- GET /users/{id}/profile endpoint
- PUT /users/{id}/profile endpoint
- Owner-only authorization
- Input validation

### Should-Have (This ERD)
- Admin read access
- Partial updates

### Won't-Have (Future ERDs)
- Email changes
- Password changes
- Profile pictures
- Account deletion

## 9. Traceability

### Links to PRD
- PRD-002: User Profile

### Links to Phase Document
- Phase 002: Frontend Readiness - Theme 2 (User Profile)

### Forward Links (to Tickets)
- 001-profile-schema-migration.md (REQ-PROFILE-001)
- 002-profile-domain-logic.md (REQ-PROFILE-002 through REQ-PROFILE-006)
- 003-profile-api-endpoints.md (REQ-PROFILE-002 through REQ-PROFILE-006)
- 004-profile-authorization.md (REQ-AUTH-006, REQ-AUTH-007)
- 005-profile-e2e-tests.md (All REQs)
