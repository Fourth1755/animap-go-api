package entities

import "gorm.io/gorm"

type Anime struct {
	gorm.Model
	ID          int
	Name        string
	Episodes    int
	Image       string
	Description string
	Seasonal    string
	Year        string
}
