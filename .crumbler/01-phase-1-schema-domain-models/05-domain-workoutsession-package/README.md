# Domain WorkoutSession Package

## Task

Create new domain package for WorkoutSession entity at `internal/domain/workoutsession/`.

## Files to Create

### 1. `workoutsession.go`

```go
package workoutsession

import "time"

type WorkoutSessionStatus string

const (
    StatusInProgress WorkoutSessionStatus = "IN_PROGRESS"
    StatusCompleted  WorkoutSessionStatus = "COMPLETED"
    StatusAbandoned  WorkoutSessionStatus = "ABANDONED"
)

type WorkoutSession struct {
    ID                 string
    UserProgramStateID string
    WeekNumber         int
    DayIndex           int
    Status             WorkoutSessionStatus
    StartedAt          time.Time
    FinishedAt         *time.Time
    CreatedAt          time.Time
    UpdatedAt          time.Time
}

func NewWorkoutSession(userProgramStateID string, weekNumber, dayIndex int) (*WorkoutSession, error) {
    // Create new session with IN_PROGRESS status
    // Generate ID, set timestamps
    // Validate and return
}

func (ws *WorkoutSession) Complete() error {
    // Set status to COMPLETED, set FinishedAt
}

func (ws *WorkoutSession) Abandon() error {
    // Set status to ABANDONED, set FinishedAt
}

func (ws *WorkoutSession) Validate() validation.Result {
    // Full validation
}
```

### 2. `validation.go`

```go
package workoutsession

// ValidateID, ValidateUserProgramStateID, ValidateWeekNumber, ValidateDayIndex, ValidateStatus
// Follow same pattern as userprogramstate validation
```

### 3. `workoutsession_test.go`

Comprehensive tests following userprogramstate_test.go pattern:
- Validation tests for each field
- Status transition tests
- Table-driven test style

## Also Create Repository

Create `internal/repository/workout_session_repository.go`:
- Pattern matches user_program_state_repository.go
- CRUD operations using generated sqlc queries

## Done When

- Package created with entity, validation, tests
- Repository created
- All tests pass
- Good test coverage
