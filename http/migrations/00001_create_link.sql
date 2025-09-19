-- +goose Up
CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    url TEXT,
    parsed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP NULL DEFAULT NULL
);

-- +goose Down
DROP TABLE IF EXISTS links;
