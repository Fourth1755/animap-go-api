package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeStudio struct {
	gorm.Model
	StudioId uuid.UUID `json:"studio_id"`
	Studio   Studio
	AnimeID  uint `json:"anime_id"`
	Anime    Anime
}
