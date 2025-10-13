package dtos

import "github.com/google/uuid"

type GetTierListResponse struct {
}

type GetTierTemplateResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	PlayedCount uint                   `json:"played_count"`
	TierList    map[string]interface{} `json:"tier_list"`
	TotalItem   uint                   `json:"total_item"`
	IsPlay      bool                   `json:"is_play"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   string                 `json:"created_at"`
}

type CreateTierTemplateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type GetTierTemplatePaginatedResponse struct {
	Data []GetTierTemplateResponse `json:"data"`
	// Pagination PaginatedResponse      `json:"pagination"`
}

type GetByIdTierTemplateResponseItem struct {
	ID    uuid.UUID `json:"id"`
	Image string    `json:"image"`
	Name  string    `json:"name"`
}

type GetByIdTierTemplateResponse struct {
	ID          uuid.UUID                         `json:"id"`
	Name        string                            `json:"name"`
	Type        string                            `json:"type"`
	PlayedCount uint                              `json:"played_count"`
	TierList    map[string]interface{}            `json:"tier_list"`
	TotalItem   uint                              `json:"total_item"`
	Items       []GetByIdTierTemplateResponseItem `json:"items"`
}
