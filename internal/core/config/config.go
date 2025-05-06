package config

import "github.com/spf13/viper"

type ConfigService interface {
	GetDatabase() *Database
	GetCommon() *Common
}

type configService struct {
	service  *viper.Viper
	database *viper.Viper
}

type AnimeSeasonalYear struct {
	FirstYear int
}

type Common struct {
	AnimeSeasonalYear AnimeSeasonalYear
}

type Database struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	DatabaseName string
}

func NewConfigService() ConfigService {
	service := viper.Sub("service")
	database := viper.Sub("database")
	return &configService{
		service:  service,
		database: database,
	}
}

func (s *configService) GetDatabase() *Database {
	return &Database{
		Host:         s.database.GetString("db.host"),
		Port:         s.database.GetInt("db.port"),
		UserName:     s.database.GetString("db.username"),
		Password:     s.database.GetString("db.password"),
		DatabaseName: s.database.GetString("db.databaseName"),
	}
}

func (s *configService) GetCommon() *Common {
	commonService := s.service.Sub("common")
	animeSeasonalYear := commonService.Sub("animeSeasonalYear")
	return &Common{
		AnimeSeasonalYear: AnimeSeasonalYear{
			FirstYear: animeSeasonalYear.GetInt("firstYear"),
		},
	}
}
