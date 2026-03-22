package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeProvider struct {
	ID         uuid.UUID      `gorm:"primaryKey"`
	AnimeID    uuid.UUID      `json:"anime_id"`
	Anime      Anime
	ProviderID uuid.UUID      `json:"provider_id"`
	Provider   Provider
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
