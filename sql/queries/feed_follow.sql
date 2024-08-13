-- name: CreateFeedFollow :one
INSERT INTO feed_follow (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: GetFeedFollowsByUserId :many
SELECT * from feed_follow where user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE from feed_follow where id = $1;