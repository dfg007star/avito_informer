-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT NOT NULL,
    uid VARCHAR(255) UNIQUE,
    title TEXT,
    description TEXT,
    url TEXT,
    preview_url TEXT,
    price INT,
    is_notify BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_items_links FOREIGN KEY (link_id)
        REFERENCES links (id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS items;
