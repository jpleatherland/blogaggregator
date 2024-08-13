-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  user_id VARCHAR(64) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (api_key) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
