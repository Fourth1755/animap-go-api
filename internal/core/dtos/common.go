package dtos

type GetSeasonalAndYearResponseData struct {
	Year     string `json:"year"`
	Seasonal string `json:"seasonal"`
}

type GetSeasonalAndYearResponse struct {
	Data []GetSeasonalAndYearResponseData `json:"data"`
}
