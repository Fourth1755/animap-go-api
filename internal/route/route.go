package route

import (
	"encoding/gob"

	adapters "github.com/Fourth1755/animap-go-api/internal/adapters/https"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRoutes(
	animeHandler *adapters.HttpAnimeHandler,
	userHandler *adapters.HttpUserHandler,
	myAnimeHandler *adapters.HttpMyAnimeHandler,
	categoryHandler *adapters.HttpCategoryHandler,
	songHandler *adapters.HttpSongHandler,
	artistHandler *adapters.HttpArtistHandler,
	studioHandler *adapters.HttpStduioHandler,
	commonHandler *adapters.HttpCommonHandler,
	categoryUniverseHandler *adapters.HttpCategoryUniverseHandler) *gin.Engine {

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
	router.GET("user/user-info", userHandler.GetUserInfo)
	router.PATCH("user/user-info", userHandler.UpdateUserInfo)

	// router.Use("animes", jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(os.GETenv("JWT_SECRET")),
	// }))
	//router.Use("animes", middleware.AuthRequired)

	//auth0
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// router.GET("/login", authHandler.Login())
	router.POST("animes", animeHandler.CreateAnime)
	router.GET("animes/:id", animeHandler.GetAnimeById)
	router.GET("animes", animeHandler.GetAnimeList)
	router.PUT("animes/:id", animeHandler.UpdateAnime)
	router.DELETE("animes/:id", animeHandler.DeleteAnime)
	router.PUT("animes/category/edit-category-anime", animeHandler.AddCategoryToAnime)
	router.PUT("animes/category-universe/edit-category-universe-anime", animeHandler.AddCategoryUniverseToAnime)
	router.GET("animes/category/:category_id", animeHandler.GetAnimeByCategory)
	router.GET("animes/category-universe/:category_id", animeHandler.GetAnimeByCategoryUniverse)
	router.POST("animes/seasonal-year", animeHandler.GetAnimeBySeasonalAndYear)

	router.POST("my-anime", myAnimeHandler.AddAnimeToList)
	router.GET("my-anime/:uuid", myAnimeHandler.GetAnimeByUserId)
	router.GET("my-anime/anime-year-list/:uuid", myAnimeHandler.GetMyAnimeYearByUserId)
	router.GET("my-anime/top-anime/:uuid", myAnimeHandler.GetMyTopAnimeByUserId)
	router.PATCH("my-anime/top-anime", myAnimeHandler.UpdateMyTopAnime)

	//router.GET("anime-list/:id", userAnimeHandler.GETAnimeByUserId)

	router.POST("category", categoryHandler.CreateCategory)
	router.GET("category", categoryHandler.Getcategorise)
	router.GET("category/:id", categoryHandler.GetCategoryById)

	router.GET("category-universe", categoryUniverseHandler.Getcategorise)

	router.POST("songs", songHandler.CreateSong)
	router.GET("songs", songHandler.GetSongAll)
	router.GET("songs/:id", songHandler.GetSongById)
	router.PUT("songs/:id", songHandler.UpdateSong)
	router.DELETE("songs/:id", songHandler.DeleteSong)
	router.GET("songs/anime/:id", songHandler.GetSongByAnimeId)
	router.POST("songs/channel", songHandler.CreateSongChannel)

	router.POST("artists", artistHandler.CreateArtist)
	router.GET("artists", artistHandler.GetArtistList)
	router.GET("artists/:id", artistHandler.GetArtistById)

	router.GET("studios", studioHandler.GetAllStduio)

	router.GET("common/seasonal-year", commonHandler.GetSeasonalAndYear)
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
