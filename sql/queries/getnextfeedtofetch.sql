-- name: GetNextFeedToFetch :many

SELECT * FROM feeds
	ORDER BY last_fetched_at NULLS FIRST, updated_at NULLS FIRST;
