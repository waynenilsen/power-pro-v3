# Migration: Backfill Existing Data

## Overview

Create migration `00025_backfill_state_machine_data.sql` to handle backwards compatibility for existing data.

## Tasks

### 1. Backfill user_program_states

For existing user_program_states that were created before the state machine fields were added:

- `enrollment_status` already defaults to 'ACTIVE' (no action needed)
- `cycle_status`: Set to 'IN_PROGRESS' for rows with `current_day_index IS NOT NULL`
- `week_status`: Set to 'IN_PROGRESS' for rows with `current_day_index IS NOT NULL`

### 2. Document logged_sets backwards compatibility

The `logged_sets` table has a `session_id` column that contains client-generated session IDs that do not correspond to the new `workout_sessions` table.

**Decision: Soft Reference (Option A)**
- No FK constraint between logged_sets.session_id and workout_sessions.id
- Existing session_ids remain orphan - they were client-generated
- New session_ids will map to workout_sessions.id
- Document this in migration comments

### 3. Migration file

Create `/migrations/00025_backfill_state_machine_data.sql` with:

```sql
-- +goose Up
-- Backfill state machine status fields for existing enrollments

-- Set cycle_status and week_status to IN_PROGRESS for enrollments that have started
-- (indicated by current_day_index being set)
UPDATE user_program_states
SET cycle_status = 'IN_PROGRESS',
    week_status = 'IN_PROGRESS'
WHERE current_day_index IS NOT NULL
  AND cycle_status = 'PENDING';

-- Note on logged_sets.session_id:
-- Pre-migration logged_sets have client-generated session_ids that do not map to
-- the workout_sessions table. This is intentional - we use a soft reference approach:
-- - No FK constraint between logged_sets.session_id and workout_sessions.id
-- - Legacy session_ids remain valid for historical data
-- - New sets can optionally reference workout_sessions

-- +goose Down
-- Revert back to PENDING status (this is a best-effort rollback)
UPDATE user_program_states
SET cycle_status = 'PENDING',
    week_status = 'PENDING'
WHERE cycle_status = 'IN_PROGRESS'
   OR week_status = 'IN_PROGRESS';
```

## Acceptance Criteria

- [ ] Migration file created at `migrations/00025_backfill_state_machine_data.sql`
- [ ] Migration runs without errors: `goose up`
- [ ] Existing enrollments with current_day_index are updated to IN_PROGRESS
- [ ] Soft reference approach documented in migration
- [ ] All tests still pass
