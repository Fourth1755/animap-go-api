package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongArtistRepository interface {
	Save([]entities.SongArtist) error
	GetByArtistId(id uuid.UUID) ([]entities.SongArtist, error)
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

func (r *GormSongArtistRepository) GetByArtistId(id uuid.UUID) ([]entities.SongArtist, error) {
	var songArtists []entities.SongArtist
	result := r.db.
		Preload("Song").Where("artist_id = ?", id).Find(&songArtists)
	if result.Error != nil {
		return nil, result.Error
	}
	return songArtists, nil
}
