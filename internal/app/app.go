package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/098765432m/config"
	"github.com/098765432m/internal/db"
	"github.com/098765432m/internal/handler"
	"github.com/098765432m/logger"
	"github.com/gorilla/mux"
)

func Run() error {
	config.InitConfig()
	port := config.AppData.Database.Port

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

	database := &db.Database{}
	conn, err := database.Connect(dsn)
	if  err != nil  {
		logger.NewLogger().Error.Fatal("Error connecting to database: ", err)
	}
	defer conn.Close()
	//Init database connection END



	homeHandler := handler.NewHomeHandler(*logger.NewLogger())

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler.Home).Methods(http.MethodGet)

	server := &http.Server {
		Addr: ":" + port,
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}