# Update API Documentation

## Overview

Update API documentation to include new state machine endpoints and response formats.

## Tasks

### 1. Update api-reference.md

Add documentation for new endpoints:

#### Workout Sessions
- `POST /workouts/start` - Start a new workout session
- `GET /workouts/{id}` - Get workout session by ID
- `POST /workouts/{id}/finish` - Complete a workout session
- `POST /workouts/{id}/abandon` - Abandon a workout session
- `GET /users/{userId}/workouts` - List user's workout history
- `GET /users/{userId}/workouts/current` - Get current in-progress workout

#### Enrollment State Management
- `POST /users/{userId}/enrollment/next-cycle` - Start next cycle
- `POST /users/{userId}/enrollment/advance-week` - Advance to next week

#### Updated Response Formats
- EnrollmentResponse now includes:
  - `enrollmentStatus`: ACTIVE | BETWEEN_CYCLES | QUIT
  - `cycleStatus`: PENDING | IN_PROGRESS | COMPLETED
  - `weekStatus`: PENDING | IN_PROGRESS | COMPLETED
  - `currentWorkoutSession`: Current in-progress session if any

### 2. Update workflows.md

Add new workflow section:

#### State Machine Workout Flow
1. User enrolled → enrollment_status: ACTIVE, cycle_status: PENDING
2. Start workout → POST /workouts/start
3. Log sets → existing logged_sets endpoints
4. Finish workout → POST /workouts/{id}/finish
5. After final day of week → POST /users/{userId}/enrollment/advance-week
6. At end of cycle → enrollment_status transitions to BETWEEN_CYCLES
7. Ready for next cycle → POST /users/{userId}/enrollment/next-cycle

### 3. Update example-requests.md

Add example requests for all new endpoints with realistic data.

### 4. Update openapi.yaml

Add OpenAPI specifications for:
- New endpoints
- New schemas (WorkoutSession, updated Enrollment)
- New error responses

## Acceptance Criteria

- [ ] api-reference.md updated with new endpoints
- [ ] workflows.md includes state machine workflow
- [ ] example-requests.md has examples for new endpoints
- [ ] openapi.yaml updated with new specifications
- [ ] Documentation is internally consistent
