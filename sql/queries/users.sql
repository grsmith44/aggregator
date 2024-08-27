-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, user_name, api_key)
VALUES ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: GetUserByAPI :one
SELECT * FROM users
WHERE api_key = $1;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE user_name = $1;

-- name: GetUserAPIKeyByName :one
SELECT api_key FROM users
WHERE user_name = $1;

