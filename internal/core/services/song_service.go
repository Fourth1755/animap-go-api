package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type SongService interface {
	CreateSong(*entities.Song) error
	GetSongById(uint) (*entities.Song, error)
	GetAllSongs() ([]entities.Song, error)
	UpdateSong(*entities.Song) error
}

type songServiceImpl struct {
	repo repositories.SongRepository
}

func NewSongService(repo repositories.SongRepository) SongService {
	return &songServiceImpl{repo: repo}
}

func (s songServiceImpl) CreateSong(song *entities.Song) error {
	if err := s.repo.Save(song); err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s songServiceImpl) GetSongById(id uint) (*entities.Song, error) {
	song, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Song not found")
	}
	return song, nil
}

func (s songServiceImpl) GetAllSongs() ([]entities.Song, error) {
	songs, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	return songs, nil
}

func (s songServiceImpl) UpdateSong(song *entities.Song) error {
	_, err := s.repo.GetById(song.ID)
	if err != nil {
		return errs.NewNotFoundError("Song not found")
	}
	if err := s.repo.Update(song); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}
