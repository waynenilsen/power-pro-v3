# Phase 3: API Layer

## Overview

Add HTTP endpoints for workout session management and enrollment state operations. Update existing endpoints to use state machines and emit events.

## Tasks

### 1. Handler: Workout session endpoints

Create workout session handlers:

```
POST   /workouts/start           # Start new workout, returns session_id
GET    /workouts/{id}            # Get workout session details
POST   /workouts/{id}/finish     # Mark workout complete
POST   /workouts/{id}/abandon    # Mark workout abandoned
GET    /users/{id}/workouts      # List user's workout history
GET    /users/{id}/workouts/current  # Get current in-progress workout
```

Implementation details:
- Starting workout creates IN_PROGRESS session
- Can't start if one already IN_PROGRESS
- Finishing triggers state transitions (week/cycle completion checks)
- Emit `WORKOUT_STARTED`, `WORKOUT_COMPLETED`, `WORKOUT_ABANDONED` events

### 2. Handler: Enrollment state endpoints

Add/update enrollment handlers:

```
GET    /users/{id}/enrollment             # Include all state fields
POST   /users/{id}/enrollment/next-cycle  # Start new cycle (when BETWEEN_CYCLES)
POST   /users/{id}/enrollment/advance-week  # Manual week advance
```

Response shape must include:
- `enrollment_status`: ACTIVE | BETWEEN_CYCLES | QUIT
- `cycle_status`: PENDING | IN_PROGRESS | COMPLETED
- `week_status`: PENDING | IN_PROGRESS | COMPLETED
- `current_workout_session`: {...} or null

### 3. Handler: Modify logged sets endpoint

Update `POST /sessions/{id}/sets`:
- Require active workout session
- Emit `SET_LOGGED` event after logging
- Include failure detection for `ON_FAILURE` trigger

### 4. Handler: Update existing enrollment endpoint

Update enrollment creation/deletion:
- Emit `ENROLLED` event on creation
- Emit `QUIT` event on deletion
- Use state machine for transitions

### 5. Error responses for invalid state transitions

Return clear errors when actions not allowed:

```json
{
  "error": "workout_already_in_progress",
  "message": "Complete or abandon current workout before starting a new one",
  "current_workout_session_id": "..."
}
```

## Acceptance Criteria

- [ ] All workout session endpoints implemented
- [ ] Enrollment state endpoints implemented
- [ ] API responses include all state fields
- [ ] Events emitted on state changes
- [ ] Clear error messages for invalid transitions
- [ ] OpenAPI/swagger docs updated
- [ ] Handler unit tests pass
