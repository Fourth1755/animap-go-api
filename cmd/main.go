package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	adapters "github.com/Fourth1755/animap-go-api/internal/adapters/https"
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	animeHandler    *adapters.HttpAnimeHandler
	userHandler     *adapters.HttpUserHandler
	myAnimeHandler  *adapters.HttpMyAnimeHandler
	categoryHandler *adapters.HttpCategoryHandler
	songHandler     *adapters.HttpSongHandler
	artistHandler   *adapters.HttpArtistHandler
	studioHandler   *adapters.HttpStduioHandler
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
	songChannelRepo := repositories.NewGormSongChannelRepository(db)
	studioRepo := repositories.NewGormStudioRepository(db)
	animeStudioRepo := repositories.NewGormAnimeStudioRepository(db)

	//create service
	
	userService := services.NewUserService(userRepo)
	myAnimeService := services.NewMyAnimeService(userAnimeRepo, animeRepo, userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	songService := services.NewSongService(songRepo, animeRepo, artistRepo, songArtistRepo, songChannelRepo)
	artistService := services.NewArtistService(artistRepo)
	studioService := services.NewStudioService(studioRepo)
	animeService := services.NewAnimeService(animeRepo, userRepo, animeCategoryRepo, animeStudioRepo, songRepo)

	//create handler
	animeHandler = adapters.NewHttpAnimeHandler(animeService)
	userHandler = adapters.NewHttpUserHandler(userService)
	myAnimeHandler = adapters.NewHttpMyAnimeHandler(myAnimeService)
	categoryHandler = adapters.NewHttpCategoryHandler(categoryService)
	songHandler = adapters.NewHttpSongHandler(songService)
	artistHandler = adapters.NewHttpArtistHandler(artistService)
	studioHandler = adapters.NewHttpStduioHandler(studioService)
	rtr := InitRoutes()

	log.Print("Server listening on http://localhost:8080/")
	if err := http.ListenAndServe("0.0.0.0:8080", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
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
		&entities.Studio{},
		&entities.AnimeStudio{},
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
