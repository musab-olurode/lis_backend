-- +goose Up

CREATE TYPE USER_ROLE as enum (
    'USER',
    'ADMIN'
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    matric_number VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role USER_ROLE NOT NULL DEFAULT 'USER',
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down

DROP TABLE users;