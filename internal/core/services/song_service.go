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
	GetSongByAnimeId(uint) (*dtos.SongDetailResponse, error)
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

func (s songServiceImpl) GetSongByAnimeId(animeId uint) (*dtos.SongDetailResponse, error) {
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
	var openingSong []dtos.SongDetail
	var endingSong []dtos.SongDetail
	var soundtrack []dtos.SongDetail
	SongType := NewSongType()
	for _, song := range songs {
		var songChannelData []dtos.SongChannel
		for _, channel := range song.SongChannel {
			songChannel := dtos.SongChannel{
				ID:      channel.ID,
				Channel: channel.Channel,
				Type:    channel.Type,
				Link:    channel.Link,
				IsMain:  channel.IsMain,
			}
			songChannelData = append(songChannelData, songChannel)
		}

		songData := dtos.SongDetail{
			ID:          song.ID,
			Name:        song.Name,
			Type:        song.Type,
			Sequence:    song.Sequence,
			Image:       song.Image,
			Description: song.Description,
			Year:        song.Year,
			AnimeID:     song.AnimeID,
			SongChannel: songChannelData,
		}
		if song.Type == SongType.Opening {
			openingSong = append(openingSong, songData)
		} else if song.Type == SongType.Ending {
			endingSong = append(endingSong, songData)
		} else if song.Type == SongType.Soundtrack {
			soundtrack = append(soundtrack, songData)
		}
	}
	songResponse := dtos.SongDetailResponse{
		OpeningSong:    openingSong,
		EndingSong:     endingSong,
		SoundtrackSong: soundtrack,
	}
	return &songResponse, nil
}
