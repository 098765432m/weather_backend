package main

import (
	"fmt"
	"net/http"

	"github.com/098765432m/internal/app"
	"github.com/098765432m/logger"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func main() {
	if err := app.Run(); err != nil {
		logger.NewLogger().Info.Fatal("Error starting server: ", err)
	}
}