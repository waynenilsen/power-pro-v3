-- +goose Up
-- +goose StatementBegin
CREATE TABLE prescriptions (
    id TEXT PRIMARY KEY,
    lift_id TEXT NOT NULL,
    load_strategy TEXT NOT NULL,
    set_scheme TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    notes TEXT,
    rest_seconds INTEGER,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE RESTRICT,
    CHECK(notes IS NULL OR length(notes) <= 500),
    CHECK(rest_seconds IS NULL OR rest_seconds >= 0)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_prescriptions_lift_id ON prescriptions(lift_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_prescriptions_lift_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS prescriptions;
-- +goose StatementEnd
