# ERD 003: Dashboard Endpoint

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for implementing an aggregated dashboard endpoint that combines multiple data sources into a single API response. This is an aggregation layer - no new data is stored, only existing data is queried and combined.

### Scope
- Dashboard aggregation service
- Single GET endpoint for user dashboard
- Aggregation of: enrollment, next workout, current session, recent workouts, current maxes
- Does NOT include: caching, real-time updates, new data storage

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Dashboard | Aggregated view of user's current program state |
| Enrollment | User's active participation in a training program |
| Training Max | User's current calculated max for a specific lift |
| Workout Session | An individual workout instance (may be in-progress or completed) |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers (Frontend Teams): Integration and usage

## 2. Business / Stakeholder Needs

### Why
Frontend clients need multiple pieces of information to render a home screen. Currently, this requires multiple API calls. An aggregated endpoint reduces latency, simplifies frontend logic, and provides a consistent snapshot of user state.

### Constraints
- SQLite only - all queries must be SQLite compatible
- Must use existing domain models and services
- No new database tables or schema changes
- Authentication required via existing session system

### Success Criteria
- Single endpoint returns all dashboard sections
- Response structure is consistent regardless of data availability
- Empty data returns empty arrays/null, not errors
- Response time < 200ms for p95

## 3. Functional Requirements

### Dashboard Aggregation

#### REQ-DASH-001: Dashboard Endpoint
- **Description**: The system shall provide an endpoint to retrieve aggregated dashboard data for a user
- **Rationale**: Reduces multiple API calls to a single request
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /users/{id}/dashboard endpoint
  - Requires valid session token in Authorization header
  - Returns 200 with dashboard object
  - Returns 401 if no valid session
  - Returns 403 if user is not owner
  - Returns 404 if user does not exist
- **Dependencies**: Authentication from Sprint 001

#### REQ-DASH-002: Dashboard Response Structure
- **Description**: The system shall return dashboard data in a consistent structure with five sections
- **Rationale**: Frontend needs predictable data structure for all possible states
- **Priority**: Must
- **Acceptance Criteria**:
  - Response includes sections: enrollment, nextWorkout, currentSession, recentWorkouts, currentMaxes
  - All sections are always present (may be null or empty array)
  - Section order is consistent
  - Response uses camelCase for JSON fields
- **Dependencies**: REQ-DASH-001

### Enrollment Section

#### REQ-DASH-003: Enrollment Status
- **Description**: The system shall aggregate current enrollment status and program position
- **Rationale**: Users need to see their current program context at a glance
- **Priority**: Must
- **Acceptance Criteria**:
  - Returns null if user has no active enrollment
  - Returns object with: status, programName, cycleIteration, cycleStatus, weekNumber, weekStatus
  - status values: "ACTIVE", "PAUSED"
  - cycleStatus/weekStatus values: "IN_PROGRESS", "COMPLETED"
  - cycleIteration is 1-indexed (first cycle = 1)
  - weekNumber is 1-indexed within the cycle
- **Dependencies**: REQ-DASH-002, Enrollment domain model

### Next Workout Section

#### REQ-DASH-004: Next Workout Preview
- **Description**: The system shall calculate and return the next scheduled workout
- **Rationale**: Users want to see upcoming workout without navigating
- **Priority**: Must
- **Acceptance Criteria**:
  - Returns null if no enrollment or program completed
  - Returns object with: dayName, daySlug, exerciseCount, estimatedSets
  - dayName is human-readable (e.g., "Volume Day")
  - daySlug is URL-friendly (e.g., "volume")
  - exerciseCount is number of distinct exercises
  - estimatedSets is total sets across all exercises
- **Dependencies**: REQ-DASH-002, REQ-DASH-003, Program template data

### Current Session Section

#### REQ-DASH-005: Current Session Detection
- **Description**: The system shall detect and return any in-progress workout session
- **Rationale**: Users mid-workout should immediately see their active session
- **Priority**: Must
- **Acceptance Criteria**:
  - Returns null if no active session
  - Returns object with: sessionId, dayName, startedAt, setsCompleted, totalSets
  - startedAt in ISO 8601 format
  - setsCompleted and totalSets are integers
  - Only one active session per user at a time
- **Dependencies**: REQ-DASH-002, WorkoutSession domain model

### Recent Workouts Section

#### REQ-DASH-006: Recent Workout History
- **Description**: The system shall return a list of recently completed workouts
- **Rationale**: Users want to see recent activity and verify logging
- **Priority**: Must
- **Acceptance Criteria**:
  - Returns empty array if no completed workouts
  - Returns array of objects with: date, dayName, setsCompleted
  - date in ISO 8601 date format (YYYY-MM-DD)
  - Maximum 5 recent workouts
  - Ordered by date descending (most recent first)
- **Dependencies**: REQ-DASH-002, WorkoutSession domain model

### Current Maxes Section

#### REQ-DASH-007: Current Training Maxes
- **Description**: The system shall return current training maxes for all tracked lifts
- **Rationale**: Powerlifters constantly reference maxes for workout planning
- **Priority**: Must
- **Acceptance Criteria**:
  - Returns empty array if no maxes set
  - Returns array of objects with: lift, value, type
  - lift is the exercise name (e.g., "Squat", "Bench Press")
  - value is the weight (in user's preferred unit from profile)
  - type is max type: "TRAINING_MAX", "ONE_REP_MAX", "ESTIMATED"
  - Order by lift name alphabetically
- **Dependencies**: REQ-DASH-002, TrainingMax domain model, User weight_unit from Sprint 002

### Authorization

#### REQ-DASH-008: Owner-Only Access
- **Description**: The system shall restrict dashboard access to the owning user only
- **Rationale**: Dashboard contains personal training data that should not be exposed to other users
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /users/{id}/dashboard returns 403 if requester is not the owner
  - Admin users cannot access other users' dashboards
  - Error message does not leak information about other users
- **Dependencies**: REQ-DASH-001, Authentication from Sprint 001

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Dashboard endpoint shall respond in < 200ms (p95)
- **NFR-002**: Queries shall use efficient joins to minimize database round-trips
- **NFR-003**: Dashboard shall execute as a single transaction for consistency

### Reliability
- **NFR-004**: Partial failures in one section shall not fail the entire response
- **NFR-005**: Missing data shall return null/empty, not errors

### Maintainability
- **NFR-006**: Dashboard aggregation logic shall be isolated in its own service
- **NFR-007**: Each section shall have its own aggregation function for testability
- **NFR-008**: Dashboard service shall depend on existing domain services, not repositories directly

## 5. External Interfaces

### System Interfaces
- Database: SQLite (queries existing tables)
- API: REST over HTTP/HTTPS
- Authentication: Session tokens from Sprint 001

### Data Formats

#### Dashboard Response (200)
```json
{
  "enrollment": {
    "status": "ACTIVE",
    "programName": "Texas Method",
    "cycleIteration": 2,
    "cycleStatus": "IN_PROGRESS",
    "weekNumber": 1,
    "weekStatus": "IN_PROGRESS"
  },
  "nextWorkout": {
    "dayName": "Volume Day",
    "daySlug": "volume",
    "exerciseCount": 2,
    "estimatedSets": 10
  },
  "currentSession": null,
  "recentWorkouts": [
    {"date": "2024-01-15", "dayName": "Intensity Day", "setsCompleted": 4},
    {"date": "2024-01-12", "dayName": "Recovery Day", "setsCompleted": 6}
  ],
  "currentMaxes": [
    {"lift": "Bench Press", "value": 227.5, "type": "TRAINING_MAX"},
    {"lift": "Squat", "value": 320, "type": "TRAINING_MAX"}
  ]
}
```

#### Dashboard with No Enrollment (200)
```json
{
  "enrollment": null,
  "nextWorkout": null,
  "currentSession": null,
  "recentWorkouts": [],
  "currentMaxes": []
}
```

#### Dashboard with Active Session (200)
```json
{
  "enrollment": {
    "status": "ACTIVE",
    "programName": "Texas Method",
    "cycleIteration": 1,
    "cycleStatus": "IN_PROGRESS",
    "weekNumber": 1,
    "weekStatus": "IN_PROGRESS"
  },
  "nextWorkout": null,
  "currentSession": {
    "sessionId": "uuid",
    "dayName": "Volume Day",
    "startedAt": "2024-01-20T14:30:00Z",
    "setsCompleted": 3,
    "totalSets": 10
  },
  "recentWorkouts": [],
  "currentMaxes": [
    {"lift": "Squat", "value": 320, "type": "TRAINING_MAX"}
  ]
}
```

#### Error Response (401/403/404)
```json
{
  "error": "Unauthorized"
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- Go backend
- SQLite database
- No new tables or schema changes
- Must use existing domain services

### Assumptions
- Authentication from Sprint 001 is implemented and working
- User profile with weight_unit from Sprint 002 is available
- Enrollment, WorkoutSession, and TrainingMax domain models exist
- Program templates with day/exercise definitions are queryable

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-DASH-001, REQ-DASH-002 | Integration tests, E2E tests | Endpoint returns correct structure |
| REQ-DASH-003 | Unit tests, E2E tests | Enrollment section populates correctly |
| REQ-DASH-004 | Unit tests, E2E tests | Next workout calculated correctly |
| REQ-DASH-005 | Unit tests, E2E tests | Active session detected correctly |
| REQ-DASH-006 | Unit tests, E2E tests | Recent workouts returned correctly |
| REQ-DASH-007 | Unit tests, E2E tests | Current maxes returned correctly |
| REQ-DASH-008 | Integration tests, E2E tests | Authorization enforced correctly |
| NFR-001 | Performance tests | < 200ms p95 response time |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Dashboard endpoint with all five sections
- Owner-only authorization
- Consistent response structure

### Should-Have (This ERD)
- Efficient query optimization
- Unit weight conversion based on user preference

### Won't-Have (Future ERDs)
- Caching layer
- Real-time updates
- Analytics aggregation
- Customizable sections

## 9. Traceability

### Links to PRD
- PRD-003: Dashboard Endpoint

### Links to Phase Document
- Phase 002: Frontend Readiness - Theme 3 (Dashboard Aggregation)

### Forward Links (to Tickets)
- 001-dashboard-service.md (REQ-DASH-001, REQ-DASH-002, NFR-006, NFR-007, NFR-008)
- 002-enrollment-aggregation.md (REQ-DASH-003)
- 003-next-workout-calculation.md (REQ-DASH-004)
- 004-current-session-query.md (REQ-DASH-005)
- 005-recent-workouts-query.md (REQ-DASH-006)
- 006-current-maxes-query.md (REQ-DASH-007)
- 007-dashboard-api-endpoint.md (REQ-DASH-001, REQ-DASH-008)
- 008-dashboard-e2e-tests.md (All REQs)
