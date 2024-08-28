package adapters

import (
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
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
	var song dtos.CreateSongRequest
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

func (h *HttpSongHandler) UpdateSong(c *fiber.Ctx) error {
	songId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	song := new(entities.Song)
	if err := c.BodyParser(&song); err != nil {
		return handleError(c, err)
	}
	song.ID = uint(songId)
	if err := h.service.UpdateSong(song); err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Update song success",
	})
}

func (h *HttpSongHandler) DeleteSong(c *fiber.Ctx) error {
	songId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	if err := h.service.DeleteSong(uint(songId)); err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Delete song success",
	})
}

func (h *HttpSongHandler) GetSongByAnimeId(c *fiber.Ctx) error {
	animeId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	songs, err := h.service.GetSongByAnimeId(uint(animeId))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(songs)
}
