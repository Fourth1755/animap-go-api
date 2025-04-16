package entities

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID          uuid.UUID     `gorm:"primarykey"`
	Name        string        `json:"name"`
	Image       string        `json:"image"`
	Description string        `json:"description"`
	Year        string        `json:"year"`
	Type        int           `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int           `json:"sequence"`
	AnimeID     uuid.UUID     `json:"anime_id"`
	Anime       Anime         `json:"anime"`
	SongChannel []SongChannel `json:"song_channel"`
	Artist      []Artist      `gorm:"many2many:song_artists;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time `gorm:"index"`
}
