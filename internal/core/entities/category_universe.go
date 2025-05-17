package entities

import (
	"time"

	"github.com/google/uuid"
)

type CategoryUniverse struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	NameTh    string    `json:"name_th"`
	Image     string    `json:"image"`
	Animes    []Anime   `gorm:"many2many:anime_category_universes;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}
