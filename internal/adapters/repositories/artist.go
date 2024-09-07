package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type ArtistRepository interface {
	Save(*entities.Artist) error
	GetAll() ([]entities.Artist, error)
	GetById(uint) (*entities.Artist, error)
	GetByIds([]uint) ([]entities.Artist, error)
}

type GormArtistRepository struct {
	db *gorm.DB
}

func NewGormArtistRepository(db *gorm.DB) ArtistRepository {
	return &GormArtistRepository{db: db}
}

func (r GormArtistRepository) Save(artist *entities.Artist) error {
	if result := r.db.Create(&artist); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormArtistRepository) GetAll() ([]entities.Artist, error) {
	var artist []entities.Artist
	if result := r.db.Find(&artist); result.Error != nil {
		return nil, result.Error
	}
	return artist, nil
}

func (r GormArtistRepository) GetById(id uint) (*entities.Artist, error) {
	var artist *entities.Artist
	if result := r.db.Preload("Song").First(&artist, id); result.Error != nil {
		return nil, result.Error
	}
	return artist, nil
}

func (r GormArtistRepository) GetByIds(ids []uint) ([]entities.Artist, error) {
	var artist []entities.Artist
	if result := r.db.Find(&artist, ids); result.Error != nil {
		return nil, result.Error
	}
	return artist, nil
}
