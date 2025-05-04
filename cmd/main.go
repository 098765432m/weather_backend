package main

import (
	"github.com/098765432m/internal/app"
	"github.com/098765432m/logger"
)

func main() {
	if err := app.Run(); err != nil {
		logger.NewLogger().Info.Fatal("Error starting server: ", err)
	}
}