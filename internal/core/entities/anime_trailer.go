package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeTrailer struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Name      string    `json:"name"`
	AnimeID   uuid.UUID `gorm:"index" json:"anime_id"`
	Anime     Anime
	VideoID   string `json:"video_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
