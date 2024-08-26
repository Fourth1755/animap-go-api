package entities

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	IsUniverse bool    `json:"is_universe"`
	Animes     []Anime `gorm:"many2many:anime_categories;"`
}
