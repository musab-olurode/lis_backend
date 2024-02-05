package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/musab-olurode/lis_backend/database"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	ImageUrl    string    `json:"image_url"`
	Venue       string    `json:"venue"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DatabaseEventToEvent(event database.Event) Event {
	return Event{
		ID:          event.ID,
		Title:       event.Title,
		ImageUrl:    event.ImageUrl,
		Description: nullStringToStringPtr(event.Description),
		CreatedAt:   event.CreatedAt,
		Date:        event.Date,
		UpdatedAt:   event.UpdatedAt,
	}
}

func DatabaseEventsToEvents(events []database.Event) []Event {
	result := make([]Event, len(events))
	for i, event := range events {
		result[i] = DatabaseEventToEvent(event)
	}
	return result
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
