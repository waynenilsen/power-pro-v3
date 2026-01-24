# Meet Date Domain Model

Add meet date tracking to UserProgramState domain model.

## Tasks

1. Add `MeetDate` field to `UserProgramState` struct (pointer to time.Time, nullable)
2. Add `ScheduleType` field to `UserProgramState` (enum: "rotation" or "days_out")
3. Add validation for MeetDate (must be in the future when set)
4. Update `EnrollUserInput` to accept optional meet date
5. Update database schema with migration (add columns to user_program_states table)
6. Update sqlc queries for the new fields
7. Regenerate sqlc models

## Schema Changes

```sql
ALTER TABLE user_program_states ADD COLUMN meet_date TEXT; -- ISO 8601 date, nullable
ALTER TABLE user_program_states ADD COLUMN schedule_type TEXT DEFAULT 'rotation'; -- 'rotation' or 'days_out'
```

## Acceptance Criteria

- [ ] `MeetDate` field exists on UserProgramState
- [ ] `ScheduleType` field exists with rotation/days_out options
- [ ] Validation rejects past dates for MeetDate
- [ ] Database schema updated
- [ ] Existing tests still pass
