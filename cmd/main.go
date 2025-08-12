package main

import (
	"log"
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/adapters/aws"
	"github.com/Fourth1755/animap-go-api/internal/adapters/external_api"
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
	awsAdapter, err := aws.NewAwsAdapter(configService.GetAWS())
	if err != nil {
		log.Fatalf("Could not create aws adapter: %v", err)
	}
	s3Service := aws.NewS3Service(awsAdapter)

	// init external api
	myAnimeListService := external_api.NewAnimeListService(configService)

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
	userService := services.NewUserService(userRepo, s3Service, configService)
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
		myAnimeListService,
	)
	commonService := services.NewCommonService(configService)
	categoryUniverseService := services.NewCategoryUniverseService(categoryUniverseRepo)
	episodeService := services.NewEpisodeService(episodeRepo, animeRepo, episodeCharacterRepo)
	characterService := services.NewCharacterService(characterRepo, animeCharacterRepo, animeRepo)

	//create handler

	rtr := route.InitRoutes(route.HttpHandler{
		AnimeHandler:            adapters.NewHttpAnimeHandler(animeService),
		UserHandler:             adapters.NewHttpUserHandler(userService),
		MyAnimeHandler:          adapters.NewHttpMyAnimeHandler(myAnimeService),
		CategoryHandler:         adapters.NewHttpCategoryHandler(categoryService),
		SongHandler:             adapters.NewHttpSongHandler(songService),
		ArtistHandler:           adapters.NewHttpArtistHandler(artistService),
		StudioHandler:           adapters.NewHttpStduioHandler(studioService),
		CommonHandler:           adapters.NewHttpCommonHandler(commonService),
		CategoryUniverseHandler: adapters.NewHttpCategoryUniverseHandler(categoryUniverseService),
		EpisodeHandler:          adapters.NewHttpEpisodeHandler(episodeService),
		CharacterHandler:        adapters.NewHttpCharacterHandler(characterService),
	})

	log.Print("Server listening on http://localhost:8080/")
	if err := http.ListenAndServe("0.0.0.0:8080", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
