-- +goose Up
ALTER TABLE feeds
ADD last_fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;