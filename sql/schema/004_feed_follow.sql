-- +goose Up
CREATE TABLE feed_follow (
  id UUID PRIMARY KEY,
  feed_id UUID NOT NULL,
  user_id VARCHAR(64) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (api_key) ON DELETE CASCADE,
  FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follow;
