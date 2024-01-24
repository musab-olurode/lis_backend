-- name: CreateUser :one
INSERT INTO users (id, matric_number, first_name, last_name, role, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByMatricNumber :one
SELECT * FROM users WHERE matric_number = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;