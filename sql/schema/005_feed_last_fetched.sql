-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP DEFAULT NULL;


-- +goose Down
ALTER TABLE
drop COLUMN last_fetched_at;
