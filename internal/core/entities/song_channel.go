package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongChannel struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Channel   int       `json:"channel"` // 1: youtube 2: spotify
	Type      int       `json:"type"`    // 1: tv_size 2: full 3: official 4 unofficial
	Link      string    `json:"link"`
	SongID    uuid.UUID `json:"song_id"`
	IsMain    int       `json:"is_main"` // 1: main 2:not main is_main for show
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
