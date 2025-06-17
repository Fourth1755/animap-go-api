package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeStudio struct {
	ID        uuid.UUID `gorm:"primarykey"`
	StudioId  uuid.UUID `json:"studio_id"`
	Studio    Studio
	AnimeID   uuid.UUID `json:"anime_id"`
	Anime     Anime
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
