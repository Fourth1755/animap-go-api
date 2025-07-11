package dtos

import (
	"time"

	"github.com/google/uuid"
)

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
	ID               uuid.UUID               `json:"id"`
	Name             string                  `json:"name"`
	NameEnglish      string                  `json:"name_english"`
	Episodes         int                     `json:"episodes"`
	Seasonal         string                  `json:"seasonal"`
	Year             string                  `json:"year"`
	Image            string                  `json:"image"`
	Description      string                  `json:"description"`
	Duration         string                  `json:"duration"`
	Type             int                     `json:"type"`
	Categories       []AnimeDetailCategories `json:"categories"`
	Wallpaper        string                  `json:"wallpaper"`
	Trailer          string                  `json:"trailer"`
	TrailerEmbed     string                  `json:"trailer_embed"`
	Studios          []AnimeDetailStduios    `json:"studios"`
	CategoryUniverse []AnimeDataUniverse     `json:"categoryUniverse"`
}

type AnimeQueryDTO struct {
	Seasonal   string `form:"seasonal"`
	Year       string `form:"year"`
	Name       string `form:"name"`
	SortBy     string `form:"sort_by"`
	OrderBy    string `form:"order_by"`
	StudioID   string `form:"studio_id"`
	CategoryID string `form:"category_id"`
}

type EditCategoryToAnimeRequest struct {
	AnimeID    uuid.UUID   `json:"anime_id"`
	CategoryID []uuid.UUID `json:"category_ids"`
}

type EditCategoryUniverseToAnimeRequest struct {
	AnimeID            uuid.UUID   `json:"anime_id"`
	CategoryUniverseID []uuid.UUID `json:"category_universe_ids"`
}

type CreateAnimeRequest struct {
	Name        string    `json:"name"`
	NameEnglish string    `json:"name_english"`
	NameThai    string    `json:"name_thai"`
	Episodes    int       `json:"episodes"`
	Seasonal    string    `json:"seasonal"`
	Image       string    `json:"image"`
	Studio      []string  `json:"studio"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	Year        string    `json:"year"`
	Type        int       `json:"type"`
	Wallpaper   string    `json:"wallpaper"`
	Trailer     string    `json:"trailer"`
	AiredAt     time.Time `json:"aired_at"`
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
type GetAnimeByCategoryIdResponseAnimeList struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Episodes int       `json:"episodes"`
	Seasonal string    `json:"seasonal"`
	Year     string    `json:"year"`
	Image    string    `json:"image"`
}

type GetAnimeByCategoryIdResponse struct {
	ID         uuid.UUID                               `json:"id"`
	Name       string                                  `json:"name"`
	IsUniverse string                                  `json:"is_universe"`
	Wallpaper  string                                  `json:"wallpaper"`
	AnimeList  []GetAnimeByCategoryIdResponseAnimeList `json:"anime_list"`
}

type GetAnimeByCategoryUniverseIdResponseAnimeList struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Episodes    int                     `json:"episodes"`
	Seasonal    string                  `json:"seasonal"`
	Year        string                  `json:"year"`
	Image       string                  `json:"image"`
	Description string                  `json:"description"`
	Duration    string                  `json:"duration"`
	Type        int                     `json:"type"`
	Categories  []AnimeDetailCategories `json:"categories"`
	Wallpaper   string                  `json:"wallpaper"`
	Studios     []AnimeDetailStduios    `json:"studios"`
}

type GetAnimeByCategoryUniverseIdResponse struct {
	ID        uuid.UUID                                       `json:"id"`
	Name      string                                          `json:"name"`
	Wallpaper string                                          `json:"wallpaper"`
	AnimeList []GetAnimeByCategoryUniverseIdResponseAnimeList `json:"anime_list"`
}

type GetAnimeByStudioRequest struct {
	StudioId uuid.UUID `json:"studio_id"`
}

type GetAnimeByStudioResponseAnimeList struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Episodes    int                     `json:"episodes"`
	Seasonal    string                  `json:"seasonal"`
	Year        string                  `json:"year"`
	Image       string                  `json:"image"`
	Description string                  `json:"description"`
	Duration    string                  `json:"duration"`
	Type        int                     `json:"type"`
	Categories  []AnimeDetailCategories `json:"categories"`
	Wallpaper   string                  `json:"wallpaper"`
	Studios     []AnimeDetailStduios    `json:"studios"`
}

type GetAnimeByStudioResponse struct {
	ID        uuid.UUID                           `json:"id"`
	Name      string                              `json:"name"`
	Wallpaper string                              `json:"wallpaper"`
	MainColor string                              `json:"main_color"`
	AnimeList []GetAnimeByStudioResponseAnimeList `json:"anime_list"`
}
