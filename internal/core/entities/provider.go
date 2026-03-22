package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Provider struct {
	ID        uuid.UUID      `gorm:"primaryKey"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Animes    []Anime        `gorm:"many2many:anime_providers;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
