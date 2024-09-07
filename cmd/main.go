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
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	animeHandler         *adapters.HttpAnimeHandler
	userHandler          *adapters.HttpUserHandler
	userAnimeHandler     *adapters.HttpUserAnimeHandler
	categoryHandler      *adapters.HttpCategoryHandler
	animeCategoryHandler *adapters.HttpAnimeCategoryHandler
	songHandler          *adapters.HttpSongHandler
	artistHandler        *adapters.HttpArtistHandler
)

func main() {
	db := InitDatabase()
	//create repository
	animeRepo := repositories.NewGormAnimeRepository(db)
	userRepo := repositories.NewGormUserRepository(db)
	userAnimeRepo := repositories.NewGormUserAnimeRepository(db)
	categoryRepo := repositories.NewGormCategoryRepository(db)
	animeCategoryRepo := repositories.NewGormAnimeCategoryRepository(db)
	songRepo := repositories.NewGormSongRepository(db)
	artistRepo := repositories.NewGormArtistRepository(db)
	songArtistRepo := repositories.NewGormSongArtistRepository(db)

	//create service
	animeService := services.NewAnimeService(animeRepo, userRepo, animeCategoryRepo)
	userService := services.NewUserService(userRepo)
	userAnimeService := services.NewUserAnimeService(userAnimeRepo, animeRepo, userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	animeCategoryService := services.NewAnimeCategoryService(animeCategoryRepo)
	songService := services.NewSongService(songRepo, animeRepo, artistRepo, songArtistRepo)
	artistService := services.NewArtistService(artistRepo)

	//create handler
	animeHandler = adapters.NewHttpAnimeHandler(animeService)
	userHandler = adapters.NewHttpUserHandler(userService)
	userAnimeHandler = adapters.NewHttpUserAnimeHandler(userAnimeService)
	categoryHandler = adapters.NewHttpCategoryHandler(categoryService)
	animeCategoryHandler = adapters.NewHttpAnimeCategoryHandler(animeCategoryService)
	songHandler = adapters.NewHttpSongHandler(songService)
	artistHandler = adapters.NewHttpArtistHandler(artistService)
	InitRoutes()
}

func InitDatabase() *gorm.DB {
	initConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.databaseName"))

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
	if err != nil {
		panic("failed to connect database")
	}
	logs.Error("failed to connect database")

	db.AutoMigrate(
		&entities.Anime{},
		&entities.User{},
		&entities.UserAnime{},
		&entities.Category{},
		&entities.AnimeCategory{},
		&entities.Song{},
		&entities.SongChannel{},
		&entities.Artist{},
		&entities.SongArtist{},
	)

	return db
}
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitRoutes() {
	app := fiber.New()
	app.Use(cors.New())
	app.Post("register", userHandler.CreateUser)
	app.Post("login", userHandler.Login)

	// app.Use("animes", jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(os.Getenv("JWT_SECRET")),
	// }))
	//app.Use("animes", middleware.AuthRequired)
	app.Post("animes", animeHandler.CreateAnime)
	app.Get("animes/:id", animeHandler.GetAnimeById)
	app.Get("animes", animeHandler.GetAnimeList)
	app.Put("animes/:id", animeHandler.UpdateAnime)
	app.Delete("animes/:id", animeHandler.DeleteAnime)
	app.Get("anime-list/:user_id", animeHandler.GetAnimeByUserId)
	app.Post("animes/category", animeCategoryHandler.AddAnimeToCategory)
	app.Get("animes/category/:category_id", animeHandler.GetAnimeByCategory)

	app.Post("anime-list", userAnimeHandler.AddAnimeToList)
	//app.Get("anime-list/:id", userAnimeHandler.GetAnimeByUserId)

	app.Post("category", categoryHandler.CreateCategory)
	app.Get("category", categoryHandler.Getcategorise)
	app.Get("category/:id", categoryHandler.GetCategoryById)

	app.Post("songs", songHandler.CreateSong)
	app.Get("songs", songHandler.GetSongAll)
	app.Get("songs/:id", songHandler.GetSongById)
	app.Put("songs/:id", songHandler.UpdateSong)
	app.Delete("songs/:id", songHandler.DeleteSong)
	app.Get("songs/anime/:id", songHandler.GetSongByAnimeId)

	app.Post("artists", artistHandler.CreateArtist)
	app.Get("artists", artistHandler.GetArtistList)
	app.Get("artists/:id", artistHandler.GetArtistById)
	app.Listen(":8080")
}
