-- name: AddNewBook :one
INSERT INTO books (id, created_at, updated_at, name, url, description, author, book_cover)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteBookByID :one
DELETE FROM books WHERE id=$1
RETURNING *;

-- name: GetAllBooks :many
SELECT * FROM books;

-- name: GetBookById :one
SELECT * FROM books WHERE id=$1;