package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Studio struct {
	ID          uuid.UUID `gorm:"primarykey"`
	Name        string    `json:"name"`
	Image       string
	Wallpaper   string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Animes      []Anime        `gorm:"many2many:anime_studios;"`
	MainColor   string         `json:"main_color"`
}
