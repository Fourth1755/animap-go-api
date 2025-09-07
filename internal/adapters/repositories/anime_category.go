package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCategoryRepository interface {
	Save(animeCategory []entities.AnimeCategory) error
	GetByCategoryId(uuid.UUID) ([]GetByCategoryId, error)
	GetByAnimeIdAndCategoryIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategory, error)
}

type GormAnimeCategoryRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeCategoryRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeCategoryRepository {
	return &GormAnimeCategoryRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormAnimeCategoryRepository) Save(animeCategory []entities.AnimeCategory) error {
	if result := r.dbPrimary.Create(&animeCategory); result.Error != nil {
		return result.Error
	}
	return nil
}

type StudioDetail struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GetByCategoryId struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Episodes int            `json:"episodes"`
	Seasonal string         `json:"seasonal"`
	Year     string         `json:"year"`
	Image    string         `json:"image"`
	Studios  []StudioDetail `json:"studios"`
}

func (r GormAnimeCategoryRepository) GetByCategoryId(category_id uuid.UUID) ([]GetByCategoryId, error) {
	type TempResult struct {
		ID         uuid.UUID
		Name       string
		Episodes   int
		Seasonal   string
		Year       string
		Image      string
		StudioID   *uuid.UUID
		StudioName *string
	}
	var tempResults []TempResult
	sql := `
		SELECT
			a.id,
			a.name AS name,
			a.episodes,
			a.seasonal,
			a.year,
			a.image AS image,
			s.id AS studio_id,
			s.name AS studio_name
		FROM
			animes a
		JOIN
			anime_categories ac ON a.id = ac.anime_id
		LEFT JOIN
			anime_studios ast ON a.id = ast.anime_id
		LEFT JOIN
			studios s ON ast.studio_id = s.id
		WHERE
			ac.category_id = ?
			AND a.is_show = true
		ORDER BY
			a.aired_at DESC, a.id`

	result := r.dbReplica.Raw(sql, category_id).Scan(&tempResults)
	if result.Error != nil {
		return nil, result.Error
	}

	animeMap := make(map[uuid.UUID]*GetByCategoryId)
	var orderedAnimes []*GetByCategoryId

	for _, res := range tempResults {
		anime, exists := animeMap[res.ID]
		if !exists {
			anime = &GetByCategoryId{
				ID:       res.ID,
				Name:     res.Name,
				Episodes: res.Episodes,
				Seasonal: res.Seasonal,
				Year:     res.Year,
				Image:    res.Image,
				Studios:  []StudioDetail{},
			}
			animeMap[res.ID] = anime
			orderedAnimes = append(orderedAnimes, anime)
		}
		if res.StudioID != nil && res.StudioName != nil {
			isDuplicate := false
			for _, s := range anime.Studios {
				if s.ID == *res.StudioID {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				anime.Studios = append(anime.Studios, StudioDetail{ID: *res.StudioID, Name: *res.StudioName})
			}
		}
	}

	finalAnimes := make([]GetByCategoryId, len(orderedAnimes))
	for i, anime := range orderedAnimes {
		finalAnimes[i] = *anime
	}

	return finalAnimes, nil
}

func (r GormAnimeCategoryRepository) GetByAnimeIdAndCategoryIds(anime_id uuid.UUID, category_ids []uuid.UUID) ([]entities.AnimeCategory, error) {
	var categoryAnime []entities.AnimeCategory
	result := r.dbReplica.Where("anime_id = ?", anime_id).Where("category_id in (?)", category_ids).Find(&categoryAnime)
	if result.Error != nil {
		return nil, result.Error
	}
	return categoryAnime, nil
}
