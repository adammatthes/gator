-- name: GetFeedNameById :many
SELECT feeds.name AS Feedname FROM feed_follows
	INNER JOIN users ON users.id = feed_follows.user_id
	INNER JOIN feeds ON feeds.id = feed_follows.feed_id
	WHERE feed_follows.user_id = $1;
