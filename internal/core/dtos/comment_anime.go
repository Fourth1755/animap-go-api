package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CommentAnimeAuthorResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image *string   `json:"image"`
}

type CommentAnimeResponse struct {
	ID        uuid.UUID                  `json:"id"`
	Content   string                     `json:"content"`
	Type      string                     `json:"type"`
	CreatedAt time.Time                  `json:"created_at"`
	Author    CommentAnimeAuthorResponse `json:"author"`
}

type CommentAnimePaginatedResponse struct {
	Data       []CommentAnimeResponse `json:"data"`
	Pagination PaginatedResponse      `json:"pagination"`
}
