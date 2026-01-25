# Domain UserProgramState Status Fields

## Task

Update the UserProgramState domain model to include status fields and update the repository.

## Changes to `internal/domain/userprogramstate/userprogramstate.go`

### 1. Add Status Types

```go
type EnrollmentStatus string
type CycleStatus string
type WeekStatus string

const (
    EnrollmentStatusActive        EnrollmentStatus = "ACTIVE"
    EnrollmentStatusBetweenCycles EnrollmentStatus = "BETWEEN_CYCLES"
    EnrollmentStatusQuit          EnrollmentStatus = "QUIT"

    CycleStatusPending    CycleStatus = "PENDING"
    CycleStatusInProgress CycleStatus = "IN_PROGRESS"
    CycleStatusCompleted  CycleStatus = "COMPLETED"

    WeekStatusPending    WeekStatus = "PENDING"
    WeekStatusInProgress WeekStatus = "IN_PROGRESS"
    WeekStatusCompleted  WeekStatus = "COMPLETED"
)
```

### 2. Add Fields to UserProgramState struct

```go
EnrollmentStatus EnrollmentStatus
CycleStatus      CycleStatus
WeekStatus       WeekStatus
```

### 3. Add Validation Functions

- `ValidateEnrollmentStatus(status EnrollmentStatus) validation.Result`
- `ValidateCycleStatus(status CycleStatus) validation.Result`
- `ValidateWeekStatus(status WeekStatus) validation.Result`

### 4. Update `EnrollUser()` function

Set initial statuses:
- EnrollmentStatus: ACTIVE
- CycleStatus: PENDING
- WeekStatus: PENDING

### 5. Update `Validate()` method

Add calls to validate status fields.

## Changes to Repository

Update `internal/repository/user_program_state_repository.go`:
- Add conversion between domain status types and string
- Update all methods to handle status fields

## Tests

Update `internal/domain/userprogramstate/userprogramstate_test.go`:
- Add tests for ValidateEnrollmentStatus
- Add tests for ValidateCycleStatus
- Add tests for ValidateWeekStatus
- Update EnrollUser tests to verify initial statuses

## Done When

- UserProgramState struct has status fields
- Validation functions exist and are tested
- EnrollUser sets correct initial statuses
- Repository handles status fields correctly
- All existing tests still pass
