package app

import (
	"net/http"
	"time"

	"github.com/098765432m/config"
	"github.com/098765432m/internal/db"
	"github.com/098765432m/internal/handler"
	redisClient "github.com/098765432m/internal/redis-client"
	"github.com/098765432m/internal/repository"
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

	//Init Redis
	redisClient ,err := redisClient.InitRedisClient("localhost:6379", "", 0)
	if err != nil {
		logger.NewLogger().Error.Fatal("Failed to connect to Redis")
	}

	//Init database connection
	cfg := config.AppData.Database
		db, err := db.NewDatabase(&cfg)
		if err != nil {
			logger.NewLogger().Error.Fatal("Failed to connect to database")
		}
		defer db.CloseDb()

		log := logger.NewLogger()
		apiKey := config.AppData.WeatherApp.ApiKey

	//Init Repository
	userRepo := repository.NewUserRepository(db.GetDb(), redisClient)

	//Init Service
	userService := service.NewUserService(userRepo)
	weatherService := service.NewWeatherService(apiKey)

	//Init handler
	homeHandler := handler.NewHomeHandler(log)
	weatherHandler := handler.NewWeatherHandler(log, weatherService)
	userHandler := handler.NewUserHandler(userService)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.Use(middleware.LogRequest)

	apiRouter.HandleFunc("/", homeHandler.Home).Methods(http.MethodGet)

	apiRouter.HandleFunc("/weather", weatherHandler.GetCityWeather).Methods(http.MethodGet)

	apiRouter.HandleFunc("/user/register", userHandler.Register).Methods(http.MethodPost)
	apiRouter.HandleFunc("/user/login", userHandler.Login).Methods(http.MethodPost)
	apiRouter.HandleFunc("/user/{id}", userHandler.DashboardUpdateUser).Methods(http.MethodPut)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      apiRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
