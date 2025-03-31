package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type StudioService interface {
	GetAllStudio() ([]dtos.StudioListResponse, error)
}

type studioServiceImpl struct {
	repo repositories.StudioRepository
}

func NewStudioService(repo repositories.StudioRepository) StudioService {
	return &studioServiceImpl{repo: repo}
}

func (s studioServiceImpl) GetAllStudio() ([]dtos.StudioListResponse, error) {
	stduios, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	var studioResponse []dtos.StudioListResponse
	for _, stduio := range stduios {
		studioResponse = append(studioResponse, dtos.StudioListResponse{
			ID:          stduio.ID,
			Name:        stduio.Name,
			Image:       stduio.Image,
			Description: stduio.Description,
		})
	}

	return studioResponse, nil
}
