package services

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/adapters/aws"
	"github.com/Fourth1755/animap-go-api/internal/adapters/external_api"
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/Fourth1755/animap-go-api/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func encodeCursor(airedAt time.Time, id uuid.UUID) string {
	raw := airedAt.UTC().Format(time.RFC3339Nano) + "|" + id.String()
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(cursor string) (*time.Time, *uuid.UUID, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, nil, err
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid cursor format")
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, nil, err
	}
	id, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, nil, err
	}
	return &t, &id, nil
}

type AnimeService interface {
	CreateAnime(anime dtos.CreateAnimeRequest) error
	GetAnimeById(id uuid.UUID) (*dtos.GetAnimeByIdResponse, error)
	GetAnimes(query dtos.AnimeQueryDTO) (*dtos.AnimeListsResponse, error)
	UpdateAnime(anime entities.Anime) error
	DeleteAnime(id uuid.UUID) error
	GetAnimeByUserId(user_id uuid.UUID) ([]entities.UserAnime, error)
	GetAnimeByCategoryId(category_id uuid.UUID, query dtos.AnimeCursorQueryDTO) (*dtos.GetAnimeByCategoryIdResponse, error)
	GetAnimeByCategoryUniverseId(category_id uuid.UUID, query dtos.AnimeCursorQueryDTO) (*dtos.GetAnimeByCategoryUniverseIdResponse, error)
	AddCategoryToAnime(request dtos.EditCategoryToAnimeRequest) error
	AddCategoryUniverseToAnime(request dtos.EditCategoryUniverseToAnimeRequest) error
	GetAnimeBySeasonalAndYear(request dtos.GetAnimeBySeasonAndYearRequest) (*dtos.GetAnimeBySeasonAndYearResponse, error)
	GetAnimeByStudio(studioId uuid.UUID, query dtos.AnimeCursorQueryDTO) (*dtos.GetAnimeByStudioResponse, error)
	AddAnimePictures(request dtos.AddAnimePicturesRequest) error
	GetAnimePictures(animeID uuid.UUID) (*dtos.AnimeMediaResponse, error)
	CreateAnimeTrailers(request dtos.CreateAnimeTrailersRequest) error
}

type animeServiceImpl struct {
	repo                        repositories.AnimeRepository
	userRepo                    repositories.UserRepository
	animeCategoryRepo           repositories.AnimeCategoryRepository
	animeStudioRepo             repositories.AnimeStudioRepository
	animeTrailerRepo            repositories.AnimeTrailerRepository
	songRepo                    repositories.SongRepository
	categoryRepo                repositories.CategoryRepository
	animeCategorryUnivserseRepo repositories.AnimeCategoryUniverseRepository
	categoryUniverseRepo        repositories.CategoryUniverseRepository
	studioRepo                  repositories.StudioRepository
	episodeRepo                 repositories.EpisodeRepository
	s3Service                   aws.S3Service
	myAnimeListService          external_api.MyAnimeListService
}

func NewAnimeService(
	repo repositories.AnimeRepository,
	userRepo repositories.UserRepository,
	animeCategoryRepo repositories.AnimeCategoryRepository,
	animeStudioRepo repositories.AnimeStudioRepository,
	animeTrailerRepo repositories.AnimeTrailerRepository,
	songRepo repositories.SongRepository,
	categoryRepo repositories.CategoryRepository,
	animeCategorryUnivserseRepo repositories.AnimeCategoryUniverseRepository,
	categoryUniverseRepo repositories.CategoryUniverseRepository,
	studioRepo repositories.StudioRepository,
	episodeRepo repositories.EpisodeRepository,
	s3Service aws.S3Service,
	myAnimeListService external_api.MyAnimeListService,
) AnimeService {
	return &animeServiceImpl{
		repo:                        repo,
		userRepo:                    userRepo,
		animeCategoryRepo:           animeCategoryRepo,
		animeStudioRepo:             animeStudioRepo,
		animeTrailerRepo:            animeTrailerRepo,
		songRepo:                    songRepo,
		categoryRepo:                categoryRepo,
		animeCategorryUnivserseRepo: animeCategorryUnivserseRepo,
		categoryUniverseRepo:        categoryUniverseRepo,
		studioRepo:                  studioRepo,
		episodeRepo:                 episodeRepo,
		s3Service:                   s3Service,
		myAnimeListService:          myAnimeListService,
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
		ID:            animeId,
		Name:          request.Name,
		NameEnglish:   request.NameEnglish,
		NameThai:      request.NameThai,
		Episodes:      request.Episodes,
		Image:         request.Image,
		Description:   request.Description,
		Seasonal:      request.Seasonal,
		Year:          request.Year,
		Type:          request.Type,
		Duration:      request.Duration,
		Wallpaper:     request.Wallpaper,
		Trailer:       request.Trailer,
		TrailerEmbed:  trailerEmbed,
		AiredAt:       request.AiredAt,
		MyAnimeListID: request.MyAnimeListID,
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
	for _, category := range anime.Categories {
		categories = append(categories, dtos.AnimeDetailCategories{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	var categoryUniverse []dtos.AnimeDataUniverse
	for _, category := range anime.CategoryUniverses {
		categoryUniverse = append(categoryUniverse, dtos.AnimeDataUniverse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	var studios []dtos.AnimeDetailStduios
	for _, studio := range anime.Studios {
		studios = append(studios, dtos.AnimeDetailStduios{
			ID:   studio.ID,
			Name: studio.Name,
		})
	}

	var providers []dtos.AnimeDetailProvider
	for _, provider := range anime.Providers {
		providers = append(providers, dtos.AnimeDetailProvider{
			ID:    provider.ID,
			Name:  provider.Name,
			Image: provider.Image,
		})
	}

	myAnimeListData, err := s.myAnimeListService.GetAnimeDetail(int(anime.MyAnimeListID))
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	animeResponse := dtos.GetAnimeByIdResponse{}
	if !anime.IsSubAnime {
		animeResponse = dtos.GetAnimeByIdResponse{
			ID:               anime.ID,
			Name:             myAnimeListData.Title,
			NameEnglish:      myAnimeListData.AlternativeTitles.En,
			NameJapan:        myAnimeListData.AlternativeTitles.Ja,
			NameThai:         anime.NameThai,
			Episodes:         int(myAnimeListData.NumEpisodes),
			Seasonal:         myAnimeListData.StartSeason.Season,
			Year:             strconv.Itoa(myAnimeListData.StartSeason.Year),
			Image:            myAnimeListData.MainPicture.Large,
			Description:      myAnimeListData.Synopsis,
			Type:             anime.Type,
			Duration:         utils.FormatDuration(myAnimeListData.AverageEpisodeDuration),
			Categories:       categories,
			Wallpaper:        anime.Wallpaper,
			TrailerEmbed:     anime.TrailerEmbed,
			Studios:          studios,
			CategoryUniverse: categoryUniverse,
			MyAnimeListScore: myAnimeListData.Mean,
			IsSubAnime:       anime.IsSubAnime,
			Providers:        providers,
		}
	} else {
		animeResponse = dtos.GetAnimeByIdResponse{
			ID:               anime.ID,
			Name:             anime.Name,
			NameEnglish:      anime.NameEnglish,
			NameJapan:        "",
			NameThai:         anime.NameThai,
			Episodes:         anime.Episodes,
			Seasonal:         anime.Seasonal,
			Year:             anime.Year,
			Image:            anime.Image,
			Description:      anime.Description,
			Type:             anime.Type,
			Duration:         anime.Duration,
			Categories:       categories,
			Wallpaper:        anime.Wallpaper,
			TrailerEmbed:     anime.TrailerEmbed,
			Studios:          studios,
			CategoryUniverse: categoryUniverse,
			MyAnimeListScore: 0,
			IsSubAnime:       anime.IsSubAnime,
			Providers:        providers,
		}
	}

	return &animeResponse, nil
}

func (s *animeServiceImpl) GetAnimes(query dtos.AnimeQueryDTO) (*dtos.AnimeListsResponse, error) {
	animes, err := s.repo.GetAll(query)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	var animesDto []dtos.AnimeListResponse
	for _, anime := range animes {
		if anime.IsShow {
			animesDto = append(animesDto, dtos.AnimeListResponse{
				ID:       anime.ID,
				Name:     anime.Name,
				Episodes: anime.Episodes,
				Seasonal: anime.Seasonal,
				Year:     anime.Year,
				Image:    anime.Image,
			})
		}
	}
	return &dtos.AnimeListsResponse{Animes: animesDto}, nil
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

func (s *animeServiceImpl) GetAnimeByCategoryId(category_id uuid.UUID, query dtos.AnimeCursorQueryDTO) (*dtos.GetAnimeByCategoryIdResponse, error) {
	category, err := s.categoryRepo.GetById(category_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	var cursorTime *time.Time
	var cursorID *uuid.UUID
	if query.Cursor != "" {
		ct, cid, err := decodeCursor(query.Cursor)
		if err != nil {
			return nil, errs.NewBadRequestError("invalid cursor")
		}
		cursorTime, cursorID = ct, cid
	}

	animeCategories, err := s.animeCategoryRepo.GetByCategoryId(category_id, cursorTime, cursorID, limit+1)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	hasMore := len(animeCategories) > limit
	if hasMore {
		animeCategories = animeCategories[:limit]
	}

	var animesReponse []dtos.GetAnimeByCategoryIdResponseAnimeList
	for _, anime := range animeCategories {
		var studios []dtos.AnimeDetailStduios
		for _, studio := range anime.Studios {
			studios = append(studios, dtos.AnimeDetailStduios{
				ID:   studio.ID,
				Name: studio.Name,
			})
		}
		animesReponse = append(animesReponse, dtos.GetAnimeByCategoryIdResponseAnimeList{
			ID:       anime.ID,
			Name:     anime.Name,
			Episodes: anime.Episodes,
			Seasonal: anime.Seasonal,
			Year:     anime.Year,
			Image:    anime.Image,
			Studios:  studios,
		})
	}

	var nextCursor *string
	if hasMore && len(animeCategories) > 0 {
		last := animeCategories[len(animeCategories)-1]
		c := encodeCursor(last.AiredAt, last.ID)
		nextCursor = &c
	}

	return &dtos.GetAnimeByCategoryIdResponse{
		ID:         category.ID,
		Name:       category.Name,
		Wallpaper:  category.Image,
		AnimeList:  animesReponse,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *animeServiceImpl) GetAnimeByCategoryUniverseId(category_id uuid.UUID, query dtos.AnimeCursorQueryDTO) (*dtos.GetAnimeByCategoryUniverseIdResponse, error) {
	category, err := s.categoryUniverseRepo.GetById(category_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	var cursorTime *time.Time
	var cursorID *uuid.UUID
	if query.Cursor != "" {
		ct, cid, err := decodeCursor(query.Cursor)
		if err != nil {
			return nil, errs.NewBadRequestError("invalid cursor")
		}
		cursorTime, cursorID = ct, cid
	}

	animeCategories, err := s.animeCategorryUnivserseRepo.GetByCategoryUniverseId(category_id, cursorTime, cursorID, limit+1)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	hasMore := len(animeCategories) > limit
	if hasMore {
		animeCategories = animeCategories[:limit]
	}

	var animeIds []uuid.UUID
	for _, animeItem := range animeCategories {
		animeIds = append(animeIds, animeItem.AnimeID)
	}

	animeList, err := s.repo.GetByIds(animeIds)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	animeMap := make(map[uuid.UUID]entities.Anime, len(animeList))
	for _, a := range animeList {
		animeMap[a.ID] = a
	}

	var animesReponse []dtos.GetAnimeByCategoryUniverseIdResponseAnimeList
	for _, animeItem := range animeCategories {
		anime, ok := animeMap[animeItem.AnimeID]
		if !ok {
			continue
		}
		var categories []dtos.AnimeDetailCategories
		for _, category := range anime.Categories {
			categories = append(categories, dtos.AnimeDetailCategories{
				ID:   category.ID,
				Name: category.Name,
			})
		}

		var studios []dtos.AnimeDetailStduios
		for _, studio := range anime.Studios {
			studios = append(studios, dtos.AnimeDetailStduios{
				ID:   studio.ID,
				Name: studio.Name,
			})
		}

		animesReponse = append(animesReponse, dtos.GetAnimeByCategoryUniverseIdResponseAnimeList{
			ID:          anime.ID,
			Name:        anime.Name,
			Episodes:    anime.Episodes,
			Seasonal:    anime.Seasonal,
			Year:        anime.Year,
			Image:       anime.Image,
			Description: anime.Description,
			Type:        anime.Type,
			Duration:    anime.Duration,
			Studios:     studios,
			Categories:  categories,
		})
	}

	var nextCursor *string
	if hasMore && len(animeCategories) > 0 {
		last := animeCategories[len(animeCategories)-1]
		c := encodeCursor(last.Anime.AiredAt, last.AnimeID)
		nextCursor = &c
	}

	return &dtos.GetAnimeByCategoryUniverseIdResponse{
		ID:         category.ID,
		Name:       category.Name,
		Wallpaper:  category.Image,
		AnimeList:  animesReponse,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// need to enhance
func (s *animeServiceImpl) AddCategoryToAnime(request dtos.EditCategoryToAnimeRequest) error {
	animeCategory := []entities.AnimeCategory{}
	// check dup
	animeCategoryDup, err := s.animeCategoryRepo.GetByAnimeIdAndCategoryIds(request.AnimeID, request.CategoryID)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	if len(animeCategoryDup) != 0 {
		errMessage := "Category in anime is duplicate."
		logs.Error(errMessage)
		return errs.NewBadRequestError(errMessage)
	}

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

// need to enhance
func (s *animeServiceImpl) AddCategoryUniverseToAnime(request dtos.EditCategoryUniverseToAnimeRequest) error {
	var allAnimeCategories []entities.AnimeCategoryUniverse

	for _, animeID := range request.AnimeIDs {
		animeCategory := []entities.AnimeCategoryUniverse{}
		// check dup
		animeCategoryDup, err := s.animeCategorryUnivserseRepo.GetByAnimeIdsAndCategoryUniverseIds([]uuid.UUID{animeID}, request.CategoryUniverseID)
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		if len(animeCategoryDup) != 0 {
			errMessage := "Category Universe in anime is duplicate."
			logs.Error(errMessage)
			return errs.NewBadRequestError(errMessage)
		}

		for _, categoryUniverseID := range request.CategoryUniverseID {
			animeCategoryId, err := uuid.NewV7()
			if err != nil {
				logs.Error(err.Error())
				return errs.NewUnexpectedError()
			}
			animeCategory = append(animeCategory, entities.AnimeCategoryUniverse{
				ID:                 animeCategoryId,
				AnimeID:            animeID,
				CategoryUniverseID: categoryUniverseID,
			})
		}
		allAnimeCategories = append(allAnimeCategories, animeCategory...)
	}

	if err := s.animeCategorryUnivserseRepo.Save(allAnimeCategories); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *animeServiceImpl) GetAnimeBySeasonalAndYear(request dtos.GetAnimeBySeasonAndYearRequest) (*dtos.GetAnimeBySeasonAndYearResponse, error) {
	seasonal := []string{"winter", "spring", "summer", "fall"}
	if !slices.Contains(seasonal, request.Seasonal) {
		errorMessage := "Invalid seasonal request."
		logs.Error(errorMessage)
		return nil, errs.NewBadRequestError(errorMessage)
	}
	if len(request.Year) != 4 {
		errorMessage := "Invalid year request."
		logs.Error(errorMessage)
		return nil, errs.NewBadRequestError(errorMessage)
	}
	animeList, err := s.repo.GetBySeasonalAndYear(request)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	var animesReponse []dtos.GetAnimeBySeasonAndYearResponseAnime
	for _, anime := range animeList {
		if anime.IsShow {
			var studios []dtos.AnimeDetailStduios
			for _, studio := range anime.Studios {
				studios = append(studios, dtos.AnimeDetailStduios{
					ID:   studio.ID,
					Name: studio.Name,
				})
			}

			animesReponse = append(animesReponse, dtos.GetAnimeBySeasonAndYearResponseAnime{
				ID:           anime.ID,
				Name:         anime.Name,
				NameEnglish:  anime.NameEnglish,
				Episodes:     anime.Episodes,
				Seasonal:     anime.Seasonal,
				Year:         anime.Year,
				Image:        anime.Image,
				Description:  anime.Description,
				Type:         anime.Type,
				Duration:     anime.Duration,
				Wallpaper:    anime.Wallpaper,
				Trailer:      anime.Trailer,
				TrailerEmbed: anime.TrailerEmbed,
				Studios:      studios,
			})
		}
	}
	return &dtos.GetAnimeBySeasonAndYearResponse{
		Year:      request.Year,
		Seasonal:  request.Seasonal,
		AnimeList: animesReponse,
	}, nil
}

func (s *animeServiceImpl) GetAnimeByStudio(studioId uuid.UUID, query dtos.AnimeCursorQueryDTO) (*dtos.GetAnimeByStudioResponse, error) {
	studio, err := s.studioRepo.GetById(studioId)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Studio Not Found" + studioId.String())
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}

	var cursorTime *time.Time
	var cursorID *uuid.UUID
	if query.Cursor != "" {
		ct, cid, err := decodeCursor(query.Cursor)
		if err != nil {
			return nil, errs.NewBadRequestError("invalid cursor")
		}
		cursorTime, cursorID = ct, cid
	}

	animeStudio, err := s.animeStudioRepo.GetByStudioId(studioId, cursorTime, cursorID, limit+1)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	hasMore := len(animeStudio) > limit
	if hasMore {
		animeStudio = animeStudio[:limit]
	}

	var animeIds []uuid.UUID
	for _, animeItem := range animeStudio {
		animeIds = append(animeIds, animeItem.AnimeID)
	}

	animeList, err := s.repo.GetByIds(animeIds)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	animeMap := make(map[uuid.UUID]entities.Anime, len(animeList))
	for _, a := range animeList {
		animeMap[a.ID] = a
	}

	var animesReponse []dtos.GetAnimeByStudioResponseAnimeList
	for _, animeItem := range animeStudio {
		anime, ok := animeMap[animeItem.AnimeID]
		if !ok {
			continue
		}
		var categories []dtos.AnimeDetailCategories
		for _, category := range anime.Categories {
			categories = append(categories, dtos.AnimeDetailCategories{
				ID:   category.ID,
				Name: category.Name,
			})
		}

		var studios []dtos.AnimeDetailStduios
		for _, s := range anime.Studios {
			studios = append(studios, dtos.AnimeDetailStduios{
				ID:   s.ID,
				Name: s.Name,
			})
		}

		animesReponse = append(animesReponse, dtos.GetAnimeByStudioResponseAnimeList{
			ID:          anime.ID,
			Name:        anime.Name,
			Episodes:    anime.Episodes,
			Seasonal:    anime.Seasonal,
			Year:        anime.Year,
			Image:       anime.Image,
			Description: anime.Description,
			Type:        anime.Type,
			Duration:    anime.Duration,
			Studios:     studios,
			Categories:  categories,
			Wallpaper:   anime.Wallpaper,
		})
	}

	var nextCursor *string
	if hasMore && len(animeStudio) > 0 {
		last := animeStudio[len(animeStudio)-1]
		c := encodeCursor(last.Anime.AiredAt, last.AnimeID)
		nextCursor = &c
	}

	return &dtos.GetAnimeByStudioResponse{
		ID:         studio.ID,
		Name:       studio.Name,
		Wallpaper:  studio.Image,
		MainColor:  studio.MainColor,
		AnimeList:  animesReponse,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
func (s *animeServiceImpl) AddAnimePictures(request dtos.AddAnimePicturesRequest) error {
	anime, err := s.repo.GetById(request.AnimeID)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}

	updated := append(anime.Pictures, request.Pictures...)
	if err := s.repo.UpdatePictures(request.AnimeID, updated); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}

func (s *animeServiceImpl) GetAnimePictures(animeID uuid.UUID) (*dtos.AnimeMediaResponse, error) {
	trailer, err := s.animeTrailerRepo.GetByAnimeId(animeID)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	var response []dtos.AnimeMediaDataResponse
	for _, t := range trailer {
		id := t.ID
		url := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", t.VideoID)
		embedURL := fmt.Sprintf("https://www.youtube.com/embed/%s", t.VideoID)
		response = append(response, dtos.AnimeMediaDataResponse{
			ID:       &id,
			Type:     "VIDEO",
			URL:      url,
			EmbedURL: embedURL,
		})
	}
	anime, err := s.repo.GetById(animeID)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	for _, picURL := range anime.Pictures {
		response = append(response, dtos.AnimeMediaDataResponse{
			Type:     "PICTURE",
			URL:      picURL,
			EmbedURL: "",
		})
	}

	return &dtos.AnimeMediaResponse{
		Data: response,
	}, nil
}

func (s *animeServiceImpl) CreateAnimeTrailers(request dtos.CreateAnimeTrailersRequest) error {
	// Check if anime exists
	if _, err := s.repo.GetById(request.AnimeID); err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Anime not found")
		}
		return errs.NewUnexpectedError()
	}

	var trailers []entities.AnimeTrailer
	for _, trailerReq := range request.Trailers {
		trailerID, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		trailers = append(trailers, entities.AnimeTrailer{
			ID:      trailerID,
			AnimeID: request.AnimeID,
			Name:    trailerReq.Name,
			VideoID: trailerReq.VideoID,
		})
	}

	if err := s.animeTrailerRepo.CreateAnimeTrailers(trailers); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}

	return nil
}
