package dtos

import "github.com/google/uuid"

type AnimeListResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Episodes int       `json:"episodes"`
	Seasonal string    `json:"seasonal"`
	Year     string    `json:"year"`
	Image    string    `json:"image"`
}
type AnimeDetailCategories struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type AnimeDetailStduios struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type AnimeDataUniverse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type GetAnimeByIdResponse struct {
	ID           uuid.UUID               `json:"id"`
	Name         string                  `json:"name"`
	NameEnglish  string                  `json:"name_english"`
	Episodes     int                     `json:"episodes"`
	Seasonal     string                  `json:"seasonal"`
	Year         string                  `json:"year"`
	Image        string                  `json:"image"`
	Description  string                  `json:"description"`
	Duration     string                  `json:"duration"`
	Type         int                     `json:"type"`
	Categories   []AnimeDetailCategories `json:"categories"`
	Wallpaper    string                  `json:"wallpaper"`
	Trailer      string                  `json:"trailer"`
	TrailerEmbed string                  `json:"trailer_embed"`
	Studios      []AnimeDetailStduios    `json:"studios"`
	Universe     AnimeDataUniverse       `json:"universe"`
}

type AnimeQueryDTO struct {
	Seasonal string
	Year     string
}

type EditCategoryToAnimeRequest struct {
	AnimeID    uuid.UUID   `json:"anime_id"`
	CategoryID []uuid.UUID `json:"category_ids"`
}

type CreateAnimeRequest struct {
	Name        string   `json:"name"`
	NameEnglish string   `json:"name_english"`
	NameThai    string   `json:"name_thai"`
	Episodes    int      `json:"episodes"`
	Seasonal    string   `json:"seasonal"`
	Image       string   `json:"image"`
	Studio      []string `json:"studio"`
	Description string   `json:"description"`
	Duration    string   `json:"duration"`
	Year        string   `json:"year"`
	Type        int      `json:"type"`
	Wallpaper   string   `json:"wallpaper"`
	Trailer     string   `json:"trailer"`
}

type GetAnimeBySeasonAndYearRequest struct {
	Year     string `json:"year"`
	Seasonal string `json:"seasonal"`
}

type GetAnimeBySeasonAndYearResponseAnime struct {
	ID           uuid.UUID               `json:"id"`
	Name         string                  `json:"name"`
	NameEnglish  string                  `json:"name_english"`
	Episodes     int                     `json:"episodes"`
	Seasonal     string                  `json:"seasonal"`
	Year         string                  `json:"year"`
	Image        string                  `json:"image"`
	Description  string                  `json:"description"`
	Duration     string                  `json:"duration"`
	Type         int                     `json:"type"`
	Categories   []AnimeDetailCategories `json:"categories"`
	Wallpaper    string                  `json:"wallpaper"`
	Trailer      string                  `json:"trailer"`
	TrailerEmbed string                  `json:"trailer_embed"`
	Studios      []AnimeDetailStduios    `json:"studios"`
	Universe     AnimeDataUniverse       `json:"universe"`
}
type GetAnimeBySeasonAndYearResponse struct {
	Year      string                                 `json:"year"`
	Seasonal  string                                 `json:"seasonal"`
	AnimeList []GetAnimeBySeasonAndYearResponseAnime `json:"anime_list"`
}
