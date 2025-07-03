package route

import (
	"net/http"
	"time"

	adapters "github.com/Fourth1755/animap-go-api/internal/adapters/https"
	"github.com/Fourth1755/animap-go-api/internal/middleware"
	"github.com/gin-contrib/cors"
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
	categoryUniverseHandler *adapters.HttpCategoryUniverseHandler,
	episodeHandler *adapters.HttpEpisodeHandler,
	characterHandler *adapters.HttpCharacterHandler) *gin.Engine {

	router := gin.Default()
	//router.Use(CORSMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("register", userHandler.CreateUser)
	router.POST("login", userHandler.Login)

	authorized := router.Group("/")
	authorized.Use(middleware.AuthRequired)
	{
		authorized.GET("user/user-info", userHandler.GetUserInfo)
		authorized.PATCH("user/user-info", userHandler.UpdateUserInfo)
	}

	router.GET("user/user-info/:uuid", userHandler.GetUserByUUID)

	// //auth0
	// gob.Register(map[string]interface{}{})

	// store := cookie.NewStore([]byte("secret"))
	// router.Use(sessions.Sessions("auth-session", store))

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
	router.GET("animes/studio/:studio_id", animeHandler.GetAnimeByStudio)

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
	router.GET("songs/artist/:id", songHandler.GetSongsByArtistId)

	router.POST("artists", artistHandler.CreateArtist)
	router.GET("artists", artistHandler.GetArtistList)
	router.GET("artists/:id", artistHandler.GetArtistById)
	router.PUT("artists/:id", artistHandler.UpdateArtist)

	router.GET("studios", studioHandler.GetAllStduio)

	router.GET("common/seasonal-year", commonHandler.GetSeasonalAndYear)

	router.POST("episodes", episodeHandler.CreateEpisode)
	router.GET("episodes/:anime_id", episodeHandler.GetEpisodesByAnimeId)
	router.PUT("episodes", episodeHandler.UpdateEpisode)
	router.POST("episodes/add-character", episodeHandler.AddCharactersToEpisode)

	router.POST("characters", characterHandler.CreateCharacter)
	router.GET("characters/:anime_id", characterHandler.GetCharacterByAnimeId)

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "http://localhost:3000" { // <-- dev origin
			c.Header("Access-Control-Allow-Origin", origin) // ต้องตรง origin
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers",
				"Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent) // header ถูกเซ็ตแล้ว
			return
		}

		c.Next()
	}
}
