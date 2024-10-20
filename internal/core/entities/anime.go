package entities

import "gorm.io/gorm"

type Anime struct {
	gorm.Model
	Name        string     `json:"name"`
	NameThai    string     `json:"name_thai"`
	NameEnglish string     `json:"name_english"`
	Episodes    int        `json:"episodes"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Seasonal    string     `json:"seasonal"`
	Year        string     `json:"year"`
	Type        int        `json:"type"` //1: TV 2: movie
	Duration    string     `json:"duration"`
	Categories  []Category `gorm:"many2many:anime_categories;"`
	Songs       []Song
	Wallpaper   string `json:"wallpaper"`
	Trailer     string `json:"trailer"`
}
