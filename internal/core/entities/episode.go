package entities

import (
	"time"

	"github.com/google/uuid"
)

type Episode struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Number      uint      `json:"number"`
	Name        string    `json:"name"`
	NameThai    string    `json:"name_thai"`
	NameEnglish string    `json:"name_english"`
	NameJapan   string    `json:"name_japan"`
	Image       string    `json:"image"`
	AnimeID     uuid.UUID `json:"anime_id"`
	Anime       Anime
	Character   []Character `gorm:"many2many:episode_characters;"`
	AiredAt     time.Time   `json:"aired_at"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time `gorm:"index"`
}
