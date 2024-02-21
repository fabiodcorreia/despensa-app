package services

import (
	"github.com/fabiodcorreia/despensa-app/internal/models"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
)

type LocationService struct {
	store *storage.Store
}

func NewLocation(s *storage.Store) *LocationService {
	return &LocationService{
		store: s,
	}
}

func (s *LocationService) GetLocation(locationId string) (models.Location, error) {
	return s.store.GetLocationById(locationId)
}

func (s *LocationService) GetLocationItems(locationId string) ([]models.Item, error) {
	return s.store.GetLocationItems(locationId)
}
