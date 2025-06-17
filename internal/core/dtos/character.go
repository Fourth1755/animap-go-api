package dtos

import "github.com/google/uuid"

type CreateCharacterRequest struct {
	AnimeId         uuid.UUID `json:"anime_id"`
	Name            string    `json:"name"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	NameThai        string    `json:"name_thai"`
	FirstNameThai   string    `json:"first_name_thai"`
	LastNameThai    string    `json:"last_name_thai"`
	FirstNameJapan  string    `json:"first_name_japan"`
	LastNameJapan   string    `json:"last_name_japan"`
	Image           string    `json:"image"`
	ImageStyleX     uint      `json:"image_style_x"`
	ImageStyleY     uint      `json:"image_style_y"`
	IsMainCharacter bool      `json:"is_main_character"`
	Description     string    `json:"description"`
}

type GetCharacterByAnimeIdResponseCharacter struct {
	CharacterID     uuid.UUID `json:"id"`
	Number          uint      `json:"number"`
	Name            string    `json:"name"`
	NameThai        string    `json:"name_thai"`
	FullName        string    `json:"full_name"`
	FullNameThai    string    `json:"full_name_thai"`
	FullNameJapan   string    `json:"full_name_japan"`
	Image           string    `json:"image"`
	ImageStyleX     uint      `json:"image_style_x"`
	ImageStyleY     uint      `json:"image_style_y"`
	IsMainCharacter bool      `json:"is_main_character"`
	Description     string    `json:"description"`
}

type GetCharacterByAnimeIdResponse struct {
	Character []GetCharacterByAnimeIdResponseCharacter `json:"character"`
}
