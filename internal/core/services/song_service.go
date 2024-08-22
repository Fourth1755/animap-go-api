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
}

type ImplSongService struct {
	repo repositories.SongRepository
}

func NewSongService(repo repositories.SongRepository) SongService {
	return &ImplSongService{repo: repo}
}

func (s ImplSongService) CreateSong(song *entities.Song) error {
	if err := s.repo.Save(song); err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s ImplSongService) GetSongById(id uint) (*entities.Song, error) {
	song, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Song not found")
	}
	return song, nil
}

func (s ImplSongService) GetAllSongs() ([]entities.Song, error) {
	songs, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	return songs, nil
}
