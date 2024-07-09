package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
	"gorm.io/gorm"
)

type GormAnimeRepository struct {
	db *gorm.DB
}

func NewGormAnimeRepository(db *gorm.DB) ports.AnimeRepository {
	return &GormAnimeRepository{db: db}
}

func (r *GormAnimeRepository) Save(anime entities.Anime) error {
	if result := r.db.Create(&anime); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) GetById(id int) (*entities.Anime, error) {
	var anime entities.Anime
	if result := r.db.First(&anime, id); result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (r *GormAnimeRepository) GetAll(query dtos.AnimeQueryDTO) ([]entities.Anime, error) {
	var animes []entities.Anime

	result := r.db.Where("seasonal = ?", query.Seasonal).Where("year = ?", query.Year).Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}

	return animes, nil
}
