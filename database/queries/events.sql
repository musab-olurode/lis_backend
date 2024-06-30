-- name: CreateEvent :one
INSERT INTO events (id, title, description, image_url, venue, date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1;

-- name: GetPaginatedEvents :many
SELECT * FROM events ORDER BY date DESC LIMIT $1 OFFSET $2;

-- name: CountEvents :one
SELECT COUNT(*) FROM events;

-- name: UpdateEvent :one
UPDATE events SET title = $2, description = $3, image_url = $4, venue = $5, date = $6, updated_at = $7
WHERE id = $1 RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;

-- name: CountUpcomingEvents :one
SELECT COUNT(*) FROM events WHERE date >= NOW();

-- name: GetUpcomingEventsPaginated :many
SELECT * FROM events WHERE date >= NOW() ORDER BY date ASC LIMIT $1 OFFSET $2;