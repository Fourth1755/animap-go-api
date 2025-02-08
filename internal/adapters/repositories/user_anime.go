package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type UserAnimeRepository interface {
	Save(userAnime *entities.UserAnime) error
	GetByUserId(id uint) ([]entities.UserAnime, error)
	GetByUserIdAndAnimeId(id uint, animeIds []uint) ([]entities.UserAnime, error)
	GetMyTopAnimeByUserId(id uint) ([]entities.UserAnime, error)
	UpdateMyTopAnime(userAnime *entities.UserAnime) error
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

func (r *GormUserAnimeRepository) GetByUserIdAndAnimeId(id uint, animeIds []uint) ([]entities.UserAnime, error) {
	var animes []entities.UserAnime

	result := r.db.Preload("Anime").Where("user_id = ?", id).Where("anime_id in (?)", animeIds).Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}

func (r *GormUserAnimeRepository) GetMyTopAnimeByUserId(id uint) ([]entities.UserAnime, error) {
	var animes []entities.UserAnime
	result := r.db.Limit(10).Preload("Anime").
		Where("user_id = ?", id).
		Where("sequence_my_top_anime <> 0").
		Order("sequence_my_top_anime ASC").
		Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}

func (r *GormUserAnimeRepository) UpdateMyTopAnime(userAnime *entities.UserAnime) error {
	result := r.db.Model(&userAnime).Updates(userAnime)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
