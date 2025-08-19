package route

import (
	"net/http"
	"time"

	adapters "github.com/Fourth1755/animap-go-api/internal/adapters/https"
	"github.com/Fourth1755/animap-go-api/internal/adapters/websocket"
	"github.com/Fourth1755/animap-go-api/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	AnimeHandler            *adapters.HttpAnimeHandler
	UserHandler             *adapters.HttpUserHandler
	MyAnimeHandler          *adapters.HttpMyAnimeHandler
	CategoryHandler         *adapters.HttpCategoryHandler
	SongHandler             *adapters.HttpSongHandler
	ArtistHandler           *adapters.HttpArtistHandler
	StudioHandler           *adapters.HttpStduioHandler
	CommonHandler           *adapters.HttpCommonHandler
	CategoryUniverseHandler *adapters.HttpCategoryUniverseHandler
	EpisodeHandler          *adapters.HttpEpisodeHandler
	CharacterHandler        *adapters.HttpCharacterHandler
}

func InitRoutes(https HttpHandler) *gin.Engine {

	router := gin.Default()
	//router.Use(CORSMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8090"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("register", https.UserHandler.CreateUser)
	router.POST("login", https.UserHandler.Login)

	authorized := router.Group("/")
	authorized.Use(middleware.AuthRequired)
	{
		authorized.POST("logout", https.UserHandler.Logout)

		authorized.GET("user/user-info", https.UserHandler.GetUserInfo)
		authorized.PATCH("user/user-info", https.UserHandler.UpdateUserInfo)
		authorized.POST("user/user-info/presign-url-avatar", https.UserHandler.GetPresignedURLAvatar)

	}

	router.GET("user/user-info/:uuid", https.UserHandler.GetUserByUUID)

	// //auth0
	// gob.Register(map[string]interface{}{})

	// store := cookie.NewStore([]byte("secret"))
	// router.Use(sessions.Sessions("auth-session", store))

	// router.GET("/login", authHandler.Login())
	router.POST("animes", https.AnimeHandler.CreateAnime)
	router.GET("animes/:id", https.AnimeHandler.GetAnimeById)
	router.GET("animes", https.AnimeHandler.GetAnimeList)
	router.PUT("animes/:id", https.AnimeHandler.UpdateAnime)
	router.DELETE("animes/:id", https.AnimeHandler.DeleteAnime)
	router.PUT("animes/category/edit-category-anime", https.AnimeHandler.AddCategoryToAnime)
	router.PUT("animes/category-universe/edit-category-universe-anime", https.AnimeHandler.AddCategoryUniverseToAnime)
	router.GET("animes/category/:category_id", https.AnimeHandler.GetAnimeByCategory)
	router.GET("animes/category-universe/:category_id", https.AnimeHandler.GetAnimeByCategoryUniverse)
	router.POST("animes/seasonal-year", https.AnimeHandler.GetAnimeBySeasonalAndYear)
	router.GET("animes/studio/:studio_id", https.AnimeHandler.GetAnimeByStudio)

	// migrate anime
	router.POST("migrate/animes", https.AnimeHandler.MigrateAnime)

	{
		authorized.POST("my-anime", https.MyAnimeHandler.AddAnimeToList)
	}
	//router.POST("my-anime", myAnimeHandler.AddAnimeToList)
	router.GET("my-anime/:uuid", https.MyAnimeHandler.GetAnimeByUserId)
	router.GET("my-anime/anime-year-list/:uuid", https.MyAnimeHandler.GetMyAnimeYearByUserId)
	router.GET("my-anime/top-anime/:uuid", https.MyAnimeHandler.GetMyTopAnimeByUserId)
	router.PATCH("my-anime/top-anime", https.MyAnimeHandler.UpdateMyTopAnime)

	//router.GET("anime-list/:id", userAnimeHandler.GETAnimeByUserId)

	router.POST("category", https.CategoryHandler.CreateCategory)
	router.GET("category", https.CategoryHandler.Getcategorise)
	router.GET("category/:id", https.CategoryHandler.GetCategoryById)

	router.GET("category-universe", https.CategoryUniverseHandler.Getcategorise)

	router.POST("songs", https.SongHandler.CreateSong)
	router.GET("songs", https.SongHandler.GetSongAll)
	router.GET("songs/:id", https.SongHandler.GetSongById)
	router.PUT("songs/:id", https.SongHandler.UpdateSong)
	router.DELETE("songs/:id", https.SongHandler.DeleteSong)
	router.GET("songs/anime/:id", https.SongHandler.GetSongByAnimeId)
	router.POST("songs/channel", https.SongHandler.CreateSongChannel)
	router.GET("songs/artist/:id", https.SongHandler.GetSongsByArtistId)

	router.POST("artists", https.ArtistHandler.CreateArtist)
	router.GET("artists", https.ArtistHandler.GetArtistList)
	router.GET("artists/:id", https.ArtistHandler.GetArtistById)
	router.PUT("artists/:id", https.ArtistHandler.UpdateArtist)

	router.GET("studios", https.StudioHandler.GetAllStduio)

	router.GET("common/seasonal-year", https.CommonHandler.GetSeasonalAndYear)

	router.POST("episodes", https.EpisodeHandler.CreateEpisode)
	router.GET("episodes/:anime_id", https.EpisodeHandler.GetEpisodesByAnimeId)
	router.PUT("episodes", https.EpisodeHandler.UpdateEpisode)
	router.POST("episodes/add-character", https.EpisodeHandler.AddCharactersToEpisode)

	router.POST("characters", https.CharacterHandler.CreateCharacter)
	router.GET("characters/:anime_id", https.CharacterHandler.GetCharacterByAnimeId)

	// websocket
	hub := websocket.H
	go hub.Run()
	{
		router.GET("rooms/:roomId/ws", websocket.HandleWebSocket)
	}

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
