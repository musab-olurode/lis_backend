-- name: CreatePost :one
INSERT INTO posts (id, cover_image_url, title, slug, description, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPostBySlug :one
SELECT * FROM posts WHERE slug = $1;

-- name: GetPaginatedPosts :many
SELECT * FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts;

-- name: UpdatePost :one
UPDATE posts SET cover_image_url = $2, title = $3, slug = $4, description = $5, content = $6, updated_at = $7
WHERE id = $1 RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;