package main

import (
	"log"
	"net/http"

	aws_adapter "github.com/Fourth1755/animap-go-api/internal/adapters/aws"
	adapters "github.com/Fourth1755/animap-go-api/internal/adapters/https"
	"github.com/Fourth1755/animap-go-api/internal/adapters/repositories"
	"github.com/Fourth1755/animap-go-api/internal/core/config"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/database"
	"github.com/Fourth1755/animap-go-api/internal/route"
)

func main() {
	// init config
	configService := config.NewConfigService()

	// init database
	dbPrimary, dbReplica := database.InitDatabase(configService)

	// init aws
	awsAdapter, err := aws_adapter.NewAwsAdapter(configService.GetAWS())
	if err != nil {
		log.Fatalf("Could not create aws adapter: %v", err)
	}
	s3Service := aws_adapter.NewS3Service(awsAdapter)

	//create repository
	animeRepo := repositories.NewGormAnimeRepository(dbPrimary, dbReplica)
	userRepo := repositories.NewGormUserRepository(dbPrimary, dbReplica)
	userAnimeRepo := repositories.NewGormUserAnimeRepository(dbPrimary, dbReplica)
	categoryRepo := repositories.NewGormCategoryRepository(dbPrimary, dbReplica)
	animeCategoryRepo := repositories.NewGormAnimeCategoryRepository(dbPrimary, dbReplica)
	songRepo := repositories.NewGormSongRepository(dbPrimary, dbReplica)
	artistRepo := repositories.NewGormArtistRepository(dbPrimary, dbReplica)
	songArtistRepo := repositories.NewGormSongArtistRepository(dbPrimary, dbReplica)
	songChannelRepo := repositories.NewGormSongChannelRepository(dbPrimary, dbReplica)
	studioRepo := repositories.NewGormStudioRepository(dbPrimary, dbReplica)
	animeStudioRepo := repositories.NewGormAnimeStudioRepository(dbPrimary, dbReplica)
	animeCategorryUnivserseRepo := repositories.NewGormAnimeCategoryUniverseRepository(dbPrimary, dbReplica)
	categoryUniverseRepo := repositories.NewGormCategoryUniverseRepository(dbPrimary, dbReplica)
	episodeRepo := repositories.NewGormEpisodeRepository(dbPrimary, dbReplica)
	characterRepo := repositories.NewGormCharacterRepository(dbPrimary, dbReplica)
	animeCharacterRepo := repositories.NewGormAnimeCharacterRepository(dbPrimary, dbReplica)
	episodeCharacterRepo := repositories.NewGormEpisodeCharacterRepository(dbPrimary, dbReplica)

	//create service
	userService := services.NewUserService(userRepo)
	myAnimeService := services.NewMyAnimeService(userAnimeRepo, animeRepo, userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	songService := services.NewSongService(songRepo, animeRepo, artistRepo, songArtistRepo, songChannelRepo)
	artistService := services.NewArtistService(artistRepo)
	studioService := services.NewStudioService(studioRepo)
	animeService := services.NewAnimeService(
		animeRepo,
		userRepo,
		animeCategoryRepo,
		animeStudioRepo,
		songRepo,
		categoryRepo,
		animeCategorryUnivserseRepo,
		categoryUniverseRepo,
		studioRepo,
		episodeRepo,
		s3Service,
	)
	commonService := services.NewCommonService(configService)
	categoryUniverseService := services.NewCategoryUniverseService(categoryUniverseRepo)
	episodeService := services.NewEpisodeService(episodeRepo, animeRepo, episodeCharacterRepo)
	characterService := services.NewCharacterService(characterRepo, animeCharacterRepo, animeRepo)

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
	episodeHandler := adapters.NewHttpEpisodeHandler(episodeService)
	characterHandler := adapters.NewHttpCharacterHandler(characterService)

	rtr := route.InitRoutes(animeHandler,
		userHandler,
		myAnimeHandler,
		categoryHandler,
		songHandler,
		artistHandler,
		studioHandler,
		commonHandler,
		categoryUniverseHandler,
		episodeHandler,
		characterHandler)

	log.Print("Server listening on http://localhost:8080/")
	if err := http.ListenAndServe("0.0.0.0:8080", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
