package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
	"gorm.io/gorm"
)

type GormUserAnimeRepository struct {
	db *gorm.DB
}

func NewGormUserAnimeRepository(db *gorm.DB) ports.UserAnimeRepository {
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
	result := r.db.Where("user_id = ?", id).Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}
