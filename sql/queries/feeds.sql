-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsAll :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :many
INSERT INTO feeds_follow (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds_follow
WHERE id = $1;

-- name: GetUserFeedAll :many
SELECT * from feeds_follow
WHERE user_id = $1;
