package entities

import (
	"time"

	"github.com/google/uuid"
)

type AnimeCategory struct {
	ID               uuid.UUID `gorm:"primarykey"`
	CategoryID       uuid.UUID `json:"category_id"`
	Category         Category
	AnimeID          uuid.UUID `json:"anime_id"`
	Anime            Anime
	Description      string `json:"description"`
	Sequence         int    `json:"sequence"`
	SequenceTimeLine int    `json:"sequence_time_line"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time `gorm:"index"`
}
