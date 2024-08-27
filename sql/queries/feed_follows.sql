-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :execresult
DELETE FROM feed_follows
WHERE id = $1;

-- name: GetAllFeedFollowsForUser :many
SELECT * FROM feed_follows
WHERE user_id = $1;

-- name: GetAllFeedNamesFollowedByUser :many
SELECT f.feed_name
FROM feeds f
INNER JOIN feed_follows ff ON f.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY ff.created_at DESC;
