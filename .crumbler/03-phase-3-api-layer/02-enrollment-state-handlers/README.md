# Enrollment State Handlers

## Overview

Add/update enrollment endpoints to expose state machine fields and support cycle/week management.

## Endpoints to Implement

### New Endpoints

```
POST   /users/{id}/enrollment/next-cycle    # Start new cycle (when BETWEEN_CYCLES)
POST   /users/{id}/enrollment/advance-week  # Manual week advance
```

### Update Existing

```
GET    /users/{id}/enrollment               # Include all state fields
```

## Implementation Details

### GET /users/{id}/enrollment (Update)
Current response needs to include:
- `enrollment_status`: ACTIVE | BETWEEN_CYCLES | QUIT
- `cycle_status`: PENDING | IN_PROGRESS | COMPLETED
- `week_status`: PENDING | IN_PROGRESS | COMPLETED
- `current_workout_session`: {...} or null (if one is in progress)

### POST /users/{id}/enrollment/next-cycle
- Validates enrollment exists and is BETWEEN_CYCLES
- Transitions enrollment to ACTIVE
- Increments cycle iteration
- Resets week to 1
- Resets week_status to PENDING, cycle_status to PENDING
- Emits `CYCLE_STARTED` event
- Returns: updated enrollment

### POST /users/{id}/enrollment/advance-week
- Validates enrollment is ACTIVE
- Advances to next week
- Handles cycle boundary (sets BETWEEN_CYCLES if final week)
- Emits `WEEK_COMPLETED` event
- Emits `CYCLE_BOUNDARY_REACHED` if applicable
- Returns: updated enrollment

## Files to Modify

- `internal/api/enrollment_handler.go` - Update existing handler
- `internal/api/enrollment_handler_test.go` - Update tests

## Response Shape

```json
{
  "id": "...",
  "userId": "...",
  "program": {...},
  "state": {
    "currentWeek": 1,
    "currentCycleIteration": 1,
    "currentDayIndex": null
  },
  "enrollmentStatus": "ACTIVE",
  "cycleStatus": "IN_PROGRESS",
  "weekStatus": "IN_PROGRESS",
  "currentWorkoutSession": null,
  "enrolledAt": "...",
  "updatedAt": "..."
}
```

## Acceptance Criteria

- [ ] GET enrollment includes all state fields
- [ ] GET enrollment includes current workout session if active
- [ ] next-cycle endpoint works from BETWEEN_CYCLES state
- [ ] advance-week endpoint advances correctly
- [ ] Events emitted on transitions
- [ ] Unit tests pass
