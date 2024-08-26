package entities

import "gorm.io/gorm"

type SongChannel struct {
	gorm.Model
	Channel uint   `json:"channel"` // 1: youtube 2: spotify
	Type    uint   `json:"type"`    // 1: tv_size 2: full 3: official 4 unofficial
	Link    string `json:"link"`
	SongID  uint   `json:"song_id"`
	IsMain  uint   `json:"is_main"` // 1: main 2:not main is_main for show
}
