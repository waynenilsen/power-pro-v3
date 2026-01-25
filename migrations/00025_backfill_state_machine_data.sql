-- +goose Up
-- Backfill state machine status fields for existing enrollments
--
-- Migration 00023 added enrollment_status, cycle_status, and week_status fields
-- with defaults of 'ACTIVE', 'PENDING', and 'PENDING' respectively.
--
-- For existing enrollments that have already started (current_day_index IS NOT NULL),
-- we need to update the status fields to reflect their actual state.

-- +goose StatementBegin
UPDATE user_program_states
SET cycle_status = 'IN_PROGRESS',
    week_status = 'IN_PROGRESS'
WHERE current_day_index IS NOT NULL
  AND cycle_status = 'PENDING';
-- +goose StatementEnd

-- Note on logged_sets.session_id backwards compatibility:
--
-- The logged_sets table has a session_id column that contains client-generated
-- session IDs from before the workout_sessions table was introduced (migration 00024).
-- These pre-migration session_ids do not map to any workout_sessions records.
--
-- Design decision: Soft Reference Approach
-- - No FK constraint between logged_sets.session_id and workout_sessions.id
-- - Legacy client-generated session_ids remain valid for historical data
-- - New sets logged through the state machine can optionally reference workout_sessions.id
-- - This preserves data integrity while allowing seamless transition to the new system

-- +goose Down
-- Revert cycle_status and week_status back to PENDING
-- Note: This is a best-effort rollback. We cannot perfectly distinguish between
-- rows that were updated by this migration vs. rows that were already IN_PROGRESS.

-- +goose StatementBegin
UPDATE user_program_states
SET cycle_status = 'PENDING',
    week_status = 'PENDING'
WHERE cycle_status = 'IN_PROGRESS'
   OR week_status = 'IN_PROGRESS';
-- +goose StatementEnd
