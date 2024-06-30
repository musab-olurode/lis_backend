-- name: CreateMaterial :one
INSERT INTO materials (id, title, file_url, category, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetMaterialByID :one
SELECT * FROM materials WHERE id = $1;

-- name: GetPaginatedMaterials :many
SELECT * FROM materials ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetPaginatedMaterialsByCategory :many
SELECT * FROM materials WHERE category = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountMaterials :one
SELECT COUNT(*) FROM materials;

-- name: CountMaterialsByCategory :one
SELECT COUNT(*) FROM materials WHERE category = $1;

-- name: UpdateMaterial :one
UPDATE materials SET title = $2, file_url = $3, category = $4, updated_at = $5
WHERE id = $1 RETURNING *;

-- name: DeleteMaterial :exec
DELETE FROM materials WHERE id = $1;