package services

import (
	"slices"

	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeService interface {
	CreateEpisode(id uuid.UUID) error
	GetEpisodesByAnimeId(anime_id uuid.UUID, filter string) (*dtos.GetEpisodeResponse, error)
	UpdateEpisode(request dtos.UpdateEpisodeRequest) error
	AddCharactersToEpisode(request dtos.AddCharacterToEpisodeRequest) error
}

type episodeServiceImpl struct {
	episodeRepo          repositories.EpisodeRepository
	animeRepo            repositories.AnimeRepository
	episodeCharacterRepo repositories.EpisodeCharacterRepository
}

func NewEpisodeService(
	episodeRepo repositories.EpisodeRepository,
	animeRepo repositories.AnimeRepository,
	episodeCharacterRepo repositories.EpisodeCharacterRepository) EpisodeService {
	return &episodeServiceImpl{
		episodeRepo:          episodeRepo,
		animeRepo:            animeRepo,
		episodeCharacterRepo: episodeCharacterRepo}
}

func (s *episodeServiceImpl) CreateEpisode(animeId uuid.UUID) error {
	anime, err := s.animeRepo.GetById(animeId)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}
	if anime.IsCreateEpisode {
		return errs.NewBadRequestError("This Anime had already been create episode.")
	}

	var episodes []entities.Episode
	for i := 1; i <= anime.Episodes; i++ {
		episodeId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		episodes = append(episodes, entities.Episode{
			ID:      episodeId,
			AnimeID: animeId,
			Number:  uint(i),
		})
	}

	err = s.episodeRepo.BulkSave(episodes)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	err = s.animeRepo.UpdateIsCreateEpisode(anime.ID)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

// filter = FIRST_APPEARANCE  -> FirstAppearance default
// filter = APPEARANCE        -> Appearance
func (s *episodeServiceImpl) GetEpisodesByAnimeId(anime_id uuid.UUID, filter string) (*dtos.GetEpisodeResponse, error) {
	showCharcaterFormatList := []string{FIRST_APPEARANCE, APPEARANCE}
	if !slices.Contains(showCharcaterFormatList, filter) {
		errorMessage := "Invalid filter."
		logs.Error(errorMessage)
		return nil, errs.NewBadRequestError(errorMessage)
	}

	episodes, err := s.episodeRepo.GetByAnimeId(anime_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	var episodeIds []uuid.UUID
	for _, item := range episodes {
		episodeIds = append(episodeIds, item.ID)
	}

	episodeCharacterList, err := s.episodeCharacterRepo.GetByEpisodeIds(episodeIds)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	showCharcaterFormat := ""
	if filter == "" {
		showCharcaterFormat = FIRST_APPEARANCE
	} else {
		showCharcaterFormat = filter
	}

	episodeCharacterMap := make(map[uuid.UUID][]dtos.GetEpisodeResponseEpisodeCharacter)
	for _, item := range episodeCharacterList {
		episodeCharacter := dtos.GetEpisodeResponseEpisodeCharacter{
			ID:              item.CharacterID,
			Name:            item.Character.Name,
			FullName:        item.Character.LastName + "" + item.Character.FirstName,
			Image:           item.Character.Image,
			ImageStyleX:     item.Character.ImageStyleX,
			ImageStyleY:     item.Character.ImageStyleY,
			Description:     item.Description,
			FirstAppearance: item.FirstAppearance,
			Appearance:      item.Appearance,
		}
		if episodeCharacter.FirstAppearance && showCharcaterFormat == FIRST_APPEARANCE {
			episodeCharacterMap[item.EpisodeID] = append(episodeCharacterMap[item.EpisodeID], episodeCharacter)
		} else if episodeCharacter.Appearance && showCharcaterFormat == APPEARANCE {
			episodeCharacterMap[item.EpisodeID] = append(episodeCharacterMap[item.EpisodeID], episodeCharacter)
		}

	}

	var episodeResponse []dtos.GetEpisodeResponseEpisode
	for _, episode := range episodes {
		episodeResponse = append(episodeResponse, dtos.GetEpisodeResponseEpisode{
			ID:          episode.ID,
			Number:      episode.Number,
			Name:        episode.Name,
			NameThai:    episode.NameThai,
			NameEnglish: episode.NameEnglish,
			Characters:  episodeCharacterMap[episode.ID],
		})
	}

	return &dtos.GetEpisodeResponse{
		Episodes: episodeResponse,
	}, nil
}

func (s *episodeServiceImpl) UpdateEpisode(request dtos.UpdateEpisodeRequest) error {
	_, err := s.episodeRepo.GetById(request.ID)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Episode not found")
		}
		return errs.NewUnexpectedError()
	}

	episodeUpdate := entities.Episode{
		ID:          request.ID,
		Name:        request.Name,
		NameThai:    request.NameThai,
		NameEnglish: request.NameEnglish,
	}

	if err := s.episodeRepo.Update(&episodeUpdate); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *episodeServiceImpl) AddCharactersToEpisode(request dtos.AddCharacterToEpisodeRequest) error {
	_, err := s.episodeRepo.GetById(request.EpisodeID)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Episode not found")
		}
		return errs.NewUnexpectedError()
	}

	var episodeCharacters []entities.EpisodeCharacter
	for _, character := range request.Characters {
		episodeCharacterId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		episodeCharacters = append(episodeCharacters, entities.EpisodeCharacter{
			ID:              episodeCharacterId,
			EpisodeID:       request.EpisodeID,
			CharacterID:     character.ID,
			Description:     character.Description,
			Appearance:      character.Appearance,
			FirstAppearance: character.FirstAppearance,
		})
	}

	if err = s.episodeCharacterRepo.BulkSave(episodeCharacters); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}
