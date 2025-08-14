-- +goose up

CREATE TABLE posts (
	id UUID UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title Text,
	url TEXT UNIQUE,
	description TEXT,
	published_at TIMESTAMP NOT NULL,
	feed_id UUID
);

-- +goose down
DROP TABLE posts;
