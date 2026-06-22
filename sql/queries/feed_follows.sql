-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES($1, $2, $3, $4, $5)
    RETURNING *
)

SELECT 
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
WITH feed_follows_of_user AS (
    SELECT * FROM feed_follows
    WHERE feed_follows.user_id = $1
)
SELECT
    feed_follows_of_user.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows_of_user
INNER JOIN feeds ON feeds.id = feed_follows_of_user.feed_id
INNER JOIN users ON users.id = feed_follows_of_user.user_id;
