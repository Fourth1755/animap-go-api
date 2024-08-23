package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type AnimeRepository interface {
	Save(anime entities.Anime) error
	GetById(id uint) (*entities.Anime, error)
	GetAll(query dtos.AnimeQueryDTO) ([]entities.Anime, error)
	Update(anime *entities.Anime) error
	Delete(id uint) error
	GetByUserId(user_id uint) ([]entities.UserAnime, error)
}

type GormAnimeRepository struct {
	db *gorm.DB
}

func NewGormAnimeRepository(db *gorm.DB) AnimeRepository {
	return &GormAnimeRepository{db: db}
}

func (r *GormAnimeRepository) Save(anime entities.Anime) error {
	if result := r.db.Create(&anime); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) GetById(id uint) (*entities.Anime, error) {
	var anime entities.Anime
	if result := r.db.Preload("Songs").Preload("Categories").First(&anime, id); result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (r *GormAnimeRepository) GetAll(query dtos.AnimeQueryDTO) ([]entities.Anime, error) {
	var animes []entities.Anime
	var result *gorm.DB
	if query.Seasonal == "" && query.Year == "" {
		result = r.db.Find(&animes)
	} else {
		result = r.db.Where("seasonal = ?", query.Seasonal).Where("year = ?", query.Year).Find(&animes)
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return animes, nil
}

func (r *GormAnimeRepository) Update(anime *entities.Anime) error {
	result := r.db.Model(&anime).Updates(anime)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) Delete(id uint) error {
	var anime entities.Anime
	result := r.db.Delete(&anime, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) GetByUserId(id uint) ([]entities.UserAnime, error) {
	var animes []entities.UserAnime
	result := r.db.Preload("Anime").Where("user_id = ?", id).Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}
