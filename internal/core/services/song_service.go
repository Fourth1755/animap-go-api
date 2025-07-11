package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
)

type SongService interface {
	CreateSong(*dtos.CreateSongRequest) error
	GetSongById(uuid.UUID) (*entities.Song, error)
	GetAllSongs() ([]dtos.SongListResponse, error)
	UpdateSong(*entities.Song) error
	DeleteSong(uuid.UUID) error
	GetSongByAnimeId(uuid.UUID) (*dtos.GetSongByAnimeIdResponse, error)
	CreateSongChannel(request *dtos.CreateSongChannelRequest) error
	GetSongsByArtistId(artistId uuid.UUID) (*dtos.GetSongsByArtistResponse, error)
}

type songServiceImpl struct {
	songRepo        repositories.SongRepository
	animeRepo       repositories.AnimeRepository
	artistRepo      repositories.ArtistRepository
	songArtistRepo  repositories.SongArtistRepository
	songChannelRepo repositories.SongChannelRepository
}

func NewSongService(
	songRepo repositories.SongRepository,
	animeRepo repositories.AnimeRepository,
	artistRepo repositories.ArtistRepository,
	songArtistRepo repositories.SongArtistRepository,
	songChannelRepo repositories.SongChannelRepository) SongService {
	return &songServiceImpl{
		songRepo:        songRepo,
		animeRepo:       animeRepo,
		artistRepo:      artistRepo,
		songArtistRepo:  songArtistRepo,
		songChannelRepo: songChannelRepo}
}

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
		songChannelId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}

		linkEmbed := ""
		if item.Link != "" {
			linkEmbed, err = convertYouTubeURLToEmbed(item.Link)
			if err != nil {
				logs.Error(err.Error())
				return errs.NewUnexpectedError()
			}
		}

		songChannel = append(songChannel, entities.SongChannel{
			ID:      songChannelId,
			Channel: item.Channel,
			Type:    item.Type,
			Link:    linkEmbed,
			IsMain:  item.IsMain,
		})
	}

	songId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	song := entities.Song{
		ID:          songId,
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
	songId, err = s.songRepo.Save(&song)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	//save song artist
	songArtist := []entities.SongArtist{}
	for _, item := range artistList {
		songArtistId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		songArtist = append(songArtist, entities.SongArtist{
			ID:       songArtistId,
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

func (s songServiceImpl) GetSongById(id uuid.UUID) (*entities.Song, error) {
	song, err := s.songRepo.GetById(id)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Song not found")
	}
	return song, nil
}

func (s songServiceImpl) GetAllSongs() ([]dtos.SongListResponse, error) {
	songs, err := s.songRepo.GetAll()
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
			Type:        song.Type,
			Sequence:    song.Sequence,
			AnimeID:     song.AnimeID,
			AnimeName:   song.Anime.Name,
		})
	}

	return songResponse, nil
}

func (s songServiceImpl) UpdateSong(song *entities.Song) error {
	_, err := s.songRepo.GetById(song.ID)
	if err != nil {
		logs.Error(err)
		return errs.NewNotFoundError("Song not found")
	}
	if err := s.songRepo.Update(song); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (s songServiceImpl) DeleteSong(id uuid.UUID) error {
	err := s.songRepo.Delete(id)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (s songServiceImpl) GetSongByAnimeId(animeId uuid.UUID) (*dtos.GetSongByAnimeIdResponse, error) {
	_, err := s.animeRepo.GetById(animeId)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Anime not found")
	}
	songs, err := s.songRepo.GetByAnimeId(animeId)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	var openingSong []dtos.GetSongByAnimeIdResponseSong
	var endingSong []dtos.GetSongByAnimeIdResponseSong
	var soundtrack []dtos.GetSongByAnimeIdResponseSong

	for _, song := range songs {
		var songChannelData []dtos.GetSongByAnimeIdResponseSongChannel
		for _, channel := range song.SongChannel {
			songChannel := dtos.GetSongByAnimeIdResponseSongChannel{
				ID:      channel.ID,
				Channel: channel.Channel,
				Type:    channel.Type,
				Link:    channel.Link,
				IsMain:  channel.IsMain,
			}
			songChannelData = append(songChannelData, songChannel)
		}

		var songArtistList []dtos.GetSongByAnimeIdResponseSongArtist
		for _, artist := range song.Artist {
			songArtist := dtos.GetSongByAnimeIdResponseSongArtist{
				ID:    artist.ID,
				Name:  artist.Name,
				Image: artist.Image,
			}
			songArtistList = append(songArtistList, songArtist)
		}
		songData := dtos.GetSongByAnimeIdResponseSong{
			ID:          song.ID,
			Name:        song.Name,
			Type:        song.Type,
			Sequence:    song.Sequence,
			Image:       song.Image,
			Description: song.Description,
			Year:        song.Year,
			AnimeID:     song.AnimeID,
			SongChannel: songChannelData,
			SongArtist:  songArtistList,
		}
		if song.Type == songTypeOpening {
			openingSong = append(openingSong, songData)
		} else if song.Type == songTypeEnding {
			endingSong = append(endingSong, songData)
		} else if song.Type == songTypeSoundtrack {
			soundtrack = append(soundtrack, songData)
		}
	}
	songResponse := dtos.GetSongByAnimeIdResponse{
		OpeningSong:    openingSong,
		EndingSong:     endingSong,
		SoundtrackSong: soundtrack,
	}
	return &songResponse, nil
}

func (s songServiceImpl) CreateSongChannel(request *dtos.CreateSongChannelRequest) error {
	songChannelId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	linkEmbed := ""
	if request.Link != "" {
		linkEmbed, err = convertYouTubeURLToEmbed(request.Link)
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
	}

	songChannel := entities.SongChannel{
		ID:      songChannelId,
		Channel: request.Channel,
		Type:    request.Type,
		Link:    linkEmbed,
		SongID:  request.SongID,
		IsMain:  request.IsMain,
	}

	//save song
	err = s.songChannelRepo.Save(&songChannel)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s songServiceImpl) GetSongsByArtistId(artistId uuid.UUID) (*dtos.GetSongsByArtistResponse, error) {
	artist, err := s.artistRepo.GetById(artistId)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Artist not found")
	}

	songArtist, err := s.songArtistRepo.GetByArtistId(artistId)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	var songIds []uuid.UUID
	for _, item := range songArtist {
		songIds = append(songIds, item.SongId)
	}

	songs, err := s.songRepo.GetByIds(songIds)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	var songsResponse []dtos.GetSongsByArtistResponseSong
	for _, song := range songs {
		var songChannels []dtos.GetSongsByArtistResponseSongChannel
		for _, songChannel := range song.SongChannel {
			songChannels = append(songChannels, dtos.GetSongsByArtistResponseSongChannel{
				ID:      songChannel.ID,
				Channel: songChannel.Channel,
				Type:    songChannel.Type,
				Link:    songChannel.Link,
				IsMain:  songChannel.IsMain,
			})
		}
		songsResponse = append(songsResponse, dtos.GetSongsByArtistResponseSong{
			ID:             song.ID,
			Name:           song.Name,
			Image:          song.Image,
			Description:    song.Description,
			Year:           song.Year,
			Type:           song.Type,
			AnimeID:        song.AnimeID,
			AnimeName:      song.Anime.Name,
			AnimeWallpaper: song.Anime.Wallpaper,
			SongChannel:    songChannels,
		})
	}
	response := &dtos.GetSongsByArtistResponse{
		Aritst: dtos.GetSongsByArtistResponseArtist{
			ID:    artist.ID,
			Name:  artist.Name,
			Image: artist.Image,
		},
		Songs: songsResponse,
	}

	return response, nil
}
