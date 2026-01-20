-- +goose Up
-- +goose StatementBegin
CREATE TABLE days (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 50),
    slug TEXT NOT NULL CHECK(slug GLOB '[a-z0-9-]*' AND length(slug) > 0 AND length(slug) <= 50),
    metadata TEXT,
    program_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CHECK(metadata IS NULL OR json_valid(metadata))
);
-- +goose StatementEnd

-- Unique constraint on slug within a program (composite unique)
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_days_program_id_slug ON days(program_id, slug);
-- +goose StatementEnd

-- Index for program_id lookups
-- +goose StatementBegin
CREATE INDEX idx_days_program_id ON days(program_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE day_prescriptions (
    id TEXT PRIMARY KEY,
    day_id TEXT NOT NULL,
    prescription_id TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    FOREIGN KEY (day_id) REFERENCES days(id) ON DELETE CASCADE,
    FOREIGN KEY (prescription_id) REFERENCES prescriptions(id) ON DELETE RESTRICT,
    CHECK("order" >= 0)
);
-- +goose StatementEnd

-- Index for efficient lookups by day
-- +goose StatementBegin
CREATE INDEX idx_day_prescriptions_day_id ON day_prescriptions(day_id);
-- +goose StatementEnd

-- Index for finding which days use a prescription
-- +goose StatementBegin
CREATE INDEX idx_day_prescriptions_prescription_id ON day_prescriptions(prescription_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_day_prescriptions_prescription_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_day_prescriptions_day_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS day_prescriptions;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_days_program_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_days_program_id_slug;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS days;
-- +goose StatementEnd
