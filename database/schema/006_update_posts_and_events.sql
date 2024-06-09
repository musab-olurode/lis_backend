-- +goose Up

ALTER TABLE posts
RENAME COLUMN image_url TO cover_image_url;
ALTER TABLE posts
RENAME COLUMN description TO content;
ALTER TABLE posts
ADD COLUMN description VARCHAR(255);
ALTER TABLE posts
ADD COLUMN slug VARCHAR(255) UNIQUE NOT NULL;
ALTER TABLE events
ALTER COLUMN description TYPE TEXT;

-- +goose Down
ALTER TABLE posts
RENAME COLUMN cover_image_url TO image_url;
ALTER TABLE posts
DROP COLUMN description;
ALTER TABLE posts
RENAME COLUMN content TO description;
ALTER TABLE posts
DROP COLUMN slug;
ALTER TABLE posts
ALTER COLUMN description TYPE VARCHAR(255);
