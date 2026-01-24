-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_program_states ADD COLUMN meet_date TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_program_states ADD COLUMN schedule_type TEXT DEFAULT 'rotation';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_program_states DROP COLUMN schedule_type;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_program_states DROP COLUMN meet_date;
-- +goose StatementEnd
