package handler

import (
	"encoding/json"
	"net/http"

	"github.com/098765432m/internal/service"
	"github.com/098765432m/logger"
)

type WeatherHandler struct {
	Log     *logger.Logger
	Service *service.WeatherService
}

func NewWeatherHandler(log *logger.Logger, service *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		Log:     log,
		Service: service,
	}
}

func (wh *WeatherHandler) GetCityWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "Missing city", http.StatusBadRequest)
		return
	}

	logger.NewLogger().Info.Printf("city: %s\n", city)

	weatherData, err := wh.Service.GetWeather(city)
	if err != nil {
		logger.NewLogger().Error.Println(err)
		http.Error(w, "Error when get city weather", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application.json")
	if err := json.NewEncoder(w).Encode(weatherData); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
