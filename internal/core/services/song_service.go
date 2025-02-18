package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type SongService interface {
	CreateSong(*dtos.CreateSongRequest) error
	GetSongById(uint) (*entities.Song, error)
	GetAllSongs() ([]dtos.SongListResponse, error)
	UpdateSong(*entities.Song) error
	DeleteSong(uint) error
	GetSongByAnimeId(uint) (*dtos.SongDetailResponse, error)
	CreateSongChannel(request *dtos.CreateSongChannelRequest) error
}

type songServiceImpl struct {
	repo            repositories.SongRepository
	animeRepo       repositories.AnimeRepository
	artistRepo      repositories.ArtistRepository
	songArtistRepo  repositories.SongArtistRepository
	songChannelRepo repositories.SongChannelRepository
}

func NewSongService(
	repo repositories.SongRepository,
	animeRepo repositories.AnimeRepository,
	artistRepo repositories.ArtistRepository,
	songArtistRepo repositories.SongArtistRepository,
	songChannelRepo repositories.SongChannelRepository) SongService {
	return &songServiceImpl{
		repo:            repo,
		animeRepo:       animeRepo,
		artistRepo:      artistRepo,
		songArtistRepo:  songArtistRepo,
		songChannelRepo: songChannelRepo}
}

var songTypeMap = []string{"none", "opening", "ending", "soundtrack"}

func (s songServiceImpl) CreateSong(songRequest *dtos.CreateSongRequest) error {
	//validate anime id
	if _, err := s.animeRepo.GetById(songRequest.AnimeID); err != nil {
		logs.Error(err)
		return errs.NewNotFoundError("Anime Not Found")
	}

	//validate artist id
	artistList, err := s.artistRepo.GetByIds(songRequest.ArtistList)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	songChannel := []entities.SongChannel{}
	for _, item := range songRequest.SongChannel {
		songChannel = append(songChannel, entities.SongChannel{
			Channel: item.Channel,
			Type:    item.Type,
			Link:    item.Link,
			IsMain:  item.IsMain,
		})
	}

	song := entities.Song{
		Name:        songRequest.Name,
		Image:       songRequest.Image,
		Description: songRequest.Description,
		Year:        songRequest.Year,
		Type:        songRequest.Type,
		Sequence:    songRequest.Sequence,
		AnimeID:     songRequest.AnimeID,
		SongChannel: songChannel,
	}

	//save song
	songId, err := s.repo.Save(&song)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	//save song artist
	songArtist := []entities.SongArtist{}
	for _, item := range artistList {
		songArtist = append(songArtist, entities.SongArtist{
			SongId:   songId,
			ArtistId: item.ID,
		})
	}
	if err := s.songArtistRepo.Save(songArtist); err != nil {
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

func (s songServiceImpl) CreateSongChannel(request *dtos.CreateSongChannelRequest) error {
	songChannel := entities.SongChannel{
		Channel: request.Channel,
		Type:    request.Type,
		Link:    request.Link,
		SongID:  request.SongID,
		IsMain:  request.IsMain,
	}

	//save song
	err := s.songChannelRepo.Save(&songChannel)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	return nil
}
