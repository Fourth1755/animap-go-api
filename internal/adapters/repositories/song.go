package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongRepository interface {
	Save(*entities.Song) (uuid.UUID, error)
	GetById(uuid.UUID) (*entities.Song, error)
	GetByIds(ids []uuid.UUID) ([]entities.Song, error)
	GetAll() ([]entities.Song, error)
	Update(*entities.Song) error
	Delete(uuid.UUID) error
	GetByAnimeId(uuid.UUID) ([]entities.Song, error)
}

type GormSongRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormSongRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) SongRepository {
	return &GormSongRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormSongRepository) Save(song *entities.Song) (uuid.UUID, error) {
	result := r.dbPrimary.Create(&song)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return song.ID, nil
}

func (r *GormSongRepository) GetById(id uuid.UUID) (*entities.Song, error) {
	song := new(entities.Song)
	if result := r.dbReplica.Preload("Artist").Preload("SongChannel").First(&song, id); result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (r *GormSongRepository) GetByIds(ids []uuid.UUID) ([]entities.Song, error) {
	var songs []entities.Song
	if result := r.dbReplica.
		Preload("Anime").
		Preload("SongChannel").
		Where("id in (?)", ids).
		Find(&songs); result.Error != nil {
		return nil, result.Error
	}
	return songs, nil
}

func (r *GormSongRepository) GetAll() ([]entities.Song, error) {
	var song []entities.Song
	if result := r.dbReplica.Preload("Anime").Find(&song); result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (r *GormSongRepository) Update(song *entities.Song) error {
	result := r.dbPrimary.Model(&song).Updates(song)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormSongRepository) Delete(id uuid.UUID) error {
	song := new(entities.Song)
	result := r.dbPrimary.Delete(&song, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormSongRepository) GetByAnimeId(id uuid.UUID) ([]entities.Song, error) {
	var songs []entities.Song
	result := r.dbReplica.Preload("Artist").Preload("SongChannel").Where("anime_id = ?", id).Find(&songs)
	if result.Error != nil {
		return nil, result.Error
	}
	return songs, nil
}
