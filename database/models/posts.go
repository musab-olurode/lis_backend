package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/musab-olurode/lis_backend/database"
)

type Post struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Description   *string   `json:"description"`
	Content       string    `json:"content"`
	CoverImageUrl string    `json:"cover_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func DatabasePostToPost(post database.Post) Post {
	return Post{
		ID:            post.ID,
		Title:         post.Title,
		Slug:          post.Slug,
		Description:   nullStringToStringPtr(post.Description),
		Content:       post.Content,
		CoverImageUrl: post.CoverImageUrl,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}
}

func DatabasePostsToPosts(posts []database.Post) []Post {
	result := make([]Post, len(posts))
	for i, post := range posts {
		result[i] = DatabasePostToPost(post)
	}
	return result
}
