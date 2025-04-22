package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/098765432m/logger"
)

type WeatherApiResponse struct {
	Location struct {
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
		Name string  `json:"name"`
	} `json:"location"`
	Current struct {
		TempC     float32 `json:"temp_c"`
		TempF     float32 `json:"temp_f"`
		IsDay     int8    `json:"is_day"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		WindMph    float32 `json:"wind_mph"`
		WindKph    float32 `json:"wind_kph"`
		Humidity   int     `json:"humidity"`
		WillItRain int8    `json:"will_it_rain"`
		WillItSnow int8    `json:"will_it_snow"`
	} `json:"current"`
}

type WeatherService struct {
	APIKey string
}

func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{APIKey: apiKey}
}

func (w *WeatherService) GetWeather(city string) (*WeatherApiResponse, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", w.APIKey, encodedCity)
	logger.NewLogger().Info.Println(url)
	res, err := http.Get(url)
	if err != nil {
		logger.NewLogger().Error.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error when get data, status : %d", res.StatusCode)
	}

	var data WeatherApiResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		logger.NewLogger().Error.Println("Error decoding response:", err)
		return nil, err
	}

	return &data, nil
}
