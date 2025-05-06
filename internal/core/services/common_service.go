package services

import (
	"strconv"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/config"
	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
)

type CommonService interface {
	GetSeasonalAndYear() (*dtos.GetSeasonalAndYearResponse, error)
}

type commonService struct {
	configService config.ConfigService
}

func NewCommonService(configService config.ConfigService) CommonService {
	return &commonService{
		configService: configService,
	}
}

const (
	WINTER = "winter"
	SPRING = "spring"
	SUMMER = "summer"
	FALL   = "fall"
)

func (s commonService) GetSeasonalAndYear() (*dtos.GetSeasonalAndYearResponse, error) {
	yearNow := time.Now().Year()
	firstYear := s.configService.GetCommon().AnimeSeasonalYear.FirstYear
	listOfSeasonal := []string{WINTER, SPRING, SUMMER, FALL}
	var responseData []dtos.GetSeasonalAndYearResponseData
	// เพิ่ม 1 ปีจากปัจจุบัน
	for year := firstYear; year <= yearNow+1; year++ {
		for _, season := range listOfSeasonal {
			yearStr := strconv.Itoa(year)
			// if params.Seasonal == season && params.Year == yearStr {
			// 	isMain = true
			// }

			responseData = append(responseData, dtos.GetSeasonalAndYearResponseData{
				Seasonal: season,
				Year:     yearStr,
			})
		}
	}
	return &dtos.GetSeasonalAndYearResponse{
		Data: responseData,
	}, nil
}
