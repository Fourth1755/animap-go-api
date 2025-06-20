package dtos

import "github.com/google/uuid"

type CreateEpisodeRequest struct {
	AnimeId uuid.UUID `json:"anime_id"`
}

type GetEpisodeResponseEpisodeCharacter struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Image           string    `json:"image"`
	ImageStyleX     uint      `json:"image_style_x"`
	ImageStyleY     uint      `json:"image_style_y"`
	Description     string    `json:"description"`
	FirstAppearance bool      `json:"first_appearance"`
	Appearance      bool      `json:"appearance"`
}

type GetEpisodeResponseEpisode struct {
	ID          uuid.UUID                            `json:"id"`
	Number      uint                                 `json:"number"`
	Name        string                               `json:"name"`
	NameThai    string                               `json:"name_thai"`
	NameEnglish string                               `json:"name_english"`
	NameJapan   string                               `json:"name_japan"`
	Characters  []GetEpisodeResponseEpisodeCharacter `json:"characters"`
}

type GetEpisodeResponse struct {
	Episodes []GetEpisodeResponseEpisode `json:"episodes"`
}

type UpdateEpisodeRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	NameThai    string    `json:"name_thai"`
	NameEnglish string    `json:"name_english"`
	NameJapan   string    `json:"name_japan"`
	Image       string    `json:"image"`
}

type AddCharacterToEpisodeRequestCharacter struct {
	ID              uuid.UUID `json:"id"`
	Description     string    `json:"description"`
	FirstAppearance bool      `json:"first_appearance"`
	Appearance      bool      `json:"appearance"`
}

type AddCharacterToEpisodeRequest struct {
	EpisodeID  uuid.UUID                               `json:"episode_id"`
	Characters []AddCharacterToEpisodeRequestCharacter `json:"characters"`
}
