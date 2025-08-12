package config

import (
	"github.com/spf13/viper"
)

type ConfigService interface {
	GetDatabase() *Database
	GetDatabaseReplica() *Database
	GetCommon() *Common
	GetAWS() *AWS
	GetMyAnimeListClient() *MyAnimeListClient
}

type configService struct {
	service  *viper.Viper
	database *viper.Viper
	client   *viper.Viper
}

type AnimeSeasonalYear struct {
	FirstYear int
}

type Common struct {
	AnimeSeasonalYear AnimeSeasonalYear
}

type AWS struct {
	Region    string
	AccessKey string
	SecretKey string
	S3Bucket  string
}

type Database struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	DatabaseName string
}

type DatabaseReplica struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	DatabaseName string
}

type MyAnimeListClient struct {
	EndPoint string
	CLientID string
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

func NewConfigService() ConfigService {
	initConfig()
	service := viper.Sub("service")
	database := viper.Sub("db")
	client := viper.Sub("client")
	return &configService{
		service:  service,
		database: database,
		client:   client,
	}
}

func (s *configService) GetDatabase() *Database {
	db := s.database.Sub("primary")
	return &Database{
		Host:         db.GetString("host"),
		Port:         db.GetInt("port"),
		UserName:     db.GetString("username"),
		Password:     db.GetString("password"),
		DatabaseName: db.GetString("databaseName"),
	}
}

func (s *configService) GetDatabaseReplica() *Database {
	db := s.database.Sub("replica")
	return &Database{
		Host:         db.GetString("host"),
		Port:         db.GetInt("port"),
		UserName:     db.GetString("username"),
		Password:     db.GetString("password"),
		DatabaseName: db.GetString("databaseName"),
	}
}

func (s *configService) GetAWS() *AWS {
	aws := s.service.Sub("aws")
	return &AWS{
		Region:    aws.GetString("region"),
		AccessKey: aws.GetString("accessKey"),
		SecretKey: aws.GetString("secretKey"),
		S3Bucket:  aws.GetString("s3Bucket"),
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

func (s *configService) GetMyAnimeListClient() *MyAnimeListClient {
	return &MyAnimeListClient{
		EndPoint: s.client.GetString("myAnimeListEndpoint"),
		CLientID: s.client.GetString("myAnimeListClientId"),
	}
}
