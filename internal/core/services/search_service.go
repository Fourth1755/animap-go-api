package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

const searchLimit = 10

type SearchService interface {
	Search(keyword string) (*dtos.SearchResponse, error)
}

type searchServiceImpl struct {
	animeRepo            repositories.AnimeRepository
	studioRepo           repositories.StudioRepository
	songRepo             repositories.SongRepository
	categoryUniverseRepo repositories.CategoryUniverseRepository
	characterRepo        repositories.CharacterRepository
}

func NewSearchService(
	animeRepo repositories.AnimeRepository,
	studioRepo repositories.StudioRepository,
	songRepo repositories.SongRepository,
	categoryUniverseRepo repositories.CategoryUniverseRepository,
	characterRepo repositories.CharacterRepository,
) SearchService {
	return &searchServiceImpl{
		animeRepo:            animeRepo,
		studioRepo:           studioRepo,
		songRepo:             songRepo,
		categoryUniverseRepo: categoryUniverseRepo,
		characterRepo:        characterRepo,
	}
}

func (s *searchServiceImpl) Search(keyword string) (*dtos.SearchResponse, error) {
	if keyword == "" {
		return &dtos.SearchResponse{Results: []dtos.SearchResultItem{}}, nil
	}

	type bucket struct {
		items []dtos.SearchResultItem
		err   error
	}

	animes := make(chan bucket, 1)
	categoryUniverses := make(chan bucket, 1)
	songs := make(chan bucket, 1)
	characters := make(chan bucket, 1)
	studios := make(chan bucket, 1)

	go func() {
		list, err := s.animeRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			animes <- bucket{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(list))
		for _, a := range list {
			items = append(items, dtos.SearchResultItem{ID: a.ID, Name: a.Name, Image: a.Image, Type: "anime"})
		}
		animes <- bucket{items: items}
	}()

	go func() {
		list, err := s.categoryUniverseRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			categoryUniverses <- bucket{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(list))
		for _, c := range list {
			items = append(items, dtos.SearchResultItem{ID: c.ID, Name: c.Name, Type: "category_universe"})
		}
		categoryUniverses <- bucket{items: items}
	}()

	go func() {
		list, err := s.songRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			songs <- bucket{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(list))
		for _, s := range list {
			items = append(items, dtos.SearchResultItem{ID: s.ID, Name: s.Name, Image: s.Image, Type: "song"})
		}
		songs <- bucket{items: items}
	}()

	go func() {
		list, err := s.characterRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			characters <- bucket{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(list))
		for _, c := range list {
			items = append(items, dtos.SearchResultItem{ID: c.ID, Name: c.Name, Image: c.Image, Type: "character"})
		}
		characters <- bucket{items: items}
	}()

	go func() {
		list, err := s.studioRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			studios <- bucket{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(list))
		for _, s := range list {
			items = append(items, dtos.SearchResultItem{ID: s.ID, Name: s.Name, Image: s.Image, Type: "studio"})
		}
		studios <- bucket{items: items}
	}()

	// collect in priority order: anime → category_universe → song → character → studio
	ordered := make([]bucket, 5)
	ordered[0] = <-animes
	ordered[1] = <-categoryUniverses
	ordered[2] = <-songs
	ordered[3] = <-characters
	ordered[4] = <-studios

	result := make([]dtos.SearchResultItem, 0, searchLimit)
	for _, b := range ordered {
		if b.err != nil {
			return nil, errs.NewUnexpectedError()
		}
		remaining := searchLimit - len(result)
		if remaining <= 0 {
			break
		}
		if len(b.items) > remaining {
			result = append(result, b.items[:remaining]...)
		} else {
			result = append(result, b.items...)
		}
	}

	return &dtos.SearchResponse{Results: result}, nil
}
