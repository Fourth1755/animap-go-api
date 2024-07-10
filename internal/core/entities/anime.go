package entities

import "gorm.io/gorm"

type Anime struct {
	gorm.Model
	Name        string `json:"name"`
	Episodes    int    `json:"episodes"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Seasonal    string `json:"seasonal"`
	Year        string `json:"year"`
}
