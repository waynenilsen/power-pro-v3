-- +goose Up
-- Add weight_unit preference column to users table
-- Stores user's preferred unit for weight display (lb or kg)

-- +goose StatementBegin
ALTER TABLE users ADD COLUMN weight_unit TEXT NOT NULL DEFAULT 'lb' CHECK(weight_unit IN ('lb', 'kg'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN weight_unit;
-- +goose StatementEnd
