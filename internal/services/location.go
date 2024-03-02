package services

import (
	"errors"

	"github.com/fabiodcorreia/despensa-app/internal/models"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
)

// ErrLocationNotFound is returned when a location is not found
var ErrLocationNotFound = errors.New("location not found")

// ErrLocationExists is returned when a location already exists
var ErrLocationExists = errors.New("location already exists")

type LocationService struct {
	store storage.LocationStore
}

func NewLocationService(store storage.LocationStore) LocationService {
	return LocationService{
		store,
	}
}

// GetLocation returns a location from the store
//
// # Returns ErrLocationNotFound if the location is not found
//
// Returns an error if the store fails to get the location
func (s LocationService) GetLocation(locationID string) (models.Location, error) {
	if locationID == "" {
		return models.Location{}, ErrLocationNotFound
	}

	loc, err := s.store.GetLocationByID(locationID)
	if errors.Is(err, storage.ErrNotFound) {
		return loc, ErrLocationNotFound
	}
	return loc, err
}

func (s LocationService) GetLocations() ([]models.Location, error) {
	return s.store.GetAllLocations()
}

// AddLocation adds a new location to the store
//
// # Returns ErrLocationExists if the location already exists
//
// Returns an error if the store fails to add the location
func (s LocationService) AddLocation(loc models.Location) error {
	if err := loc.Validate(); err != nil {
		return err
	}

	_, err := s.store.GetLocationByID(loc.ID)
	if err == nil {
		return ErrLocationExists
	}

	return s.store.AddLocation(loc)
}
