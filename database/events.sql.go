// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: events.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const countEvents = `-- name: CountEvents :one
SELECT COUNT(*) FROM events
`

func (q *Queries) CountEvents(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countEvents)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countUpcomingEvents = `-- name: CountUpcomingEvents :one
SELECT COUNT(*) FROM events WHERE date >= NOW()
`

func (q *Queries) CountUpcomingEvents(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUpcomingEvents)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (id, title, description, image_url, venue, date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, title, description, image_url, venue, date, created_at, updated_at
`

type CreateEventParams struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	ImageUrl    string         `json:"image_url"`
	Venue       string         `json:"venue"`
	Date        time.Time      `json:"date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, createEvent,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.ImageUrl,
		arg.Venue,
		arg.Date,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.ImageUrl,
		&i.Venue,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteEvent = `-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1
`

func (q *Queries) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteEvent, id)
	return err
}

const getEventByID = `-- name: GetEventByID :one
SELECT id, title, description, image_url, venue, date, created_at, updated_at FROM events WHERE id = $1
`

func (q *Queries) GetEventByID(ctx context.Context, id uuid.UUID) (Event, error) {
	row := q.db.QueryRowContext(ctx, getEventByID, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.ImageUrl,
		&i.Venue,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPaginatedEvents = `-- name: GetPaginatedEvents :many
SELECT id, title, description, image_url, venue, date, created_at, updated_at FROM events ORDER BY created_at DESC LIMIT $1 OFFSET $2
`

type GetPaginatedEventsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetPaginatedEvents(ctx context.Context, arg GetPaginatedEventsParams) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getPaginatedEvents, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Event{}
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.ImageUrl,
			&i.Venue,
			&i.Date,
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

const getUpcomingEventsPaginated = `-- name: GetUpcomingEventsPaginated :many
SELECT id, title, description, image_url, venue, date, created_at, updated_at FROM events WHERE date >= NOW() ORDER BY date ASC LIMIT $1 OFFSET $2
`

type GetUpcomingEventsPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUpcomingEventsPaginated(ctx context.Context, arg GetUpcomingEventsPaginatedParams) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getUpcomingEventsPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Event{}
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.ImageUrl,
			&i.Venue,
			&i.Date,
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

const updateEvent = `-- name: UpdateEvent :one
UPDATE events SET title = $2, description = $3, image_url = $4, venue = $5, created_at = $6, updated_at = $7
WHERE id = $1 RETURNING id, title, description, image_url, venue, date, created_at, updated_at
`

type UpdateEventParams struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	ImageUrl    string         `json:"image_url"`
	Venue       string         `json:"venue"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, updateEvent,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.ImageUrl,
		arg.Venue,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.ImageUrl,
		&i.Venue,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
