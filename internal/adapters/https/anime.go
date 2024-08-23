package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gofiber/fiber/v2"
)

type HttpAnimeHandler struct {
	service services.AnimeService
}

func NewHttpAnimeHandler(service services.AnimeService) *HttpAnimeHandler {
	return &HttpAnimeHandler{service: service}
}

func (h *HttpAnimeHandler) CreateAnime(c *fiber.Ctx) error {
	var anime entities.Anime
	if err := c.BodyParser(&anime); err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	if err := h.service.CreateAnime(anime); err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create anime success",
	})
}

func (h *HttpAnimeHandler) GetAnimeById(c *fiber.Ctx) error {
	animeId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	anime, err := h.service.GetAnimeById(uint(animeId))
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(anime)
}

func (h *HttpAnimeHandler) GetAnimeList(c *fiber.Ctx) error {
	query := c.Queries()
	animeQuery := dtos.AnimeQueryDTO{
		Seasonal: query["seasonal"],
		Year:     query["year"],
	}

	animes, err := h.service.GetAnimes(animeQuery)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(animes)
}

func (h *HttpAnimeHandler) UpdateAnime(c *fiber.Ctx) error {
	animeId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	animeUpdate := new(entities.Anime)
	if err := c.BodyParser(animeUpdate); err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	animeUpdate.ID = uint(animeId)
	err = h.service.UpdateAnime(*animeUpdate)
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update anime success",
	})
}

func (h *HttpAnimeHandler) DeleteAnime(c *fiber.Ctx) error {
	animeId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	if err = h.service.DeleteAnime(uint(animeId)); err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delete anime success",
	})
}

func (h *HttpAnimeHandler) GetAnimeByUserId(c *fiber.Ctx) error {
	user_id, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	animes, err := h.service.GetAnimeByUserId(uint(user_id))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(animes)
}
