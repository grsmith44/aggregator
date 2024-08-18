-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, feed_name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: SelectAllFeeds :many
SELECT * FROM feeds;


-- name: GetNextFeedToFetch :many
SELECT * 
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT sqlc.arg(n);

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(idx)
RETURNING *;
