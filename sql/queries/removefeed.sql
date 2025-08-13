-- name: RemoveFeed :one
DELETE FROM feed_follows
	WHERE feed_id = $1 AND user_id = $2
	RETURNING *;
