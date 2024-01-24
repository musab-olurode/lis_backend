-- +goose Up

CREATE TABLE events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255),
    image_url VARCHAR(255) NOT NULL,
    venue VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down

DROP TABLE events;