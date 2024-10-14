package main

import (
	"encoding/gob"
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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	animeHandler     *adapters.HttpAnimeHandler
	userHandler      *adapters.HttpUserHandler
	userAnimeHandler *adapters.HttpUserAnimeHandler
	categoryHandler  *adapters.HttpCategoryHandler
	songHandler      *adapters.HttpSongHandler
	artistHandler    *adapters.HttpArtistHandler
	authHandler      *adapters.HttpAuthHandler
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
	songService := services.NewSongService(songRepo, animeRepo, artistRepo, songArtistRepo)
	artistService := services.NewArtistService(artistRepo)

	auth, err := services.NewAuthenticator()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	//create handler
	animeHandler = adapters.NewHttpAnimeHandler(animeService)
	userHandler = adapters.NewHttpUserHandler(userService)
	userAnimeHandler = adapters.NewHttpUserAnimeHandler(userAnimeService)
	categoryHandler = adapters.NewHttpCategoryHandler(categoryService)
	songHandler = adapters.NewHttpSongHandler(songService)
	artistHandler = adapters.NewHttpArtistHandler(artistService)
	authHandler = adapters.NewHttpAuthHandler(auth)
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

func InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:8090", "http://localhost:3000"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	router.POST("register", userHandler.CreateUser)
	router.POST("login", userHandler.Login)

	// router.Use("animes", jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(os.GETenv("JWT_SECRET")),
	// }))
	//router.Use("animes", middleware.AuthRequired)

	//auth0
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// router.GET("/login", authHandler.Login())
	router.GET("/callback", authHandler.Callback)
	router.GET("/logout", authHandler.Logout)

	router.POST("animes", animeHandler.CreateAnime)
	router.GET("animes/:id", animeHandler.GetAnimeById)
	router.GET("animes", animeHandler.GetAnimeList)
	router.PUT("animes/:id", animeHandler.UpdateAnime)
	router.DELETE("animes/:id", animeHandler.DeleteAnime)
	router.POST("animes/category/:anime_id", animeHandler.AddCategoryToAnime)
	router.GET("animes/category/:category_id", animeHandler.GetAnimeByCategory)

	router.POST("user/anime-list", userAnimeHandler.AddAnimeToList)
	router.GET("user/anime-list/:uuid", userAnimeHandler.GetAnimeByUserId)
	//router.GET("anime-list/:id", userAnimeHandler.GETAnimeByUserId)

	router.POST("category", categoryHandler.CreateCategory)
	router.GET("category", categoryHandler.Getcategorise)
	router.GET("category/:id", categoryHandler.GetCategoryById)

	router.POST("songs", songHandler.CreateSong)
	router.GET("songs", songHandler.GetSongAll)
	router.GET("songs/:id", songHandler.GetSongById)
	router.PUT("songs/:id", songHandler.UpdateSong)
	router.DELETE("songs/:id", songHandler.DeleteSong)
	router.GET("songs/anime/:id", songHandler.GetSongByAnimeId)

	router.POST("artists", artistHandler.CreateArtist)
	router.GET("artists", artistHandler.GetArtistList)
	router.GET("artists/:id", artistHandler.GetArtistById)
	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
