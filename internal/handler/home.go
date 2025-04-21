package handler

import (
	"net/http"

	"github.com/098765432m/internal/service"
	"github.com/098765432m/logger"
)

type HomeHandler struct {
	Log *logger.Logger
	Service *service.HomeService
}

func NewHomeHandler(Log logger.Logger) *HomeHandler {
	return &HomeHandler {
		Log: &Log,
	}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	display := h.Service.GetHomeService()
	if display == "" {
		h.Log.Error.Println("Error: Home service is empty!")
		display = "There is no data in service!"
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(display))
}