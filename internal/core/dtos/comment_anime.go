package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateCommentAnimeRequest struct {
	Message string    `json:"message" binding:"required"`
	AnimeID uuid.UUID `json:"anime_id" binding:"required"`
}

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
