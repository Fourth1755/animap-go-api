package entities

import (
	"time"

	"gorm.io/gorm"
)

type UserAnime struct {
	gorm.Model
	UserID             uint `json:"user_id"`
	User               User
	AnimeID            uint `json:"anime_id"`
	Anime              Anime
	Score              float32   `json:"score"`
	Status             int       `json:"status"`
	WatchedAt          time.Time `json:"watched_at"`
	WatchedYearAt      string    `json:"watched_year_at"`
	SequenceMyTopAnime int       `gorm:"default:0" json:"sequence_my_top_anime"`
}
