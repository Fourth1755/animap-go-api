package dtos

type AnimeListResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Episodes int    `json:"episodes"`
	Seasonal string `json:"seasonal"`
	Year     string `json:"year"`
	Image    string `json:"image"`
}
type AnimeDetailCategories struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GetAnimeByIdResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	NameEnglish string                  `json:"name_english"`
	Episodes    int                     `json:"episodes"`
	Seasonal    string                  `json:"seasonal"`
	Year        string                  `json:"year"`
	Image       string                  `json:"image"`
	Description string                  `json:"description"`
	Duration    string                  `json:"duration"`
	Type        int                     `json:"type"`
	Categories  []AnimeDetailCategories `json:"categories"`
	Wallpaper   string                  `json:"wallpaper"`
	Trailer     string                  `json:"trailer"`
}

type AnimeQueryDTO struct {
	Seasonal string
	Year     string
}

type AddCategoryToAnimeRequest struct {
	AnimeID    uint   `json:"anime_id"`
	CategoryID []uint `json:"category_id"`
}
