package dtos

import "github.com/Fourth1755/animap-go-api/internal/core/entities"

type GetArtistsResponse struct {
	Artists []entities.Artist `json:"artists"`
}
