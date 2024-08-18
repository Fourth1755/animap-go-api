package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type UserAnimeRepository interface {
	Save(userAnime *entities.UserAnime) error
	GetByUserId(id uint) ([]entities.UserAnime, error)
}

type GormUserAnimeRepository struct {
	db *gorm.DB
}

func NewGormUserAnimeRepository(db *gorm.DB) UserAnimeRepository {
	return &GormUserAnimeRepository{db: db}
}

func (r *GormUserAnimeRepository) Save(userAnime *entities.UserAnime) error {
	if err := r.db.Create(&userAnime); err != nil {
		return err.Error
	}
	return nil
}

func (r *GormUserAnimeRepository) GetByUserId(id uint) ([]entities.UserAnime, error) {
	var animes []entities.UserAnime
	result := r.db.Preload("Anime").Where("user_id = ?", id).Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}
