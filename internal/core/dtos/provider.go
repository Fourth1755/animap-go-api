package dtos

import "github.com/google/uuid"

type CreateProviderRequest struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image"`
}

type AddProviderToAnimeRequest struct {
	AnimeID     uuid.UUID   `json:"anime_id"`
	ProviderIDs []uuid.UUID `json:"provider_ids"`
}

type ProviderResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image string    `json:"image"`
}
