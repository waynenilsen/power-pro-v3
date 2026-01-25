# Phase 2: Event System

## Overview

Create an in-memory event bus to decouple state transitions from progression triggers. Events will fire when state changes occur, and the progression system will subscribe to relevant events.

## Tasks

### 1. Domain: Create event types

Create `internal/domain/event/event.go`:

```go
type EventType string

const (
    EventEnrolled             EventType = "ENROLLED"
    EventCycleBoundaryReached EventType = "CYCLE_BOUNDARY_REACHED"
    EventQuit                 EventType = "QUIT"
    EventCycleStarted         EventType = "CYCLE_STARTED"
    EventCycleCompleted       EventType = "CYCLE_COMPLETED"
    EventWeekStarted          EventType = "WEEK_STARTED"
    EventWeekCompleted        EventType = "WEEK_COMPLETED"
    EventWorkoutStarted       EventType = "WORKOUT_STARTED"
    EventWorkoutCompleted     EventType = "WORKOUT_COMPLETED"
    EventWorkoutAbandoned     EventType = "WORKOUT_ABANDONED"
    EventSetLogged            EventType = "SET_LOGGED"
)

type StateEvent struct {
    Type       EventType
    UserID     string
    ProgramID  string
    Timestamp  time.Time
    Payload    map[string]interface{}
}
```

### 2. Domain: Create event bus

Create `internal/domain/event/bus.go`:
- In-memory pub/sub implementation
- `Subscribe(eventType EventType, handler EventHandler)`
- `Publish(event StateEvent)`
- Thread-safe with mutex

### 3. Domain: Create event handlers registry

Create `internal/domain/event/handlers.go`:
- Registration pattern for handlers
- Handler interface definition

### 4. Integration: Wire events to progression system

Update progression service to:
- Subscribe to `SET_LOGGED`, `WORKOUT_COMPLETED`, `WEEK_COMPLETED`, `CYCLE_COMPLETED` events
- Map events to progression trigger types:
  - `SET_LOGGED` → `AFTER_SET`, `ON_FAILURE`
  - `WORKOUT_COMPLETED` → `AFTER_SESSION`
  - `WEEK_COMPLETED` → `AFTER_WEEK`
  - `CYCLE_COMPLETED` → `AFTER_CYCLE`

## Acceptance Criteria

- [ ] Event types defined for all state transitions
- [ ] Event bus supports subscribe/publish pattern
- [ ] Event bus is thread-safe
- [ ] Progression service subscribes to relevant events
- [ ] Events replace manual progression trigger calls
- [ ] Unit tests for event bus
- [ ] Integration tests for event → progression flow
