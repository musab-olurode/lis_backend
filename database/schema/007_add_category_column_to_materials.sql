-- +goose Up

create type category as ENUM('100 Level', '200 Level', '300 Level', '400 Level', '500 Level', 'Law Literatures', 'Our Journals', 'Articles');
ALTER TABLE materials ADD COLUMN category category NOT NULL DEFAULT 'Articles';

-- +goose Down

ALTER TABLE materials DROP COLUMN category;