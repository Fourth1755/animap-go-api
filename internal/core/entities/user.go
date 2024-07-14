package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Name     string
	Animes   []Anime `gorm:"many2many:user_animes;"`
}
