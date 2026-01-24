-- +goose Up
-- +goose StatementBegin
ALTER TABLE logged_sets ADD COLUMN rpe REAL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE logged_sets DROP COLUMN rpe;
-- +goose StatementEnd
