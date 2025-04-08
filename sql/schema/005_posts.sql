-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id SERIAL NOT NULL,
    CONSTRAINT fk_posts_feeds
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
