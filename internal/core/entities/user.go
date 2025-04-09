package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID         string `gorm:"unique"`
	Email        string `gorm:"unique"`
	Password     string
	Name         string
	SID          string
	Animes       []Anime `gorm:"many2many:user_animes;"`
	ProfileImage string
	Description  string `json:"description"`
}
