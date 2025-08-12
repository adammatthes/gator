-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name as username FROM feeds
	INNER JOIN users ON users.id = feeds.user_id;
