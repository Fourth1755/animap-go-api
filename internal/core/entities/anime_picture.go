package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimePicture struct {
	gorm.Model
	ID         uuid.UUID `gorm:"primarykey"`
	AnimeID    uuid.UUID `json:"anime_id"`
	Anime      Anime
	PictureURL string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
