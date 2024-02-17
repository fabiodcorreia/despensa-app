package storage

import (
	"database/sql"
	"fmt"

	"github.com/fabiodcorreia/despensa-app/internal/model"
)

//TODO: Replace the Exec and Query with Prepared Statment

func (s *Store) GetLocationById(id string) (model.Location, error) {
	var location model.Location

	row := s.db.QueryRow("SELECT id, name FROM location WHERE id = ?", id)
	if err := row.Scan(&location.Id, &location.Name); err != nil {
		if err == sql.ErrNoRows {
			return location, fmt.Errorf("get location by id: location not found for id %s", id)
		}
		return location, fmt.Errorf("get location by id: %v", err)
	}

	return location, nil
}

func (s *Store) GetAllLocations() ([]model.Location, error) {
	var locations []model.Location

	rows, err := s.db.Query("SELECT id, name FROM location ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("get all locations: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var loc model.Location
		if err := rows.Scan(&loc.Id, &loc.Name); err != nil {
			return nil, fmt.Errorf("get all locations: %v", err)
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all locations: %v", err)
	}

	return locations, nil
}

func (s *Store) CreateLocation(loc model.Location) error {
	result, err := s.db.Exec("INSERT INTO location VALUES (?,?)", loc.Id, loc.Name)
	if err != nil {
		return fmt.Errorf("create location: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("create location: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("create location: affected %d rows instead of 1", count)
	}

	return nil
}

func (s *Store) UpdateLocation(loc model.Location) error {
	result, err := s.db.Exec("UPDATE location SET name = ? WHERE id = ?", loc.Name, loc.Id)
	if err != nil {
		return fmt.Errorf("update location: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update location: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("update location: affected %d rows instead of 1", count)
	}

	return nil
}

func (s *Store) DeleteLocation(id string) error {
	result, err := s.db.Exec("DELETE location WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete location: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete location: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("delete location: affected %d rows instead of 1", count)
	}

	return nil
}