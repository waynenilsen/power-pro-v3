-- +goose Up
-- Fix idempotency index to include lift_id
-- This allows the same progression to be applied to different lifts at the same timestamp

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_progression_logs_idempotency;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX idx_progression_logs_idempotency
    ON progression_logs(user_id, progression_id, lift_id, trigger_type, applied_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_progression_logs_idempotency;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX idx_progression_logs_idempotency
    ON progression_logs(user_id, progression_id, trigger_type, applied_at);
-- +goose StatementEnd
