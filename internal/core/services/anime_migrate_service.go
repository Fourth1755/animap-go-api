package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/adapters/external_api"
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeMigrateService interface {
	MigrateAnime(req dtos.MigrateAnimeRequest) error
}

type animeMigrateServiceImpl struct {
	repo               repositories.AnimeRepository
	animeCategoryRepo  repositories.AnimeCategoryRepository
	animeStudioRepo    repositories.AnimeStudioRepository
	categoryRepo       repositories.CategoryRepository
	studioRepo         repositories.StudioRepository
	episodeRepo        repositories.EpisodeRepository
	myAnimeListService external_api.MyAnimeListService
}

func NewAnimeMigrateService(
	repo repositories.AnimeRepository,
	animeCategoryRepo repositories.AnimeCategoryRepository,
	animeStudioRepo repositories.AnimeStudioRepository,
	categoryRepo repositories.CategoryRepository,
	studioRepo repositories.StudioRepository,
	episodeRepo repositories.EpisodeRepository,
	myAnimeListService external_api.MyAnimeListService,
) AnimeMigrateService {
	return &animeMigrateServiceImpl{
		repo:               repo,
		animeCategoryRepo:  animeCategoryRepo,
		animeStudioRepo:    animeStudioRepo,
		categoryRepo:       categoryRepo,
		studioRepo:         studioRepo,
		episodeRepo:        episodeRepo,
		myAnimeListService: myAnimeListService,
	}
}

// ConvertDateStringToTime แปลง string "YYYY-MM-DD" -> time.Time
func ConvertDateStringToTime(dateStr string) (time.Time, error) {
	// layout ต้องเป็น "2006-01-02" ถ้าข้อมูลเป็น "2025-07-06"
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

func (s *animeMigrateServiceImpl) insertAnimeFromMyAnimeList(i int, animeMal *external_api.GetAnimeDetailResponse) error {
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
		Image:         animeMal.MainPicture.Medium,
	}

	_, err = s.repo.Save(anime)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	fmt.Println(strconv.Itoa(i) + "Id that are created")
	return nil
}

func (s *animeMigrateServiceImpl) updateStudioAnime(i int, animeMal *external_api.GetAnimeDetailResponse) error {
	animeStudios := []entities.AnimeStudio{}
	for _, item := range animeMal.Studios {
		animeStudioId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		var studioId uuid.UUID
		studio, err := s.studioRepo.GetByMyAnimeListId(item.ID)
		if err != nil {
			logs.Error(err.Error())
			if err == gorm.ErrRecordNotFound {
				studioNewId, err := uuid.NewV7()
				if err != nil {
					logs.Error(err.Error())
					return errs.NewUnexpectedError()
				}
				studioNew := entities.Studio{
					ID:            studioNewId,
					Name:          item.Name,
					MyAnimeListID: item.ID,
				}
				result, err := s.studioRepo.Save(studioNew)
				if err != nil {
					logs.Error(err.Error())
					return errs.NewUnexpectedError()
				}
				studioId = result.ID
			} else {
				return errs.NewUnexpectedError()
			}
		} else {
			studioId = studio.ID
		}

		anime, err := s.repo.GetByMyAnimeListId(i)
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		animeStudios = append(animeStudios, entities.AnimeStudio{
			ID:       animeStudioId,
			StudioId: studioId,
			AnimeID:  anime.ID,
		})
	}
	if err := s.animeStudioRepo.Save(animeStudios); err != nil {
		return err
	}
	return nil
}

func (s *animeMigrateServiceImpl) updateCategoryAnime(i int, animeMal *external_api.GetAnimeDetailResponse) error {
	animeCategory := []entities.AnimeCategory{}
	for _, item := range animeMal.Genres {
		animeCategoryId, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		var categoryId uuid.UUID
		category, err := s.categoryRepo.GetByMyAnimeListId(item.ID)
		if err != nil {
			logs.Error(err.Error())
			if err == gorm.ErrRecordNotFound {
				categoryNewId, err := uuid.NewV7()
				if err != nil {
					logs.Error(err.Error())
					return errs.NewUnexpectedError()
				}
				categoryNew := entities.Category{
					ID:            categoryNewId,
					Name:          item.Name,
					MyAnimeListID: item.ID,
				}
				result, err := s.categoryRepo.Save(categoryNew)
				if err != nil {
					logs.Error(err.Error())
					return errs.NewUnexpectedError()
				}
				categoryId = result.ID
			} else {
				return errs.NewUnexpectedError()
			}
		} else {
			categoryId = category.ID
		}
		anime, err := s.repo.GetByMyAnimeListId(i)
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		var categoryList []uuid.UUID
		categoryList = append(categoryList, categoryId)
		animeCategoryDup, err := s.animeCategoryRepo.GetByAnimeIdAndCategoryIds(anime.ID, categoryList)
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		if len(animeCategoryDup) != 0 {
			errMessage := "Category in anime is duplicate."
			fmt.Println(errMessage)
		} else {
			animeCategory = append(animeCategory, entities.AnimeCategory{
				ID:         animeCategoryId,
				CategoryID: categoryId,
				AnimeID:    anime.ID,
			})
		}
	}

	if err := s.animeCategoryRepo.Save(animeCategory); err != nil {
		return err
	}
	return nil
}

func (s *animeMigrateServiceImpl) MigrateAnime(req dtos.MigrateAnimeRequest) error {
	fmt.Println("Strat migrate")
	count := 0
	for i := req.StartAnimeId; i <= req.EndAnimeId; i++ {
		animeMal, err := s.myAnimeListService.GetAnimeDetail(i)
		fmt.Println(animeMal)
		fmt.Println(err)
		if err != nil {
			logs.Error(err.Error())
		} else if animeMal != nil {
			animeData, err := s.repo.GetByMyAnimeListId(i)
			if err != nil {
				logs.Error(err.Error())
			}
			if animeData != nil {
				fmt.Println("data has been in database")
			}

			if animeData == nil && animeMal.MediaType != "music" && animeMal.MediaType != "cm" && animeMal.MediaType != "ona" && animeMal.MediaType != "pv" {
				fmt.Println("Insert anime")
				err = s.insertAnimeFromMyAnimeList(i, animeMal)
				if err != nil {
					logs.Error(err.Error())
				}

				err = s.updateStudioAnime(i, animeMal)
				if err != nil {
					logs.Error(err.Error())
				}

				err = s.updateCategoryAnime(i, animeMal)
				if err != nil {
					logs.Error(err.Error())
				}
				count++
			}
		}
	}
	fmt.Println("count")
	fmt.Println(count)
	fmt.Println("End migrate")
	return nil
}

func (s *animeMigrateServiceImpl) MigrateAnimeOld(req dtos.MigrateAnimeRequest) error {
	fmt.Println("Strat migrate")
	count := 0
	for i := req.StartAnimeId; i <= req.EndAnimeId; i++ {
		animeData, err := s.repo.GetByMyAnimeListId(i)
		if err != nil {
			logs.Error(err.Error())
		}
		if animeData != nil {
			fmt.Println("data has been in database")
			animeMal, err := s.myAnimeListService.GetAnimeDetail(i)
			count++
			if err != nil {
				logs.Error(err.Error())
			}
			err = s.updateCategoryAnime(i, animeMal)
			if err != nil {
				logs.Error(err.Error())
			}
			fmt.Println("update success")
		}
	}
	fmt.Println(count)
	fmt.Println("End migrate")
	return nil
}
