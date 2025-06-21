package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	ID          uuid.UUID     `gorm:"primarykey"`
	Name        string        `json:"name"`
	Image       string        `json:"image"`
	Description string        `json:"description"`
	Year        string        `json:"year"`
	Type        string        `json:"type"` //OPENING, ENDING, SOUNDSTRACK
	Sequence    int           `json:"sequence"`
	AnimeID     uuid.UUID     `json:"anime_id"`
	Anime       Anime         `json:"anime"`
	SongChannel []SongChannel `json:"song_channel"`
	Artist      []Artist      `gorm:"many2many:song_artists;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
