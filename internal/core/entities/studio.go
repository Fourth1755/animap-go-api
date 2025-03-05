package entities

import (
	"time"

	"github.com/google/uuid"
)

type Studio struct {
	ID          uuid.UUID `gorm:"primarykey"`
	Name        string    `json:"name"`
	Image       string
	Wallpaper   string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time `gorm:"index"`
}
