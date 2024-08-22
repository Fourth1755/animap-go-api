package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gofiber/fiber/v2"
)

type HttpSongHandler struct {
	service services.SongService
}

func NewHttpSongHandler(service services.SongService) *HttpSongHandler {
	return &HttpSongHandler{service: service}
}

func (h *HttpSongHandler) CreateSong(c *fiber.Ctx) error {
	var song entities.Song
	if err := c.BodyParser(&song); err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	if err := h.service.CreateSong(&song); err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create song success",
	})
}

func (h *HttpSongHandler) GetSongById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	song, err := h.service.GetSongById(uint(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(song)
}

func (h *HttpSongHandler) GetSongAll(c *fiber.Ctx) error {
	song, err := h.service.GetAllSongs()
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(song)
}
