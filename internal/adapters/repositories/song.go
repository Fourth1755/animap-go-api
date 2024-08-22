package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type SongRepository interface {
	Save(*entities.Song) error
	GetById(uint) (*entities.Song, error)
	GetAll() ([]entities.Song, error)
}

type GormSongRepository struct {
	db *gorm.DB
}

func NewGormSongRepository(db *gorm.DB) SongRepository {
	return &GormSongRepository{db: db}
}

func (r *GormSongRepository) Save(song *entities.Song) error {
	if result := r.db.Create(&song); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormSongRepository) GetById(id uint) (*entities.Song, error) {
	song := new(entities.Song)
	if result := r.db.First(&song, id); result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (r *GormSongRepository) GetAll() ([]entities.Song, error) {
	var song []entities.Song
	if result := r.db.Find(&song); result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}
