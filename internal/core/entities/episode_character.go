package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeCharacter struct {
	ID              uuid.UUID `gorm:"primaryKey" json:"id"`
	EpisodeID       uuid.UUID `json:"episode_id"`
	Episode         Episode
	CharacterID     uuid.UUID `json:"character_id"`
	Character       Character
	Description     string `json:"description"`
	FirstAppearance bool   `json:"firstAppearance"`
	Appearance      bool   `json:"appearance"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
