package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.service.AddAnimeToList(userAnime); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Add anime to list success.",
	})
}

func (h *HttpUserAnimeHandler) GetAnimeByUserId(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userAnimes, err := h.service.GetAnimeByUserId(uint(userId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var animeList []dtos.UserAnimeListDTO
	for _, useranime := range userAnimes {
		animeList = append(animeList, dtos.UserAnimeListDTO{
			AnimeID:     useranime.AnimeID,
			AnimeName:   useranime.Anime.Name,
			Score:       useranime.Score,
			Description: useranime.Anime.Description,
			Episodes:    useranime.Anime.Description,
			Image:       useranime.Anime.Image,
			Status:      useranime.Status,
			WatchAt:     useranime.WatchAt,
			CreatedAt:   useranime.CreatedAt,
		})
	}
	return c.JSON(animeList)
}
