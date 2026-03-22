package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
)

type ProviderService interface {
	CreateProvider(req dtos.CreateProviderRequest) (*dtos.ProviderResponse, error)
	GetAllProviders() ([]dtos.ProviderResponse, error)
	AddProviderToAnime(req dtos.AddProviderToAnimeRequest) error
}

type providerServiceImpl struct {
	repo          repositories.ProviderRepository
	animeProvRepo repositories.AnimeProviderRepository
}

func NewProviderService(repo repositories.ProviderRepository, animeProvRepo repositories.AnimeProviderRepository) ProviderService {
	return &providerServiceImpl{repo: repo, animeProvRepo: animeProvRepo}
}

func (s *providerServiceImpl) CreateProvider(req dtos.CreateProviderRequest) (*dtos.ProviderResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	provider := &entities.Provider{
		ID:    id,
		Name:  req.Name,
		Image: req.Image,
	}

	saved, err := s.repo.Save(provider)
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	return &dtos.ProviderResponse{
		ID:    saved.ID,
		Name:  saved.Name,
		Image: saved.Image,
	}, nil
}

func (s *providerServiceImpl) AddProviderToAnime(req dtos.AddProviderToAnimeRequest) error {
	dups, err := s.animeProvRepo.GetByAnimeIdAndProviderIds(req.AnimeID, req.ProviderIDs)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	if len(dups) != 0 {
		return errs.NewBadRequestError("provider already mapped to this anime")
	}

	var records []entities.AnimeProvider
	for _, providerID := range req.ProviderIDs {
		id, err := uuid.NewV7()
		if err != nil {
			logs.Error(err.Error())
			return errs.NewUnexpectedError()
		}
		records = append(records, entities.AnimeProvider{
			ID:         id,
			AnimeID:    req.AnimeID,
			ProviderID: providerID,
		})
	}

	if err := s.animeProvRepo.Save(records); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s *providerServiceImpl) GetAllProviders() ([]dtos.ProviderResponse, error) {
	providers, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err.Error())
		return nil, errs.NewUnexpectedError()
	}

	result := make([]dtos.ProviderResponse, 0, len(providers))
	for _, p := range providers {
		result = append(result, dtos.ProviderResponse{
			ID:    p.ID,
			Name:  p.Name,
			Image: p.Image,
		})
	}
	return result, nil
}
