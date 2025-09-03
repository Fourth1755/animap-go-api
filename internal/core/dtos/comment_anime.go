package dtos

import (
	"github.com/google/uuid"
)

type CreateCommentAnimeRequest struct {
	Message string    `json:"message" binding:"required"`
	AnimeID uuid.UUID `json:"anime_id" binding:"required"`
}

type CommentAnimeAuthorResponse struct {
}

type CommentAnimeResponse struct {
	ID          uuid.UUID `json:"id"`
	Message     string    `json:"message"`
	Type        string    `json:"type"`
	CreatedAt   string    `json:"created_at"`
	AuthorID    uuid.UUID `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	AuthorImage *string   `json:"author_image"`
}

type CommentAnimePaginatedResponse struct {
	Data       []CommentAnimeResponse `json:"data"`
	Pagination PaginatedResponse      `json:"pagination"`
}
