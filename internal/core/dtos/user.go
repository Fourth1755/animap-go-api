package dtos

import "github.com/google/uuid"

type LoginResponse struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

type GetUserInfoRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type GetUserInfoResponse struct {
	ID           uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	ProfileImage string    `json:"profile_image"`
	Description  string    `json:"description"`
}

type UpdateUserInfoRequest struct {
	ID           uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	ProfileImage string    `json:"profile_image"`
	Description  string    `json:"description"`
}
