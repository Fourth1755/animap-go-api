package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/gofiber/fiber/v2"
)

type HttpCategoryHandler struct {
	service services.CategoryService
}

func NewHttpCategoryHandler(service services.CategoryService) *HttpCategoryHandler {
	return &HttpCategoryHandler{service: service}
}

func (h *HttpCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	category := new(entities.Category)
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	if err := h.service.CreateCategory(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create category success",
	})
}

func (h *HttpCategoryHandler) Getcategorise(c *fiber.Ctx) error {
	category, err := h.service.Getcategorise()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

func (h *HttpCategoryHandler) GetCategoryById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	category, err := h.service.GetCategoryById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}