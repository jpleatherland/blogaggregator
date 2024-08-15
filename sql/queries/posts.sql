-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPostsByUser :many
SELECT feed_follow.feed_id, posts.*
from feed_follow 
INNER JOIN posts ON feed_follow.feed_id = posts.feed_id
where feed_follow.user_id = $1 
LIMIT $2;
