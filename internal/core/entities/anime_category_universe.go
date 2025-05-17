package entities

import (
	"time"

	"github.com/google/uuid"
)

type AnimeCategoryUniverse struct {
	ID                 uuid.UUID `gorm:"primarykey"`
	CategoryUniverseID uuid.UUID `json:"category_universe_id"`
	CategoryUniverse   CategoryUniverse
	AnimeID            uuid.UUID `json:"anime_id"`
	Anime              Anime
	Description        string `json:"description"`
	Sequence           int    `json:"sequence"`
	SequenceTimeLine   int    `json:"sequence_time_line"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time `gorm:"index"`
}
