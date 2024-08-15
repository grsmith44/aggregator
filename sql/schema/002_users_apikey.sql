-- +goose Up
ALTER TABLE users
ADD COLUMN api_key VARCHAR(64) NOT NULL DEFAULT encode(
    sha256(random()::TEXT::BYTEA), 'hex'
);

ALTER TABLE users
ADD constraint API_KEY_UNIQUE UNIQUE (api_key);

-- +goose Down
ALTER TABLE
drop COLUMN api_key;
