-- +goose Up
    CREATE TABLE posts (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        title TEXT NOT NULL,
        url TEXT UNIQUE NOT NULL,
        description TEXT NULL,
        published_at TEXT NULL,
        feed_id UUID NOT NULL,
        CONSTRAINT fk_feed
        FOREIGN KEY (feed_id)
        REFERENCES feeds(id) ON DELETE CASCADE
    );
-- +goose Down
   DROP TABLE posts;
