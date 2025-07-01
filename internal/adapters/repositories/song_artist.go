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
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormSongArtistRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) SongArtistRepository {
	return &GormSongArtistRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormSongArtistRepository) Save(songArtist []entities.SongArtist) error {
	if result := r.dbPrimary.Create(&songArtist); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormSongArtistRepository) GetByArtistId(id uuid.UUID) ([]entities.SongArtist, error) {
	var songArtists []entities.SongArtist
	result := r.dbReplica.
		Preload("Song").Where("artist_id = ?", id).Find(&songArtists)
	if result.Error != nil {
		return nil, result.Error
	}
	return songArtists, nil
}
