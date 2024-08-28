package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type SongRepository interface {
	Save(*entities.Song) (uint, error)
	GetById(uint) (*entities.Song, error)
	GetAll() ([]entities.Song, error)
	Update(*entities.Song) error
	Delete(uint) error
	GetByAnimeId(uint) ([]entities.Song, error)
}

type GormSongRepository struct {
	db *gorm.DB
}

func NewGormSongRepository(db *gorm.DB) SongRepository {
	return &GormSongRepository{db: db}
}

func (r *GormSongRepository) Save(song *entities.Song) (uint, error) {
	result := r.db.Create(&song)
	if result.Error != nil {
		return 0, result.Error
	}
	return song.ID, nil
}

func (r *GormSongRepository) GetById(id uint) (*entities.Song, error) {
	song := new(entities.Song)
	if result := r.db.Preload("SongChannel").First(&song, id); result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (r *GormSongRepository) GetAll() ([]entities.Song, error) {
	var song []entities.Song
	if result := r.db.Preload("Anime").Find(&song); result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (r *GormSongRepository) Update(song *entities.Song) error {
	result := r.db.Model(&song).Updates(song)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormSongRepository) Delete(id uint) error {
	song := new(entities.Song)
	result := r.db.Delete(&song, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormSongRepository) GetByAnimeId(id uint) ([]entities.Song, error) {
	var songs []entities.Song
	result := r.db.Preload("SongChannel").Where("anime_id = ?", id).Find(&songs)
	if result.Error != nil {
		return nil, result.Error
	}
	return songs, nil
}
