package adapters

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
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
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	if err := h.service.AddAnimeToCategory(animeCategory); err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Add anime to category success.",
	})
}
