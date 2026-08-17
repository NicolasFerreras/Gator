-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6
)
RETURNING *;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE name = $1;

-- name: GetFeedsByUserID :many
SELECT * FROM feeds WHERE user_id = $1;

-- name: GetFeedByID :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetUserNameByFeedID :one
SELECT u.username FROM users u
JOIN feeds f ON u.id = f.user_id
WHERE f.id = $1;

-- name: GetFeedWithUserName :many
SELECT f.*, u.username FROM feeds f
INNER JOIN users u ON f.user_id = u.id;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;

-- name: DeleteFeedByID :exec
DELETE FROM feeds WHERE id = $1;


-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $1, last_fetched_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT id, url
FROM feeds 
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
