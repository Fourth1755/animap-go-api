package dtos

type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}
