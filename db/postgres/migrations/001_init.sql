-- +goose Up
CREATE TABLE feature_flags (
    id uuid NOT NULL,
    description text,
    name text,
    enabled bool,
    tenant text,
    created_at TIMESTAMPTZ DEFAULT TIMEZONE('utc', NOW()),
    created_by text,
    updated_at TIMESTAMPTZ DEFAULT TIMEZONE('utc', NOW()),
    updated_by text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE feature_flags;
