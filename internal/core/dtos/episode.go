package dtos

import "github.com/google/uuid"

type CreateEpisodeRequest struct {
	AnimeId uuid.UUID `json:"anime_id"`
}

type GetEpisodeResponseEpisode struct {
	ID          uuid.UUID `json:"id"`
	Number      uint      `json:"number"`
	Name        string    `json:"name"`
	NameThai    string    `json:"name_thai"`
	NameEnglish string    `json:"name_english"`
}

type GetEpisodeResponse struct {
	Episodes []GetEpisodeResponseEpisode `json:"episodes"`
}

type UpdateEpisodeRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	NameThai    string    `json:"name_thai"`
	NameEnglish string    `json:"name_english"`
	Image       string    `json:"image"`
}

type AddCharacterToEpisodeRequestCharacter struct {
	ID              uuid.UUID `json:"id"`
	Description     string    `json:"description"`
	FirstAppearance bool      `json:"firstAppearance"`
	Appearance      bool      `json:"appearance"`
}

type AddCharacterToEpisodeRequest struct {
	EpisodeID  uuid.UUID `json:"episode_id"`
	Characters []AddCharacterToEpisodeRequestCharacter
}
