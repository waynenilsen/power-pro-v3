# Invalid Transition Errors

## Overview

Return clear, structured error responses when actions are not allowed due to state constraints.

## Error Response Format

```json
{
  "error": "error_code",
  "message": "Human readable message explaining why action is not allowed",
  "details": {
    "field": "value"
  }
}
```

## Error Cases to Handle

### Workout Already In Progress
When trying to start a workout but one is already active:
```json
{
  "error": "workout_already_in_progress",
  "message": "Complete or abandon current workout before starting a new one",
  "details": {
    "current_workout_session_id": "..."
  }
}
```

### No Active Workout
When trying to finish/abandon but no active workout:
```json
{
  "error": "no_active_workout",
  "message": "No workout is currently in progress"
}
```

### Invalid Enrollment State
When trying to start cycle but not BETWEEN_CYCLES:
```json
{
  "error": "invalid_enrollment_state",
  "message": "Cannot start new cycle - enrollment is not between cycles",
  "details": {
    "current_status": "ACTIVE"
  }
}
```

### Session Not Active
When trying to log sets to completed/abandoned session:
```json
{
  "error": "session_not_active",
  "message": "Cannot log sets to a session that is not in progress",
  "details": {
    "session_status": "COMPLETED"
  }
}
```

### Not Enrolled
When enrollment required but user not enrolled:
```json
{
  "error": "not_enrolled",
  "message": "User must be enrolled in a program to perform this action"
}
```

## Implementation

- Create error types in `internal/errors/` if needed
- Use consistent error structure across all handlers
- HTTP status codes:
  - 400 Bad Request for invalid transitions
  - 404 Not Found for missing resources
  - 409 Conflict for state conflicts

## Files to Create/Modify

- `internal/errors/state_errors.go` - New state-specific errors
- Update all Phase 3 handlers to use these errors

## Acceptance Criteria

- [ ] All error cases return structured JSON
- [ ] Error codes are machine-readable (snake_case)
- [ ] Messages are human-readable
- [ ] Details provide context for debugging
- [ ] Consistent across all handlers
