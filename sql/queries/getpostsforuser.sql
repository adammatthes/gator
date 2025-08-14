-- name: GetPostsForUser :many

SELECT * FROM posts
	WHERE posts.feed_id IN
	(
	SELECT feed_follows.feed_id FROM feed_follows
		WHERE feed_follows.user_id = $1
	)
	ORDER BY posts.published_at DESC
	LIMIT $2;
