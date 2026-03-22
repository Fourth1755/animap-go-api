package services

import (
	"sync"

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

	type result struct {
		items []dtos.SearchResultItem
		err   error
	}

	ch := make(chan result, 5)
	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		defer wg.Done()
		animes, err := s.animeRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			ch <- result{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(animes))
		for _, a := range animes {
			items = append(items, dtos.SearchResultItem{
				ID:    a.ID,
				Name:  a.Name,
				Image: a.Image,
				Type:  "anime",
			})
		}
		ch <- result{items: items}
	}()

	go func() {
		defer wg.Done()
		studios, err := s.studioRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			ch <- result{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(studios))
		for _, s := range studios {
			items = append(items, dtos.SearchResultItem{
				ID:    s.ID,
				Name:  s.Name,
				Image: s.Image,
				Type:  "studio",
			})
		}
		ch <- result{items: items}
	}()

	go func() {
		defer wg.Done()
		songs, err := s.songRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			ch <- result{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(songs))
		for _, s := range songs {
			items = append(items, dtos.SearchResultItem{
				ID:    s.ID,
				Name:  s.Name,
				Image: s.Image,
				Type:  "song",
			})
		}
		ch <- result{items: items}
	}()

	go func() {
		defer wg.Done()
		categories, err := s.categoryUniverseRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			ch <- result{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(categories))
		for _, c := range categories {
			items = append(items, dtos.SearchResultItem{
				ID:    c.ID,
				Name:  c.Name,
				Image: c.Image,
				Type:  "category_universe",
			})
		}
		ch <- result{items: items}
	}()

	go func() {
		defer wg.Done()
		characters, err := s.characterRepo.Search(keyword, searchLimit)
		if err != nil {
			logs.Error(err.Error())
			ch <- result{err: err}
			return
		}
		items := make([]dtos.SearchResultItem, 0, len(characters))
		for _, c := range characters {
			items = append(items, dtos.SearchResultItem{
				ID:    c.ID,
				Name:  c.Name,
				Image: c.Image,
				Type:  "character",
			})
		}
		ch <- result{items: items}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	var allItems []dtos.SearchResultItem
	for r := range ch {
		if r.err != nil {
			return nil, errs.NewUnexpectedError()
		}
		allItems = append(allItems, r.items...)
	}

	if allItems == nil {
		allItems = []dtos.SearchResultItem{}
	}

	return &dtos.SearchResponse{Results: allItems}, nil
}
