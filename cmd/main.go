package main

import (
	"fmt"
	"log"
	"os"
	"time"

	adapters "github.com/Fourth1755/animap-go-api/internal/adapters/https"
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host         = "localhost"
	port         = 5432
	databaseName = "postgres"
	username     = "postgres"
	password     = "12131415"
)

var (
	animeHandler     *adapters.HttpAnimeHandler
	userHandler      *adapters.HttpUserHandler
	userAnimeHandler *adapters.HttpUserAnimeHandler
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, databaseName)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	print(db)

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entities.Anime{}, &entities.User{}, &entities.UserAnime{})
	animeRepo := repositories.NewGormAnimeRepository(db)
	animeService := services.NewAnimeService(animeRepo)
	animeHandler = adapters.NewHttpAnimeHandler(animeService)

	userRepo := repositories.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler = adapters.NewHttpUserHandler(userService)

	userAnimeRepo := repositories.NewGormUserAnimeRepository(db)
	userAnimeService := services.NewUserAnimeService(userAnimeRepo, animeRepo, userRepo)
	userAnimeHandler = adapters.NewHttpUserAnimeHandler(userAnimeService)
	InitRoutes()
}

func InitRoutes() {
	app := fiber.New()
	app.Post("register", userHandler.CreateUser)
	app.Post("login", userHandler.Login)

	app.Use("animes", jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	app.Post("animes", animeHandler.CreateAnime)
	app.Get("animes/:id", animeHandler.GetAnimeById)
	app.Get("animes", animeHandler.GetAnimeList)
	app.Put("animes/:id", animeHandler.UpdateAnime)
	app.Delete("animes/:id", animeHandler.DeleteAnime)

	app.Post("anime-list", userAnimeHandler.AddAnimeToList)
	app.Get("anime-list/:id", userAnimeHandler.GetAnimeByUserId)
	app.Listen(":8080")
}
