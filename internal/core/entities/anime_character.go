package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCharacter struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	AnimeID     uuid.UUID `json:"anime_id"`
	Anime       Anime
	CharacterID uuid.UUID `json:"character_id"`
	Character   Character
	Number      uint   `json:"number"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
