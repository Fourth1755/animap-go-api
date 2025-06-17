package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CharacterService interface {
	CreateCharacter(dtos.CreateCharacterRequest) error
	GetCharacterByAnimeId(animeId uuid.UUID) (*dtos.GetCharacterByAnimeIdResponse, error)
}

type characterServiceImpl struct {
	characterRepo      repositories.CharacterRepository
	animeCharacterRepo repositories.AnimeCharacterRepository
	animeRepo          repositories.AnimeRepository
}

func NewCharacterService(
	characterRepo repositories.CharacterRepository,
	animeCharacterRepo repositories.AnimeCharacterRepository,
	animeRepo repositories.AnimeRepository) CharacterService {
	return &characterServiceImpl{
		characterRepo:      characterRepo,
		animeCharacterRepo: animeCharacterRepo,
		animeRepo:          animeRepo,
	}
}

func (s *characterServiceImpl) CreateCharacter(request dtos.CreateCharacterRequest) error {
	characterId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	characterCreate := entities.Character{
		ID:              characterId,
		Name:            request.Name,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		NameThai:        request.NameThai,
		FirstNameThai:   request.FirstNameThai,
		LastNameThai:    request.LastNameThai,
		FirstNameJapan:  request.FirstNameJapan,
		LastNameJapan:   request.LastNameJapan,
		Image:           request.Image,
		ImageStyleX:     request.ImageStyleX,
		ImageStyleY:     request.ImageStyleY,
		IsMainCharacter: request.IsMainCharacter,
	}

	character, err := s.characterRepo.Save(&characterCreate)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	animeCharacterId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	animeCharacter := entities.AnimeCharacter{
		ID:          animeCharacterId,
		AnimeID:     request.AnimeId,
		CharacterID: character.ID,
	}

	if err = s.animeCharacterRepo.Save(&animeCharacter); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *characterServiceImpl) GetCharacterByAnimeId(animeId uuid.UUID) (*dtos.GetCharacterByAnimeIdResponse, error) {
	_, err := s.animeRepo.GetById(animeId)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError("Anime not found")
		}
		return nil, errs.NewUnexpectedError()
	}

	animeCharacter, err := s.animeCharacterRepo.GetByAnimeId(animeId)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	characterResponse := []dtos.GetCharacterByAnimeIdResponseCharacter{}
	for _, character := range animeCharacter {
		characterItem := dtos.GetCharacterByAnimeIdResponseCharacter{
			CharacterID:     character.Character.ID,
			Number:          character.Number,
			Name:            character.Character.Name,
			FullName:        character.Character.FirstName + " " + character.Character.LastName,
			NameThai:        character.Character.NameThai,
			FullNameThai:    character.Character.FirstNameThai + " " + character.Character.LastNameThai,
			FullNameJapan:   character.Character.LastNameJapan + " " + character.Character.FirstNameJapan,
			Image:           character.Character.Image,
			ImageStyleX:     character.Character.ImageStyleX,
			ImageStyleY:     character.Character.ImageStyleY,
			IsMainCharacter: character.Character.IsMainCharacter,
		}
		characterResponse = append(characterResponse, characterItem)
	}

	return &dtos.GetCharacterByAnimeIdResponse{
		Character: characterResponse,
	}, nil
}
