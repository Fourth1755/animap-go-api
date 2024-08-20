package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gofiber/fiber/v2"
)

type HttpUserAnimeHandler struct {
	service services.UserAnimeService
}

func NewHttpUserAnimeHandler(service services.UserAnimeService) *HttpUserAnimeHandler {
	return &HttpUserAnimeHandler{service: service}
}

func (h *HttpUserAnimeHandler) AddAnimeToList(c *fiber.Ctx) error {
	userAnime := new(entities.UserAnime)
	if err := c.BodyParser(userAnime); err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	if err := h.service.AddAnimeToList(userAnime); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Add anime to list success.",
	})
}

func (h *HttpUserAnimeHandler) GetAnimeByUserId(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	animeList, err := h.service.GetAnimeByUserId(uint(userId))
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(animeList)
}
