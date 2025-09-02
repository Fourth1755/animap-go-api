package dtos

type GetSeasonalAndYearResponseData struct {
	Year     string `json:"year"`
	Seasonal string `json:"seasonal"`
}

type GetSeasonalAndYearResponse struct {
	Data []GetSeasonalAndYearResponseData `json:"data"`
}

type PaginatedResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
}
