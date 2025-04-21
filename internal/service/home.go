package service

import "github.com/098765432m/internal/db"

type HomeService struct {
	db *db.Database
}

func NewHomeService(database *db.Database) *HomeService {
	return &HomeService{
		db: database,
	}
}

func (h *HomeService) GetHomeService() string {
	return "There is nothing in home service!"
}