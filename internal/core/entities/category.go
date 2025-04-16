package entities

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID         uuid.UUID `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	IsUniverse bool      `json:"is_universe"`
	Animes     []Anime   `gorm:"many2many:anime_categories;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time `gorm:"index"`
}
