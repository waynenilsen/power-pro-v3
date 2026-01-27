-- +goose Up
-- +goose StatementBegin
-- Backfill Training Maxes for all 1RMs that don't have a corresponding TM.
-- TM = 1RM * 0.9, rounded to nearest 0.25
INSERT INTO lift_maxes (id, user_id, lift_id, type, value, effective_date, created_at, updated_at)
SELECT 
    lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)), 2) || '-' || substr('89ab', abs(random()) % 4 + 1, 1) || substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6))) as id,
    orm.user_id,
    orm.lift_id,
    'TRAINING_MAX' as type,
    ROUND(orm.value * 0.9 / 0.25) * 0.25 as value,
    orm.effective_date,
    datetime('now') as created_at,
    datetime('now') as updated_at
FROM lift_maxes orm
WHERE orm.type = 'ONE_RM'
AND NOT EXISTS (
    SELECT 1 FROM lift_maxes tm 
    WHERE tm.user_id = orm.user_id 
    AND tm.lift_id = orm.lift_id 
    AND tm.type = 'TRAINING_MAX'
    AND tm.effective_date = orm.effective_date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove auto-generated TMs (those created by this migration)
-- We can identify them by having same effective_date as a 1RM
DELETE FROM lift_maxes 
WHERE type = 'TRAINING_MAX'
AND EXISTS (
    SELECT 1 FROM lift_maxes orm
    WHERE orm.user_id = lift_maxes.user_id
    AND orm.lift_id = lift_maxes.lift_id
    AND orm.type = 'ONE_RM'
    AND orm.effective_date = lift_maxes.effective_date
);
-- +goose StatementEnd
