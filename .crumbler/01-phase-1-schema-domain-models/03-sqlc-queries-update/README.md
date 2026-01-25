# SQLC Queries Update

## Task

Update SQLC queries to include new status columns and add workout_sessions queries.

## Updates Required

### 1. Update `internal/db/queries/user_program_states.sql`

Add new status columns to all SELECT queries:
- `enrollment_status`
- `cycle_status`
- `week_status`

Update these queries:
- GetUserProgramStateByUserID
- GetUserProgramStateByID
- GetEnrollmentWithProgram
- GetStateAdvancementContext

Update INSERT to include status fields with defaults.
Update UPDATE to allow updating status fields.

### 2. Create `internal/db/queries/workout_sessions.sql`

New queries needed:
```sql
-- name: CreateWorkoutSession :exec
-- name: GetWorkoutSessionByID :one
-- name: GetActiveWorkoutSession :one (get IN_PROGRESS session for user_program_state_id)
-- name: GetWorkoutSessionsByState :many
-- name: UpdateWorkoutSessionStatus :exec
-- name: CompleteWorkoutSession :exec (sets finished_at and status=COMPLETED)
-- name: AbandonWorkoutSession :exec (sets finished_at and status=ABANDONED)
-- name: DeleteWorkoutSession :exec
```

### 3. Regenerate SQLC

Run `sqlc generate` to regenerate the Go code.

## Verification

- `sqlc generate` runs without errors
- Generated code includes new status fields
- Generated code includes workout_sessions queries

## Done When

- All queries updated with status fields
- workout_sessions.sql created with CRUD queries
- `sqlc generate` succeeds
- Repository layer can be updated to use new fields (next crumbs)
