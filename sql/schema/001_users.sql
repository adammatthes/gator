-- +goose up
CREATE TABLE users(
	id INTEGER PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name VARCHAR(255) UNIQUE NOT NULL	
);

-- +goose down
DROP TABLE users;
