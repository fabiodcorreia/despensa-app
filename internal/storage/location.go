package storage

import (
	"database/sql"
	"fmt"

	"github.com/fabiodcorreia/despensa-app/internal/models"
)

type LocationStore interface {
	AddLocation(loc models.Location) error
	GetLocationByID(id string) (models.Location, error)
	GetAllLocations() ([]models.Location, error)
}

const queryAddLocation = `
  INSERT INTO location (id, name) VALUES (?, ?)
`
const queryGetLocationByID = `
  SELECT id, name FROM location WHERE id = ? LIMIT 1
`
const queryGetAllLocations = `
  SELECT id, name FROM location
`

// GetLocationByID will search for a location with the given id.
//
// Returns storage.ErrNotFound if location not found.
//
// Returns an error if the query fails for other reasons.
func (s *Store) GetLocationByID(id string) (models.Location, error) {
	var location models.Location
	stmt, err := s.db.Prepare(queryGetLocationByID)
	if err != nil {
		return location, fmt.Errorf("store get location by id statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(s.ctx, id)
	if err = row.Scan(&location.ID, &location.Name); err != nil {
		if err == sql.ErrNoRows {
			return location, ErrNotFound
		}
		return location, fmt.Errorf("store get location by id query: %v", err)
	}

	return location, nil
}

// GetAllLocations returns all locations from the database
//
// Returns an error if the query fails for other reasons.
func (s *Store) GetAllLocations() ([]models.Location, error) {
	stmt, err := s.db.Prepare(queryGetAllLocations)
	if err != nil {
		return nil, fmt.Errorf("store get all locations statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("store get all locations query: %w", err)
	}
	defer rows.Close()

	locations := make([]models.Location, 0)
	for rows.Next() {
		var location models.Location
		if err = rows.Scan(&location.ID, &location.Name); err != nil {
			return nil, fmt.Errorf("store get all locations scan: %w", err)
		}
		locations = append(locations, location)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("store get all locations rows: %w", err)
	}
	return locations, nil
}

// AddLocation adds a new location to the database
//
// Returns an error if the query fails for other reasons.
func (s *Store) AddLocation(loc models.Location) error {
	stmt, err := s.db.Prepare(queryAddLocation)
	if err != nil {
		return fmt.Errorf("store add location statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(s.ctx, loc.ID, loc.Name)
	if err != nil {
		return fmt.Errorf("store add location exec: %w", err)
	}
	return nil
}
