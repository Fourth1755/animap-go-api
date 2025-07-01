package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArtistRepository interface {
	Save(*entities.Artist) error
	Update(artist *entities.Artist) error
	GetAll() ([]entities.Artist, error)
	GetById(uuid.UUID) (*entities.Artist, error)
	GetByIds([]uuid.UUID) ([]entities.Artist, error)
}

type GormArtistRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormArtistRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) ArtistRepository {
	return &GormArtistRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormArtistRepository) Save(artist *entities.Artist) error {
	if result := r.dbPrimary.Create(&artist); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormArtistRepository) Update(artist *entities.Artist) error {
	result := r.dbPrimary.Model(&artist).Updates(artist)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormArtistRepository) GetAll() ([]entities.Artist, error) {
	var artist []entities.Artist
	if result := r.dbReplica.Find(&artist); result.Error != nil {
		return nil, result.Error
	}
	return artist, nil
}

func (r GormArtistRepository) GetById(id uuid.UUID) (*entities.Artist, error) {
	var artist *entities.Artist
	if result := r.dbReplica.Preload("Song").First(&artist, id); result.Error != nil {
		return nil, result.Error
	}
	return artist, nil
}

func (r GormArtistRepository) GetByIds(ids []uuid.UUID) ([]entities.Artist, error) {
	var artist []entities.Artist
	if result := r.dbReplica.Find(&artist, ids); result.Error != nil {
		return nil, result.Error
	}
	return artist, nil
}
