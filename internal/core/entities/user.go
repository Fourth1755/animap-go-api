package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"unique"`
	Password     string
	Name         string
	SID          string
	Animes       []Anime `gorm:"many2many:user_animes;"`
	ProfileImage string
	Description  string `json:"description"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time `gorm:"index"`
}
