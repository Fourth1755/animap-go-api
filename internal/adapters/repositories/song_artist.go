package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type SongArtistRepository interface {
	Save([]entities.SongArtist) error
}

type GormSongArtistRepository struct {
	db *gorm.DB
}

func NewGormSongArtistRepository(db *gorm.DB) SongArtistRepository {
	return &GormSongArtistRepository{db: db}
}

func (r *GormSongArtistRepository) Save(songArtist []entities.SongArtist) error {
	if result := r.db.Create(&songArtist); result.Error != nil {
		return result.Error
	}
	return nil
}
