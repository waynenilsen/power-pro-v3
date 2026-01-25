# Update Enrollment Events

## Overview

Update enrollment creation and deletion to emit events and use state machines.

## Endpoints to Update

```
POST   /users/{id}/program    # Enrollment creation
DELETE /users/{id}/program    # Unenrollment
```

## Implementation Details

### POST /users/{id}/program (Enroll)
- Use state machine for initial state setup
- Initial states: enrollment=ACTIVE, cycle=PENDING, week=PENDING
- Emit `ENROLLED` event after successful creation
- Event payload:
  - `programId`
  - `userId`
  - `enrolledAt`

### DELETE /users/{id}/program (Unenroll)
- Use state machine to transition to QUIT state (if not already)
- Emit `QUIT` event before deletion
- Event payload:
  - `programId`
  - `userId`
  - `cyclesCompleted`
  - `weeksCompleted`

## Files to Modify

- `internal/api/enrollment_handler.go` - Update handlers
- `internal/api/enrollment_handler_test.go` - Update tests

## Event Bus Integration

The event bus is at `internal/domain/event/bus.go`. Use:
```go
bus := event.NewEventBus()
bus.Publish(ctx, event.StateEvent{
    Type:      event.EventEnrolled,
    UserID:    userID,
    ProgramID: programID,
    Timestamp: time.Now(),
    Payload:   map[string]interface{}{...},
})
```

## Acceptance Criteria

- [ ] ENROLLED event emitted on enrollment creation
- [ ] QUIT event emitted on unenrollment
- [ ] State machine used for state transitions
- [ ] Event payloads include required fields
- [ ] Unit tests pass
