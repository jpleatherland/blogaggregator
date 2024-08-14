-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * from feeds;

-- name: GetNextFeedsToFetch :many
SELECT * from feeds order by last_fetched_at asc limit $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = TO_CHAR(CURRENT_TIMESTAMP, 'YYYY-MM-DD"T"HH24:MI:SS'),
    last_fetched_at = TO_CHAR(CURRENT_TIMESTAMP, 'YYYY-MM-DD"T"HH24:MI:SS')
WHERE id = $1;
