-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;


-- name: GetPostsByUser :many
SELECT p.*
FROM posts p
INNER JOIN feed_follows ff ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY
  CASE 
    WHEN p.published_at IS NOT NULL THEN p.published_at
    ELSE p.created_at
  END DESC;
