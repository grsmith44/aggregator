-- +goose Up
ALTER TABLE users
ADD constraint USER_NAME_UNIQUE UNIQUE (user_name);

-- +goose Down
ALTER TABLE users
DROP constraint USER_NAME_UNIQUE;
