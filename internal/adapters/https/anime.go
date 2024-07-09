package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if err := h.service.CreateAnime(anime); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create anime success",
	})
}

func (h *HttpAnimeHandler) GetAnimeById(c *fiber.Ctx) error {
	animeId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	anime, err := h.service.GetAnimeById(animeId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	animeDto := dtos.AnimeDTO{
		ID:       anime.ID,
		Name:     anime.Name,
		Episodes: anime.Episodes,
		Seasonal: anime.Seasonal,
		Year:     anime.Year,
	}
	return c.JSON(animeDto)
}

func (h *HttpAnimeHandler) GetAnimeList(c *fiber.Ctx) error {
	query := c.Queries()
	animeQuery := dtos.AnimeQueryDTO{
		Seasonal: query["seasonal"],
		Year:     query["year"],
	}

	animes, err := h.service.GetAnimes(animeQuery)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var animesDto []dtos.AnimeDTO
	for _, anime := range animes {
		animesDto = append(animesDto, dtos.AnimeDTO{
			ID:       anime.ID,
			Name:     anime.Name,
			Episodes: anime.Episodes,
			Seasonal: anime.Seasonal,
			Year:     anime.Year,
		})
	}
	return c.JSON(animesDto)
}
