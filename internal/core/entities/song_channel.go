package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongChannel struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Channel   string    `json:"channel"` // YOUTUBE,SPOTIFY
	Type      string    `json:"type"`    // TV_SIZE, FULL_SIZE_OFFICIAL, FULL_SIZE_UNOFFICIAL, FIRST_TAKE
	Link      string    `json:"link"`
	SongID    uuid.UUID `json:"song_id"`
	IsMain    bool      `json:"is_main"` // true: main false:not main is_main for show
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
