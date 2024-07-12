package entities

import (
	"time"

	"gorm.io/gorm"
)

type UserAnime struct {
	gorm.Model
	UserID  uint `json:"user_id"`
	User    User
	AnimeID uint `json:"anime_id"`
	Anime   Anime
	Score   float32   `json:"score"`
	Status  int       `json:"status"`
	WatchAt time.Time `json:"watch_at"`
}
