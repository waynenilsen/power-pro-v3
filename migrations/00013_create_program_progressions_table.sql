-- +goose Up
-- +goose StatementBegin
CREATE TABLE program_progressions (
    id TEXT PRIMARY KEY,
    program_id TEXT NOT NULL,
    progression_id TEXT NOT NULL,
    lift_id TEXT,
    priority INTEGER NOT NULL DEFAULT 0,
    enabled INTEGER NOT NULL DEFAULT 1 CHECK(enabled IN (0, 1)),
    override_increment REAL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE CASCADE,
    FOREIGN KEY (progression_id) REFERENCES progressions(id) ON DELETE CASCADE,
    FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- Unique constraint: same progression can't be assigned twice for the same lift in a program
-- Note: SQLite treats NULL as distinct in unique constraints, so (prog, progression, NULL) is unique from other NULLs
-- We use COALESCE to treat NULL lift_id as a special value for uniqueness
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_program_progressions_unique ON program_progressions(
    program_id,
    progression_id,
    COALESCE(lift_id, '00000000-0000-0000-0000-000000000000')
);
-- +goose StatementEnd

-- Index for program lookup (get all progressions for a program)
-- +goose StatementBegin
CREATE INDEX idx_program_progressions_program_id ON program_progressions(program_id);
-- +goose StatementEnd

-- Index for lift-specific progression lookup (get progressions for a specific lift in a program)
-- +goose StatementBegin
CREATE INDEX idx_program_progressions_program_lift ON program_progressions(program_id, lift_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_program_progressions_program_lift;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_program_progressions_program_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_program_progressions_unique;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS program_progressions;
-- +goose StatementEnd
