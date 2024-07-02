package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	API_KEY        = "<YOU_API_KEY>"
	WeatherAPI_URL = "https://api.weatherapi.com/v1/current.json"
)

type WeatherService interface {
	GetWeatherByCity(city string) (*WeatherAPIResponse, error)
}

// WeatherAPIService is a service to interact with the WeatherAPI API
type WeatherAPIService struct {
	apiKey string
	BaseHttpService
}

type WeatherAPIResponseCurrent struct {
	LastUpdatedEpoch int     `json:"last_updated_epoch"`
	LastUpdated      string  `json:"last_updated"`
	TempC            float64 `json:"temp_c"`
	TempF            float64 `json:"temp_f"`
	IsDay            int     `json:"is_day"`
	Condition        struct {
		Text string `json:"text"`
		Icon string `json:"icon"`
		Code int    `json:"code"`
	} `json:"condition"`
	WindMph    float64 `json:"wind_mph"`
	WindKph    float64 `json:"wind_kph"`
	WindDegree int     `json:"wind_degree"`
	WindDir    string  `json:"wind_dir"`
	PressureMb float64 `json:"pressure_mb"`
	PressureIn float64 `json:"pressure_in"`
	PrecipMm   float64 `json:"precip_mm"`
	PrecipIn   float64 `json:"precip_in"`
	Humidity   int     `json:"humidity"`
	Cloud      int     `json:"cloud"`
	FeelslikeC float64 `json:"feelslike_c"`
	FeelslikeF float64 `json:"feelslike_f"`
	VisKm      float64 `json:"vis_km"`
	VisMiles   float64 `json:"vis_miles"`
	Uv         float64 `json:"uv"`
	GustMph    float64 `json:"gust_mph"`
	GustKph    float64 `json:"gust_kph"`
}

type WeatherAPIResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current WeatherAPIResponseCurrent `json:"current"`
}

// NewWeatherAPIService creates a new WeatherAPIService
func NewWeatherAPIService(apiKey string) WeatherService {
	return &WeatherAPIService{
		apiKey,
		BaseHttpService{Client: &http.Client{}},
	}
}

func (w *WeatherAPIService) GetWeatherByCity(city string) (*WeatherAPIResponse, error) {
	base, _ := url.Parse(WeatherAPI_URL)
	params := url.Values{}
	params.Add("key", w.apiKey)
	params.Add("q", city)
	base.RawQuery = params.Encode()
	resp, err := w.Client.Get(base.String())
	if err != nil {
		log.Println("error getting weather: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("error getting weather: statusCode:%d Response:%s\n", resp.StatusCode, resp.Body)
		return nil, fmt.Errorf("error getting weather: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body: ", err)
		return nil, err
	}

	var weatherResponse WeatherAPIResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Println("error on Unmarshal response body: ", err)
		return nil, err
	}

	return &weatherResponse, nil
}
