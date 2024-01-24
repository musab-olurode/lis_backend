// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: materials.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countMaterials = `-- name: CountMaterials :one
SELECT COUNT(*) FROM materials
`

func (q *Queries) CountMaterials(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countMaterials)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createMaterial = `-- name: CreateMaterial :one
INSERT INTO materials (id, title, file_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, file_url, created_at, updated_at
`

type CreateMaterialParams struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	FileUrl   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateMaterial(ctx context.Context, arg CreateMaterialParams) (Material, error) {
	row := q.db.QueryRowContext(ctx, createMaterial,
		arg.ID,
		arg.Title,
		arg.FileUrl,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Material
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.FileUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMaterial = `-- name: DeleteMaterial :exec
DELETE FROM materials WHERE id = $1
`

func (q *Queries) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteMaterial, id)
	return err
}

const getMaterialByID = `-- name: GetMaterialByID :one
SELECT id, title, file_url, created_at, updated_at FROM materials WHERE id = $1
`

func (q *Queries) GetMaterialByID(ctx context.Context, id uuid.UUID) (Material, error) {
	row := q.db.QueryRowContext(ctx, getMaterialByID, id)
	var i Material
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.FileUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPaginatedMaterials = `-- name: GetPaginatedMaterials :many
SELECT id, title, file_url, created_at, updated_at FROM materials ORDER BY created_at DESC LIMIT $1 OFFSET $2
`

type GetPaginatedMaterialsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetPaginatedMaterials(ctx context.Context, arg GetPaginatedMaterialsParams) ([]Material, error) {
	rows, err := q.db.QueryContext(ctx, getPaginatedMaterials, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Material
	for rows.Next() {
		var i Material
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.FileUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMaterial = `-- name: UpdateMaterial :one
UPDATE materials SET title = $2, file_url = $3, updated_at = $4
WHERE id = $1 RETURNING id, title, file_url, created_at, updated_at
`

type UpdateMaterialParams struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	FileUrl   string    `json:"file_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateMaterial(ctx context.Context, arg UpdateMaterialParams) (Material, error) {
	row := q.db.QueryRowContext(ctx, updateMaterial,
		arg.ID,
		arg.Title,
		arg.FileUrl,
		arg.UpdatedAt,
	)
	var i Material
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.FileUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
