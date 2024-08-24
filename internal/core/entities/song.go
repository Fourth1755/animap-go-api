package entities

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Name        string        `json:"name"`
	Image       string        `json:"image"`
	Description string        `json:"description"`
	Year        string        `json:"year"`
	Type        int           `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int           `json:"sequence"`
	AnimeID     uint          `json:"anime_id"`
	Anime       Anime         `json:"anime"`
	SongChannel []SongChannel `json:"song_channel"`
}
