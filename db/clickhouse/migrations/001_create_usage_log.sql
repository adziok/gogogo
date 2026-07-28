-- +goose Up
-- The table is ordered for the primary analytics query: tenant and time range,
-- with flag_name available as a further filter within that range.
CREATE TABLE IF NOT EXISTS usage_log
(
    tenant      LowCardinality(String),
    flag_name   LowCardinality(String),
    enabled     Bool,
    called_at   DateTime64(3, 'UTC'),
    ingested_at DateTime64(3, 'UTC') DEFAULT now64(3)
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(called_at)
ORDER BY (tenant, called_at, flag_name);

-- +goose Down
DROP TABLE IF EXISTS usage_log;
