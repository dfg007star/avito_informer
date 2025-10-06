-- +goose Up
ALTER TABLE links
    ADD COLUMN min_price INTEGER,
ADD COLUMN max_price INTEGER;

-- +goose Down
ALTER TABLE links
DROP COLUMN min_price,
DROP COLUMN max_price;
