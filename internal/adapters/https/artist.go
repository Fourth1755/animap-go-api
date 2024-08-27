package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gofiber/fiber/v2"
)

type HttpArtistHandler struct {
	service services.ArtistService
}

func NewHttpArtistHandler(service services.ArtistService) *HttpArtistHandler {
	return &HttpArtistHandler{service: service}
}

func (h *HttpArtistHandler) CreateArtist(c *fiber.Ctx) error {
	var artist *entities.Artist
	if err := c.BodyParser(&artist); err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	if err := h.service.CreateArtist(artist); err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create artist success",
	})
}

func (h *HttpArtistHandler) GetArtistById(c *fiber.Ctx) error {
	artistId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	Artist, err := h.service.GetArtistById(uint(artistId))
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(Artist)
}

func (h *HttpArtistHandler) GetArtistList(c *fiber.Ctx) error {
	artists, err := h.service.GetArtists()
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	return c.JSON(artists)
}
