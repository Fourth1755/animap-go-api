package entities

import (
	"time"

	"github.com/google/uuid"
)

type AnimeStudio struct {
	ID        uuid.UUID `gorm:"primarykey"`
	StudioId  uuid.UUID `json:"studio_id"`
	Studio    Studio
	AnimeID   uuid.UUID `json:"anime_id"`
	Anime     Anime
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}
