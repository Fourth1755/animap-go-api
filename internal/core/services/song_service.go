package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type SongService interface {
	CreateSong(*entities.Song) error
	GetSongById(uint) (*entities.Song, error)
	GetAllSongs() ([]dtos.SongListResponse, error)
	UpdateSong(*entities.Song) error
	DeleteSong(uint) error
	GetSongByAnimeId(uint) ([]entities.Song, error)
}

type songServiceImpl struct {
	repo      repositories.SongRepository
	animeRepo repositories.AnimeRepository
}

func NewSongService(repo repositories.SongRepository, animeRepo repositories.AnimeRepository) SongService {
	return &songServiceImpl{repo: repo, animeRepo: animeRepo}
}

var songTypeMap = []string{"none", "opening", "ending", "soundtrack"}

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

func (s songServiceImpl) GetAllSongs() ([]dtos.SongListResponse, error) {
	songs, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	var songResponse []dtos.SongListResponse
	for _, song := range songs {
		songResponse = append(songResponse, dtos.SongListResponse{
			ID:          song.ID,
			Name:        song.Name,
			Image:       song.Image,
			Description: song.Description,
			Year:        song.Year,
			Type:        songTypeMap[song.Type],
			Sequence:    song.Sequence,
			AnimeID:     song.AnimeID,
			AnimeName:   song.Anime.Name,
		})
	}

	return songResponse, nil
}

func (s songServiceImpl) UpdateSong(song *entities.Song) error {
	_, err := s.repo.GetById(song.ID)
	if err != nil {
		logs.Error(err)
		return errs.NewNotFoundError("Song not found")
	}
	if err := s.repo.Update(song); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (s songServiceImpl) DeleteSong(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (s songServiceImpl) GetSongByAnimeId(animeId uint) ([]entities.Song, error) {
	_, err := s.animeRepo.GetById(animeId)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Anime not found")
	}
	songs, err := s.repo.GetByAnimeId(animeId)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return songs, nil
}
