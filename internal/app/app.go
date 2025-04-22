package app

import (
	"fmt"
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

	database_env := config.AppData.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		database_env.User,
		database_env.Password,
		database_env.Host,
		database_env.Port,
		database_env.Name,
	)

	// logger.NewLogger().Info.Println(dsn)

	database := &db.Database{}
	conn, err := database.Connect(dsn)
	if err != nil {
		logger.NewLogger().Error.Fatal("Error connecting to database: ", err)
	}
	logger.NewLogger().Info.Println("Connect to database successfully!")
	defer conn.Close()
	//Init database connection END

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
