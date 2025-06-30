-- +goose Up
CREATE TABLE books (
    id  UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    author TEXT NOT NULL,
    book_cover TEXT NOT NULL
);

-- +goose Down
DROP TABLE books;