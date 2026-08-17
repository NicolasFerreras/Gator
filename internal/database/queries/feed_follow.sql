-- name: CreateFeedFollow :one 
WITH inserted_feed_follow AS (
    INSERT INTO feed_follow (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.username AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;



-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follow.*,
    feeds.name AS feed_name,
    users.username AS user_name
FROM feed_follow
INNER JOIN feeds ON feed_follow.feed_id = feeds.id
INNER JOIN users ON feed_follow.user_id = users.id
WHERE feed_follow.user_id = $1;

-- name: DeleteFeedFollowByUserID :exec
DELETE FROM feed_follow
WHERE user_id = $1 AND feed_id = $2;