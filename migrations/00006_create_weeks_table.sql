-- +goose Up
-- +goose StatementBegin
CREATE TABLE weeks (
    id TEXT PRIMARY KEY,
    week_number INTEGER NOT NULL CHECK(week_number >= 1),
    variant TEXT CHECK(variant IS NULL OR variant IN ('A', 'B')),
    cycle_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- Unique constraint on week_number within a cycle (composite unique)
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_weeks_cycle_id_week_number ON weeks(cycle_id, week_number);
-- +goose StatementEnd

-- Index for cycle_id lookups
-- +goose StatementBegin
CREATE INDEX idx_weeks_cycle_id ON weeks(cycle_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE week_days (
    id TEXT PRIMARY KEY,
    week_id TEXT NOT NULL,
    day_id TEXT NOT NULL,
    day_of_week TEXT NOT NULL CHECK(day_of_week IN ('MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY')),
    created_at TEXT NOT NULL,
    FOREIGN KEY (week_id) REFERENCES weeks(id) ON DELETE CASCADE,
    FOREIGN KEY (day_id) REFERENCES days(id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- Index for efficient lookups by week
-- +goose StatementBegin
CREATE INDEX idx_week_days_week_id ON week_days(week_id);
-- +goose StatementEnd

-- Index for finding which weeks use a day
-- +goose StatementBegin
CREATE INDEX idx_week_days_day_id ON week_days(day_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_week_days_day_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_week_days_week_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS week_days;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_weeks_cycle_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_weeks_cycle_id_week_number;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS weeks;
-- +goose StatementEnd
