package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	adapters "github.com/Fourth1755/animap-go-api/internal/adapters/https"
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/config"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/logs"
	"github.com/Fourth1755/animap-go-api/internal/route"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db := InitDatabase()
	// init config
	configService := config.NewConfigService()

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
	animeCategorryUnivserseRepo := repositories.NewGormAnimeCategoryUniverseRepository(db)
	categoryUniverseRepo := repositories.NewGormCategoryUniverseRepository(db)

	//create service
	userService := services.NewUserService(userRepo)
	myAnimeService := services.NewMyAnimeService(userAnimeRepo, animeRepo, userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	songService := services.NewSongService(songRepo, animeRepo, artistRepo, songArtistRepo, songChannelRepo)
	artistService := services.NewArtistService(artistRepo)
	studioService := services.NewStudioService(studioRepo)
	animeService := services.NewAnimeService(animeRepo, userRepo, animeCategoryRepo, animeStudioRepo, songRepo, categoryRepo, animeCategorryUnivserseRepo)
	commonService := services.NewCommonService(configService)
	categoryUniverseService := services.NewCategoryUniverseService(categoryUniverseRepo)

	//create handler
	animeHandler := adapters.NewHttpAnimeHandler(animeService)
	userHandler := adapters.NewHttpUserHandler(userService)
	myAnimeHandler := adapters.NewHttpMyAnimeHandler(myAnimeService)
	categoryHandler := adapters.NewHttpCategoryHandler(categoryService)
	songHandler := adapters.NewHttpSongHandler(songService)
	artistHandler := adapters.NewHttpArtistHandler(artistService)
	studioHandler := adapters.NewHttpStduioHandler(studioService)
	commonHandler := adapters.NewHttpCommonHandler(commonService)
	categoryUniverseHandler := adapters.NewHttpCategoryUniverseHandler(categoryUniverseService)

	rtr := route.InitRoutes(animeHandler,
		userHandler,
		myAnimeHandler,
		categoryHandler,
		songHandler,
		artistHandler,
		studioHandler,
		commonHandler,
		categoryUniverseHandler)

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
		&entities.CategoryUniverse{},
		&entities.AnimeCategoryUniverse{},
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
