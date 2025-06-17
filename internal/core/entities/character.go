package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Character struct {
	ID              uuid.UUID `gorm:"primaryKey" json:"id"`
	Name            string    `json:"name"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	NameThai        string    `json:"name_thai"`
	FirstNameThai   string    `json:"first_name_thai"`
	LastNameThai    string    `json:"last_name_thai"`
	FirstNameJapan  string    `json:"first_name_japan"`
	LastNameJapan   string    `json:"last_name_japan"`
	Image           string    `json:"image"`
	ImageStyleX     uint      `json:"image_style_x"`
	ImageStyleY     uint      `json:"image_style_y"`
	IsMainCharacter bool      `json:"is_main_character"`
	Description     string    `json:"description"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
