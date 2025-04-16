package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserAnime struct {
	ID                 uuid.UUID `gorm:"primaryKey" json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	User               User
	AnimeID            uuid.UUID `json:"anime_id"`
	Anime              Anime
	Score              float32   `json:"score"`
	Status             int       `json:"status"`
	WatchedAt          time.Time `json:"watched_at"`
	WatchedYearAt      string    `json:"watched_year_at"`
	SequenceMyTopAnime int       `gorm:"default:0" json:"sequence_my_top_anime"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time `gorm:"index"`
}
