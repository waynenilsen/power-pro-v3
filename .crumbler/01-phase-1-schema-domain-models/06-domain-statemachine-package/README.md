# Domain StateMachine Package

## Task

Create state machine package at `internal/domain/statemachine/` with generic interface and specific implementations.

## Files to Create

### 1. `statemachine.go` - Generic Interface

```go
package statemachine

// State represents any state in a state machine
type State string

// Transition represents a valid state transition
type Transition struct {
    From State
    To   State
}

// StateMachine defines the interface for all state machines
type StateMachine interface {
    CurrentState() State
    ValidTransitions() []Transition
    CanTransitionTo(target State) bool
    TransitionTo(target State) error
}

// InvalidTransitionError for invalid state transitions
type InvalidTransitionError struct {
    From State
    To   State
}
```

### 2. `enrollment.go` - Enrollment State Machine

States: ACTIVE, BETWEEN_CYCLES, QUIT

Valid transitions:
- ACTIVE -> BETWEEN_CYCLES (cycle completed)
- ACTIVE -> QUIT (user quits)
- BETWEEN_CYCLES -> ACTIVE (new cycle started)
- BETWEEN_CYCLES -> QUIT (user quits)

### 3. `cycle.go` - Cycle State Machine

States: PENDING, IN_PROGRESS, COMPLETED

Valid transitions:
- PENDING -> IN_PROGRESS (cycle started)
- IN_PROGRESS -> COMPLETED (all weeks done)
- COMPLETED -> PENDING (reset for new cycle)

### 4. `week.go` - Week State Machine

States: PENDING, IN_PROGRESS, COMPLETED

Valid transitions:
- PENDING -> IN_PROGRESS (week started)
- IN_PROGRESS -> COMPLETED (all workouts done)
- COMPLETED -> PENDING (reset for next week)

### 5. `workout.go` - Workout State Machine

States: IN_PROGRESS, COMPLETED, ABANDONED

Valid transitions:
- IN_PROGRESS -> COMPLETED (workout finished)
- IN_PROGRESS -> ABANDONED (workout cancelled)

Note: No transitions FROM COMPLETED or ABANDONED (terminal states)

### 6. Test Files

Create test file for each:
- `statemachine_test.go`
- `enrollment_test.go`
- `cycle_test.go`
- `week_test.go`
- `workout_test.go`

Test all valid and invalid transitions.

## Done When

- All state machine files created
- All tests pass
- Invalid transitions return proper errors
- Good test coverage for all state machines
