// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, user_name, api_key)
VALUES ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex'))
RETURNING id, created_at, updated_at, user_name, api_key
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserName  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserName,
		&i.ApiKey,
	)
	return i, err
}

const getUserAPIKeyByName = `-- name: GetUserAPIKeyByName :one
SELECT api_key FROM users
WHERE user_name = $1
`

func (q *Queries) GetUserAPIKeyByName(ctx context.Context, userName string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserAPIKeyByName, userName)
	var api_key string
	err := row.Scan(&api_key)
	return api_key, err
}

const getUserByAPI = `-- name: GetUserByAPI :one
SELECT id, created_at, updated_at, user_name, api_key FROM users
WHERE api_key = $1
`

func (q *Queries) GetUserByAPI(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByAPI, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserName,
		&i.ApiKey,
	)
	return i, err
}

const getUserByUserName = `-- name: GetUserByUserName :one
SELECT id, created_at, updated_at, user_name, api_key FROM users
WHERE user_name = $1
`

func (q *Queries) GetUserByUserName(ctx context.Context, userName string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUserName, userName)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserName,
		&i.ApiKey,
	)
	return i, err
}
