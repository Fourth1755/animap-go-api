package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Artist struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Name           string         `json:"name"`
	NameJapan      string         `json:"name_japan"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	FirstNameJapan string         `json:"first_name_japan"`
	LastNameJapan  string         `json:"last_name_japan"`
	Image          string         `json:"image"`
	Description    string         `json:"description"`
	RaceOne        string         `json:"race _one"`
	RaceTwo        string         `json:"race _two"`
	RecordLabel    string         `json:"record_label"`
	IsMusicBand    bool           `json:"is_music_band"`
	DateOfBirth    time.Time      `json:"date_of_birth"`
	Member         pq.StringArray `gorm:"type:text[]" json:"member"`
	Song           []Song         `gorm:"many2many:song_artists;"`
}
