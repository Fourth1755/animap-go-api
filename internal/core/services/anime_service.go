package services

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeService interface {
	CreateAnime(anime dtos.CreateAnimeRequest) error
	GetAnimeById(id uuid.UUID) (*dtos.GetAnimeByIdResponse, error)
	GetAnimes(query dtos.AnimeQueryDTO) ([]dtos.AnimeListResponse, error)
	UpdateAnime(anime entities.Anime) error
	DeleteAnime(id uuid.UUID) error
	GetAnimeByUserId(user_id uuid.UUID) ([]entities.UserAnime, error)
	GetAnimeByCategoryId(category_id uuid.UUID) ([]dtos.AnimeListResponse, error)
	AddCategoryToAnime(request dtos.EditCategoryToAnimeRequest) error
}

type animeServiceImpl struct {
	repo              repositories.AnimeRepository
	userRepo          repositories.UserRepository
	animeCategoryRepo repositories.AnimeCategoryRepository
	animeStudioRepo   repositories.AnimeStudioRepository
	songRepo          repositories.SongRepository
}

func NewAnimeService(
	repo repositories.AnimeRepository,
	userRepo repositories.UserRepository,
	animeCategoryRepo repositories.AnimeCategoryRepository,
	animeStudioRepo repositories.AnimeStudioRepository,
	songRepo repositories.SongRepository) AnimeService {
	return &animeServiceImpl{
		repo:              repo,
		userRepo:          userRepo,
		animeCategoryRepo: animeCategoryRepo,
		animeStudioRepo:   animeStudioRepo,
		songRepo:          songRepo,
	}
}

const (
	ANIME_TYPE_TV    = 1
	ANIME_TYPE_MOVIE = 2
)

func (s *animeServiceImpl) CreateAnime(request dtos.CreateAnimeRequest) error {
	animeId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	trailerEmbed := ""
	if request.Trailer != "" {
		trailerEmbed, err = convertYouTubeURLToEmbed(request.Trailer)
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
	}
	anime := entities.Anime{
		ID:           animeId,
		Name:         request.Name,
		NameEnglish:  request.NameEnglish,
		NameThai:     request.NameThai,
		Episodes:     request.Episodes,
		Image:        request.Image,
		Description:  request.Description,
		Seasonal:     request.Seasonal,
		Year:         request.Year,
		Type:         request.Type,
		Duration:     request.Duration,
		Wallpaper:    request.Wallpaper,
		Trailer:      request.Trailer,
		TrailerEmbed: trailerEmbed,
	}

	animeCreate, err := s.repo.Save(anime)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	animeStudios := []entities.AnimeStudio{}
	for _, studio := range request.Studio {
		animeStudioId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}

		animeStudios = append(animeStudios, entities.AnimeStudio{
			ID:       animeStudioId,
			StudioId: uuid.MustParse(studio),
			AnimeID:  animeCreate.ID,
		})
	}

	if err := s.animeStudioRepo.Save(animeStudios); err != nil {
		return err
	}

	return nil
}

func convertYouTubeURLToEmbed(youtubeURL string) (string, error) {
	parsedURL, err := url.Parse(youtubeURL)
	if err != nil {
		return "", err
	}

	var videoID string

	switch parsedURL.Host {
	case "www.youtube.com", "youtube.com":
		queryParams := parsedURL.Query()
		videoID = queryParams.Get("v")
	case "youtu.be":
		videoID = strings.TrimPrefix(parsedURL.Path, "/")
	default:
		return "", fmt.Errorf("unsupported YouTube URL")
	}

	if videoID == "" {
		return "", fmt.Errorf("could not extract video ID")
	}

	embedURL := fmt.Sprintf("https://www.youtube.com/embed/%s", videoID)
	return embedURL, nil
}

func (s *animeServiceImpl) GetAnimeById(id uuid.UUID) (*dtos.GetAnimeByIdResponse, error) {
	anime, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError("Anime not found")
		}
		return nil, errs.NewUnexpectedError()
	}

	var categories []dtos.AnimeDetailCategories
	var universe dtos.AnimeDataUniverse
	for _, category := range anime.Categories {
		if category.IsUniverse {
			universe = dtos.AnimeDataUniverse{
				ID:   category.ID,
				Name: category.Name,
			}
		} else {
			categories = append(categories, dtos.AnimeDetailCategories{
				ID:   category.ID,
				Name: category.Name,
			})
		}
	}

	var studios []dtos.AnimeDetailStduios
	for _, studio := range anime.Studios {
		studios = append(studios, dtos.AnimeDetailStduios{
			ID:   studio.ID,
			Name: studio.Name,
		})
	}

	animeResponse := dtos.GetAnimeByIdResponse{
		ID:          anime.ID,
		Name:        anime.Name,
		NameEnglish: anime.NameEnglish,
		Episodes:    anime.Episodes,
		Seasonal:    anime.Seasonal,
		Year:        anime.Year,
		Image:       anime.Image,
		Description: anime.Description,
		Type:        anime.Type,
		Duration:    anime.Duration,
		Categories:  categories,
		Wallpaper:   anime.Wallpaper,
		Trailer:     anime.Trailer,
		Studios:     studios,
		Universe:    universe,
	}
	return &animeResponse, nil
}

func (s *animeServiceImpl) GetAnimes(query dtos.AnimeQueryDTO) ([]dtos.AnimeListResponse, error) {
	animes, err := s.repo.GetAll(query)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	var animesDto []dtos.AnimeListResponse
	for _, anime := range animes {
		animesDto = append(animesDto, dtos.AnimeListResponse{
			ID:       anime.ID,
			Name:     anime.Name,
			Episodes: anime.Episodes,
			Seasonal: anime.Seasonal,
			Year:     anime.Year,
			Image:    anime.Image,
		})
	}
	return animesDto, nil
}

func (s *animeServiceImpl) UpdateAnime(anime entities.Anime) error {
	_, err := s.repo.GetById(anime.ID)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}

	if err := s.repo.Update(&anime); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *animeServiceImpl) DeleteAnime(id uuid.UUID) error {
	_, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}

	if err := s.repo.Delete(id); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *animeServiceImpl) GetAnimeByUserId(userId uuid.UUID) ([]entities.UserAnime, error) {
	if _, err := s.userRepo.GetById(userId); err != nil {
		logs.Error(err.Error())
		return nil, errs.NewNotFoundError("User not found")
	}
	result, err := s.repo.GetByUserId(userId)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	return result, nil
}

func (s *animeServiceImpl) GetAnimeByCategoryId(category_id uuid.UUID) ([]dtos.AnimeListResponse, error) {
	animeCategories, err := s.animeCategoryRepo.GetByCategoryId(category_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	var animesReponse []dtos.AnimeListResponse
	for _, anime := range animeCategories {
		animesReponse = append(animesReponse, dtos.AnimeListResponse{
			ID:       anime.ID,
			Name:     anime.Anime.Name,
			Episodes: anime.Anime.Episodes,
			Seasonal: anime.Anime.Seasonal,
			Year:     anime.Anime.Year,
			Image:    anime.Anime.Image,
		})
	}
	return animesReponse, nil
}

// need to enhance
func (s *animeServiceImpl) AddCategoryToAnime(request dtos.EditCategoryToAnimeRequest) error {
	animeCategory := []entities.AnimeCategory{}
	for _, catrgory := range request.CategoryID {
		animeCategoryId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		animeCategory = append(animeCategory, entities.AnimeCategory{
			ID:         animeCategoryId,
			AnimeID:    request.AnimeID,
			CategoryID: catrgory,
		})
	}
	if err := s.animeCategoryRepo.Save(animeCategory); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}
