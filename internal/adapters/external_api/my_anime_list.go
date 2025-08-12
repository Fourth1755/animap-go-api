package external_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/config"
	"github.com/Fourth1755/animap-go-api/internal/logs"
)

type MyAnimeListService interface {
	GetAnimeDetail(myAnimeListId uint) (*GetAnimeDetailResponse, error)
}

type myAnimeListService struct {
	configService config.ConfigService
}

func NewAnimeListService(configService config.ConfigService) MyAnimeListService {
	return &myAnimeListService{configService: configService}
}

type MainPicture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type AlternativeTitles struct {
	Synonyms []string `json:"synonyms"`
	En       string   `json:"en"`
	Ja       string   `json:"ja"`
}

type StartSeason struct {
	Year   int    `json:"year"`
	Season string `json:"season"`
}

type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week"`
	StartTime    string `json:"start_time"`
}

type GetAnimeDetailResponse struct {
	ID                     uint              `json:"id"`
	Title                  string            `json:"title"`
	MainPicture            MainPicture       `json:"main_picture"`
	AlternativeTitles      AlternativeTitles `json:"alternative_titles"`
	StartDate              string            `json:"istart_dated"`
	EndDate                string            `json:"end_date"`
	Synopsis               string            `json:"synopsis"`
	Mean                   float64           `json:"mean"`
	Rank                   uint              `json:"rank"`
	Popularity             uint              `json:"popularity"`
	NumListUsers           uint              `json:"num_list_users"`
	NumScoringUsers        uint              `json:"num_scoring_users"`
	MediaYype              string            `json:"media_type"`
	Status                 string            `json:"finished_airing"`
	NumEpisodes            uint              `json:"num_episodes"`
	StartSeason            StartSeason       `json:"start_season"`
	Broadcast              Broadcast         `json:"broadcast"`
	Source                 string            `json:"source"`
	AverageEpisodeDuration uint              `json:"average_episode_duration"`
	Rating                 string            `json:"rating"`
}

const (
	animeColumn string = "?fields=id,title,main_picture,alternative_titles,start_date,end_date,synopsis,mean,rank,popularity,num_list_users,num_scoring_users,nsfw,created_at,updated_at,media_type,status,genres,my_list_status,num_episodes,start_season,broadcast,source,average_episode_duration,rating,pictures,background,related_anime,related_manga,recommendations,studios,statistics"
)

func (m myAnimeListService) GetAnimeDetail(myAnimeListId uint) (*GetAnimeDetailResponse, error) {
	// 2. สร้าง HTTP client พร้อม timeout
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// 3. สร้าง request และใส่ header
	animeId := strconv.Itoa(int(myAnimeListId))
	resource := fmt.Sprintf(m.configService.GetMyAnimeListClient().EndPoint + "/" + animeId + animeColumn)
	req, err := http.NewRequest(http.MethodGet, resource, nil)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create request"})
		logs.Error(err.Error() + "cannot create request")
		return nil, err
	}

	//req.Header.Set("Authorization", "Bearer your-token-here")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-MAL-CLIENT-ID", m.configService.GetMyAnimeListClient().CLientID)

	// 4. ส่ง request
	resp, err := client.Do(req)
	if err != nil {
		//c.JSON(http.StatusGatewayTimeout, gin.H{"error": fmt.Sprintf("request failed: %v", err)})
		return nil, err
	}
	defer resp.Body.Close()

	// 5. ตรวจสอบสถานะ
	if resp.StatusCode != http.StatusOK {
		//c.JSON(resp.StatusCode, gin.H{"error": "external api returned non-200 status"})
		return nil, err
	}

	// 6. Decode JSON ลง struct
	var data GetAnimeDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot decode json"})
		return nil, err
	}

	return &data, nil
}
