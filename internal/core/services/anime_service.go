package services

import (
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
	GetAnimeByCategoryId(category_id uuid.UUID) (*dtos.GetAnimeByCategoryIdResponse, error)
	GetAnimeByCategoryUniverseId(category_id uuid.UUID) (*dtos.GetAnimeByCategoryUniverseIdResponse, error)
	AddCategoryToAnime(request dtos.EditCategoryToAnimeRequest) error
	AddCategoryUniverseToAnime(request dtos.EditCategoryUniverseToAnimeRequest) error
	GetAnimeBySeasonalAndYear(request dtos.GetAnimeBySeasonAndYearRequest) (*dtos.GetAnimeBySeasonAndYearResponse, error)
	GetAnimeByStudio(studioId uuid.UUID) (*dtos.GetAnimeByStudioResponse, error)
	MigrateAnime(req dtos.MigrateAnimeRequest) error
}

type animeServiceImpl struct {
	repo                        repositories.AnimeRepository
	userRepo                    repositories.UserRepository
	animeCategoryRepo           repositories.AnimeCategoryRepository
	animeStudioRepo             repositories.AnimeStudioRepository
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

	myAnimeListData, err := s.myAnimeListService.GetAnimeDetail(int(anime.MyAnimeListID))
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	animeResponse := dtos.GetAnimeByIdResponse{
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
		Duration:         anime.Duration,
		Categories:       categories,
		Wallpaper:        anime.Wallpaper,
		TrailerEmbed:     anime.TrailerEmbed,
		Studios:          studios,
		CategoryUniverse: categoryUniverse,
		MyAnimeListScore: myAnimeListData.Mean,
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

func (s *animeServiceImpl) GetAnimeByCategoryId(category_id uuid.UUID) (*dtos.GetAnimeByCategoryIdResponse, error) {
	category, err := s.categoryRepo.GetById(category_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	animeCategories, err := s.animeCategoryRepo.GetByCategoryId(category_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}
	var animesReponse []dtos.GetAnimeByCategoryIdResponseAnimeList
	for _, anime := range animeCategories {
		animesReponse = append(animesReponse, dtos.GetAnimeByCategoryIdResponseAnimeList{
			ID:       anime.Anime.ID,
			Name:     anime.Anime.Name,
			Episodes: anime.Anime.Episodes,
			Seasonal: anime.Anime.Seasonal,
			Year:     anime.Anime.Year,
			Image:    anime.Anime.Image,
		})
	}
	return &dtos.GetAnimeByCategoryIdResponse{
		ID:        category.ID,
		Name:      category.Name,
		Wallpaper: category.Image,
		AnimeList: animesReponse,
	}, nil
}

func (s *animeServiceImpl) GetAnimeByCategoryUniverseId(category_id uuid.UUID) (*dtos.GetAnimeByCategoryUniverseIdResponse, error) {
	category, err := s.categoryUniverseRepo.GetById(category_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	animeCategories, err := s.animeCategorryUnivserseRepo.GetByCategoryUniverseId(category_id)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
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
	var animesReponse []dtos.GetAnimeByCategoryUniverseIdResponseAnimeList
	for _, anime := range animeList {
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
	return &dtos.GetAnimeByCategoryUniverseIdResponse{
		ID:        category.ID,
		Name:      category.Name,
		Wallpaper: category.Image,
		AnimeList: animesReponse,
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
	animeCategory := []entities.AnimeCategoryUniverse{}
	// check dup
	animeCategoryDup, err := s.animeCategorryUnivserseRepo.GetByAnimeIdAndCategoryUniverseIds(request.AnimeID, request.CategoryUniverseID)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	if len(animeCategoryDup) != 0 {
		errMessage := "Category Universe in anime is duplicate."
		logs.Error(errMessage)
		return errs.NewBadRequestError(errMessage)
	}

	for _, catrgory := range request.CategoryUniverseID {
		animeCategoryId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		animeCategory = append(animeCategory, entities.AnimeCategoryUniverse{
			ID:                 animeCategoryId,
			AnimeID:            request.AnimeID,
			CategoryUniverseID: catrgory,
		})
	}
	if err := s.animeCategorryUnivserseRepo.Save(animeCategory); err != nil {
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
		})
	}
	return &dtos.GetAnimeBySeasonAndYearResponse{
		Year:      request.Year,
		Seasonal:  request.Seasonal,
		AnimeList: animesReponse,
	}, nil
}

func (s *animeServiceImpl) GetAnimeByStudio(studioId uuid.UUID) (*dtos.GetAnimeByStudioResponse, error) {
	studio, err := s.studioRepo.GetById(studioId)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Studio Not Found" + studioId.String())
	}
	animeStudio, err := s.animeStudioRepo.GetByStudioId(studioId)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
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
	var animesReponse []dtos.GetAnimeByStudioResponseAnimeList
	for _, anime := range animeList {
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
	return &dtos.GetAnimeByStudioResponse{
		ID:        studio.ID,
		Name:      studio.Name,
		Wallpaper: studio.Image,
		MainColor: studio.MainColor,
		AnimeList: animesReponse,
	}, nil
}

// ConvertDateStringToTime แปลง string "YYYY-MM-DD" -> time.Time
func ConvertDateStringToTime(dateStr string) (time.Time, error) {
	// layout ต้องเป็น "2006-01-02" ถ้าข้อมูลเป็น "2025-07-06"
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

func (s *animeServiceImpl) insertAnimeFromMyAnimeList(i int, animeMal *external_api.GetAnimeDetailResponse) error {
	animeId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	nameEnglish := ""
	if animeMal.AlternativeTitles.En != "" {
		nameEnglish = animeMal.AlternativeTitles.En
	}
	fmt.Println(animeMal.StartDate)
	airedAt, err := ConvertDateStringToTime(animeMal.StartDate)
	if err != nil {
		logs.Error("can cont convert data to string tine " + err.Error())
		return errs.NewUnexpectedError()
	}

	anime := entities.Anime{
		ID:            animeId,
		Name:          animeMal.Title,
		NameEnglish:   nameEnglish,
		NameThai:      "",
		Episodes:      int(animeMal.NumEpisodes),
		Seasonal:      animeMal.StartSeason.Season,
		Year:          strconv.Itoa(animeMal.StartSeason.Year),
		Rating:        animeMal.Rating,
		MediaType:     animeMal.MediaType,
		AiredAt:       airedAt,
		MyAnimeListID: uint64(i),
	}

	_, err = s.repo.Save(anime)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	fmt.Println(strconv.Itoa(i) + "Id that are created")
	return nil
}

func (s *animeServiceImpl) MigrateAnime(req dtos.MigrateAnimeRequest) error {
	fmt.Println("Strat migrate")
	fmt.Printf("start at %d", req.StartAnimeId)
	fmt.Printf("end at %d", req.EndAnimeId)
	for i := req.StartAnimeId; i <= req.EndAnimeId; i++ {
		animeMal, err := s.myAnimeListService.GetAnimeDetail(i)
		fmt.Println(animeMal)
		fmt.Println(err)
		if err != nil {
			logs.Error(err.Error())
		} else if animeMal != nil {
			if animeMal.MediaType != "music" && animeMal.MediaType != "cm" && animeMal.MediaType != "ona" {
				err = s.insertAnimeFromMyAnimeList(i, animeMal)
				if err != nil {
					logs.Error(err.Error())
				}
			}
		}
	}
	fmt.Println("End migrate")
	return nil
}
