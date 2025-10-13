package repositories

import (
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeRepository interface {
	Save(anime entities.Anime) (*entities.Anime, error)
	GetById(id uuid.UUID) (*entities.Anime, error)
	GetByIds(ids []uuid.UUID) ([]entities.Anime, error)
	GetAll(query dtos.AnimeQueryDTO) ([]entities.Anime, error)
	Update(anime *entities.Anime) error
	Delete(id uuid.UUID) error
	GetByUserId(user_id uuid.UUID) ([]entities.UserAnime, error)
	GetBySeasonalAndYear(request dtos.GetAnimeBySeasonAndYearRequest) ([]entities.Anime, error)
	UpdateIsCreateEpisode(animeId uuid.UUID) error
	UpdadteImage(image string, myAnimeListId int) error
	GetByMyAnimeListId(id int) (*entities.Anime, error)
	UpdadteAiredAt(airedAt time.Time, myAnimeListId int) error
}

type GormAnimeRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeRepository {
	return &GormAnimeRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormAnimeRepository) Save(anime entities.Anime) (*entities.Anime, error) {
	if result := r.dbPrimary.Create(&anime); result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (r *GormAnimeRepository) GetById(id uuid.UUID) (*entities.Anime, error) {
	var anime entities.Anime
	if result := r.dbReplica.
		Preload("Songs").
		Preload("Categories").
		Preload("Studios").
		Preload("CategoryUniverses").
		First(&anime, id); result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (r *GormAnimeRepository) GetByIds(ids []uuid.UUID) ([]entities.Anime, error) {
	var animes []entities.Anime
	if result := r.dbReplica.
		Preload("Songs").
		Preload("Categories").
		Preload("Studios").
		Preload("CategoryUniverses").
		Order("aired_at desc").
		Where("id in (?)", ids).Find(&animes); result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}

func (r *GormAnimeRepository) GetAll(query dtos.AnimeQueryDTO) ([]entities.Anime, error) {
	var animes []entities.Anime
	db := r.dbReplica.Model(&entities.Anime{})

	if query.Seasonal != "" {
		db = db.Where("seasonal = ?", query.Seasonal)
	}

	if query.Year != "" {
		db = db.Where("year = ?", query.Year)
	}

	if query.Name != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?) OR LOWER(name_thai) LIKE LOWER(?) OR LOWER(name_english) LIKE LOWER(?) ", "%"+query.Name+"%", "%"+query.Name+"%", "%"+query.Name+"%")
	}

	if query.StudioID != "" {
		db = db.Joins("JOIN anime_studios ON anime_studios.anime_id = animes.id").Where("anime_studios.studio_id = ?", query.StudioID)
	}

	if query.CategoryID != "" {
		db = db.Joins("JOIN anime_categories ON anime_categories.anime_id = animes.id").Where("anime_categories.category_id = ?", query.CategoryID)
	}

	if query.SortBy != "" {
		orderBy := "asc"
		if query.OrderBy != "" {
			orderBy = query.OrderBy
		}
		db = db.Order(query.SortBy + " " + orderBy)
	}

	result := db.Find(&animes)

	if result.Error != nil {
		return nil, result.Error
	}

	return animes, nil
}

func (r *GormAnimeRepository) Update(anime *entities.Anime) error {
	result := r.dbPrimary.Model(&anime).Updates(anime)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) UpdateIsCreateEpisode(animeId uuid.UUID) error {
	var animes []entities.Anime
	result := r.dbPrimary.Raw("UPDATE animes SET is_create_episode = ? WHERE id = ? ", true, animeId).Scan(&animes)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) Delete(id uuid.UUID) error {
	var anime entities.Anime
	result := r.dbPrimary.Delete(&anime, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) GetByUserId(id uuid.UUID) ([]entities.UserAnime, error) {
	var animes []entities.UserAnime
	result := r.dbReplica.Preload("Anime").Where("user_id = ?", id).Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}
	return animes, nil
}

func (r *GormAnimeRepository) GetBySeasonalAndYear(request dtos.GetAnimeBySeasonAndYearRequest) ([]entities.Anime, error) {
	var animes []entities.Anime
	result := r.dbReplica.Preload("Studios").Where("seasonal = ?", request.Seasonal).Where("year = ?", request.Year).Find(&animes)
	if result.Error != nil {
		return nil, result.Error
	}

	return animes, nil
}

func (r *GormAnimeRepository) UpdadteImage(image string, myAnimeListId int) error {
	var animes []entities.Anime
	result := r.dbPrimary.Raw("UPDATE animes SET image = ? WHERE my_anime_list_id = ? ", image, myAnimeListId).Scan(&animes)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeRepository) UpdadteAiredAt(airedAt time.Time, myAnimeListId int) error {
	var animes []entities.Anime
	result := r.dbPrimary.Raw("UPDATE animes SET aired_at = ? WHERE my_anime_list_id = ? ", airedAt, myAnimeListId).Scan(&animes)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *GormAnimeRepository) GetByMyAnimeListId(id int) (*entities.Anime, error) {
	var anime entities.Anime
	if result := r.dbReplica.Where("my_anime_list_id = ?", id).
		First(&anime); result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}
