package services

import (
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArtistService interface {
	CreateArtist(Artist *entities.Artist) error
	UpdateArtist(artist *entities.Artist) error
	GetArtists() (*dtos.GetArtistsResponse, error)
	GetArtistById(id uuid.UUID) (*entities.Artist, error)
}

type ArtistServiceImpl struct {
	repo repositories.ArtistRepository
}

func NewArtistService(repo repositories.ArtistRepository) ArtistService {
	return &ArtistServiceImpl{repo: repo}
}

func (s ArtistServiceImpl) CreateArtist(artist *entities.Artist) error {
	artistId, err := uuid.NewV7()
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	artist.ID = artistId
	if err := s.repo.Save(artist); err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s ArtistServiceImpl) UpdateArtist(artist *entities.Artist) error {
	_, err := s.repo.GetById(artist.ID)
	if err != nil {
		logs.Error(err.Error())
		if err == gorm.ErrRecordNotFound {
			return errs.NewNotFoundError("Artist not found")
		}
		return errs.NewUnexpectedError()
	}

	if err := s.repo.Update(artist); err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

func (s ArtistServiceImpl) GetArtists() (*dtos.GetArtistsResponse, error) {
	artists, err := s.repo.GetAll()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	return &dtos.GetArtistsResponse{
		Artists: artists,
	}, nil
}

func (s ArtistServiceImpl) GetArtistById(id uuid.UUID) (*entities.Artist, error) {
	artist, err := s.repo.GetById(id)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewNotFoundError("Not found artist")
	}
	return artist, nil
}
