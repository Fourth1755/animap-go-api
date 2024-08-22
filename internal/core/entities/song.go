package entities

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Year        string `json:"year"`
	Type        uint   `json:"type"` // 1: opening, 2: ending, 3:other
	AnimeID     uint   `json:"anime_id"`
	Anime       Anime
	SongChannel []SongChannel
}
