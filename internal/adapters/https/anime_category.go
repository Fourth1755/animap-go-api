package adapters

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/gofiber/fiber/v2"
)

type HttpAnimeCategoryHandler struct {
	service services.AnimeCategoryService
}

func NewHttpAnimeCategoryHandler(service services.AnimeCategoryService) *HttpAnimeCategoryHandler {
	return &HttpAnimeCategoryHandler{service: service}
}

func (h HttpAnimeCategoryHandler) AddAnimeToCategory(c *fiber.Ctx) error {
	var animeCategory *entities.AnimeCategory
	if err := c.BodyParser(&animeCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if err := h.service.AddAnimeToCategory(animeCategory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Add anime to category success.",
	})
}
