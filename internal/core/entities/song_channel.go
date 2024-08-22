package entities

import "gorm.io/gorm"

type SongChannel struct {
	gorm.Model
	Type   uint   `json:"type"` // 1: youtube 2: spotify
	Link   string `json:"link"`
	SongID uint   `json:"song_id"`
	Song   Song
}
