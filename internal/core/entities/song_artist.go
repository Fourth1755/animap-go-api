package entities

import (
	"time"

	"github.com/google/uuid"
)

type SongArtist struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	SongId    uuid.UUID `json:"song_id"`
	ArtistId  uuid.UUID `json:"artist_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}
