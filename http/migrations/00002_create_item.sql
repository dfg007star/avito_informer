-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT NOT NULL,
    uid VARCHAR(255),
    title TEXT,
    description TEXT,
    url TEXT,
    preview_url TEXT,
    price INT,
    need_notify BOOLEAN,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_items_links FOREIGN KEY (link_id)
        REFERENCES links (id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS items;
