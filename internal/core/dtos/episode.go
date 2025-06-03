package dtos

import "github.com/google/uuid"

type CreateEpisodeRequest struct {
	AnimeId uuid.UUID `json:"anime_id"`
}
