package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongArtist struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	SongId    uuid.UUID `json:"song_id"`
	ArtistId  uuid.UUID `json:"artist_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
