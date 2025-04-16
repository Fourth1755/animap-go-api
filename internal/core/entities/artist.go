package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Artist struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `json:"name"`
	Image       string         `json:"image"`
	Description string         `json:"description"`
	RecordLabel string         `json:"record_label"`
	IsMusicBand bool           `json:"is_music_band"`
	Member      pq.StringArray `gorm:"type:string[]" json:"member"`
	Song        []Song         `gorm:"many2many:song_artists;"`
}
