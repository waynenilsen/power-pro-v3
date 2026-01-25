# Workout Session Handlers

## Overview

Create HTTP handlers for workout session management endpoints.

## Endpoints to Implement

```
POST   /workouts/start                    # Start new workout, returns session_id
GET    /workouts/{id}                     # Get workout session details
POST   /workouts/{id}/finish              # Mark workout complete
POST   /workouts/{id}/abandon             # Mark workout abandoned
GET    /users/{id}/workouts               # List user's workout history
GET    /users/{id}/workouts/current       # Get current in-progress workout
```

## Implementation Details

### POST /workouts/start
- Requires active enrollment
- Creates IN_PROGRESS session with current week/day from UserProgramState
- Fails if user already has IN_PROGRESS session (return clear error)
- Emits `WORKOUT_STARTED` event
- Returns: session details with ID

### GET /workouts/{id}
- Returns full session details
- Include: status, week number, day index, started/finished timestamps
- Authorization: owner or admin

### POST /workouts/{id}/finish
- Validates session is IN_PROGRESS
- Marks session COMPLETED
- Triggers state advancement (week/cycle completion checks)
- Emits `WORKOUT_COMPLETED` event
- Returns: updated session details

### POST /workouts/{id}/abandon
- Validates session is IN_PROGRESS
- Marks session ABANDONED
- Emits `WORKOUT_ABANDONED` event
- Does NOT trigger state advancement
- Returns: updated session details

### GET /users/{id}/workouts
- List all workout sessions for user
- Support pagination (limit/offset)
- Optional status filter query param
- Returns: list with meta

### GET /users/{id}/workouts/current
- Get active IN_PROGRESS session if exists
- Returns 404 if no active session
- Returns: session details

## Files to Create/Modify

- `internal/api/workout_session_handler.go` - New handler file
- `internal/api/routes.go` - Register new routes
- `internal/api/workout_session_handler_test.go` - Unit tests

## Acceptance Criteria

- [ ] All 6 endpoints implemented
- [ ] Cannot start workout if one already in progress
- [ ] Events emitted on start/finish/abandon
- [ ] State advancement triggered on finish
- [ ] Unit tests pass
