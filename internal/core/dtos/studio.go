package dtos

import "github.com/google/uuid"

type StudioListRequest struct {
}

type StudioListResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
}
