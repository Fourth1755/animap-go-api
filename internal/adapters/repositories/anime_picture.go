package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimePictureRepository interface {
	Save(pictures []entities.AnimePicture) error
	GetByAnimeId(animeID uuid.UUID) ([]entities.AnimePicture, error)
}

type GormAnimePictureRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimePictureRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimePictureRepository {
	return &GormAnimePictureRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormAnimePictureRepository) Save(pictures []entities.AnimePicture) error {
	if result := r.dbPrimary.Create(&pictures); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimePictureRepository) GetByAnimeId(animeID uuid.UUID) ([]entities.AnimePicture, error) {
	var pictures []entities.AnimePicture
	if err := r.dbReplica.Where("anime_id = ?", animeID).Find(&pictures).Error; err != nil {
		return nil, err
	}
	return pictures, nil
}
