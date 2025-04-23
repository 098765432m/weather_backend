package app

import (
	"net/http"
	"time"

	"github.com/098765432m/config"
	"github.com/098765432m/internal/db"
	"github.com/098765432m/internal/handler"
	"github.com/098765432m/internal/service"
	"github.com/098765432m/logger"
	"github.com/098765432m/middleware"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func Run() error {
	config.InitConfig()
	port := viper.GetString("app.port")

	if port == "" {
		port = "8500" //Default port
	}

	logger.NewLogger().Info.Printf("Server is running on port: %s\n", port)

	//Init database connection
	cfg := config.AppData.Database
	db, err := db.NewDatabase(&cfg)
	if err != nil {
		logger.NewLogger().Error.Fatal("Failed to connect to database")
	}
	defer db.CloseDb()

	log := logger.NewLogger()
	apiKey := config.AppData.WeatherApp.ApiKey

	//Init Service
	weatherService := service.NewWeatherService(apiKey)

	//Init handler
	homeHandler := handler.NewHomeHandler(log)
	weatherHandler := handler.NewWeatherHandler(log, weatherService)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.Use(middleware.LogRequest)

	apiRouter.HandleFunc("/", homeHandler.Home).Methods(http.MethodGet)

	apiRouter.HandleFunc("/weather", weatherHandler.GetCityWeather).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      apiRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
