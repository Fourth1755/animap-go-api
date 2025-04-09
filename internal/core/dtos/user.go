package dtos

type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

type GetUserInfoRequest struct {
	UUID string `json:"uuid"`
}

type GetUserInfoResponse struct {
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	Email        string `json:"emdil"`
	ProfileImage string `json:"profile_image"`
	Description  string `json:"description"`
}

type UpdateUserInfoRequest struct {
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	Email        string `json:"emdil"`
	ProfileImage string `json:"profile_image"`
	Description  string `json:"description"`
}
