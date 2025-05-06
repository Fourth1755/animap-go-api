package dtos

type GetSeasonalAndYearRequest struct {
	Year     string
	Seasonal string
}

type GetSeasonalAndYearResponseData struct {
	Year     string
	Seasonal string
	IsMain   bool
}

type GetSeasonalAndYearResponse struct {
	Data []GetSeasonalAndYearResponseData
}
