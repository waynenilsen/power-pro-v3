-- +goose Up
-- +goose StatementBegin
CREATE TABLE lifts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    slug TEXT NOT NULL UNIQUE CHECK(slug GLOB '[a-z0-9-]*' AND length(slug) > 0 AND length(slug) <= 100),
    is_competition_lift INTEGER NOT NULL DEFAULT 0,
    parent_lift_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (parent_lift_id) REFERENCES lifts(id) ON DELETE SET NULL,
    CHECK(parent_lift_id IS NULL OR parent_lift_id != id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_lifts_parent_lift_id ON lifts(parent_lift_id);
-- +goose StatementEnd

-- Seed competition lifts with deterministic UUIDs
-- +goose StatementBegin
INSERT INTO lifts (id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at) VALUES
    ('00000000-0000-0000-0000-000000000001', 'Squat', 'squat', 1, NULL, datetime('now'), datetime('now')),
    ('00000000-0000-0000-0000-000000000002', 'Bench Press', 'bench-press', 1, NULL, datetime('now'), datetime('now')),
    ('00000000-0000-0000-0000-000000000003', 'Deadlift', 'deadlift', 1, NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_lifts_parent_lift_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS lifts;
-- +goose StatementEnd
