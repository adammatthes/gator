-- +goose Up
ALTER TABLE users ALTER COLUMN id TYPE UUID USING id::text::uuid;

-- +goose Down
ALTER TABLE users ALTER COLUMN id TYPE VARCHAR(255);
