package storage

import (
	"database/sql"
	"fmt"

	"github.com/fabiodcorreia/despensa-app/internal/models"
)

type ItemStore interface {
	GetItemByID(id string) (models.Item, error)
	AddItem(item models.Item) error
	GetAllItems() ([]models.Item, error)
	AddItemStored(item models.Item) error
}

//

const queryGetItemByID = `
  SELECT i.id AS id, i.name AS name, s.quantity AS quantity, s.locationId AS locationId 
  FROM item AS i LEFT JOIN itemStored AS s ON i.id = s.itemId 
  WHERE i.id = ?
`
const queryAddItem = `
  INSERT INTO item (id, name) VALUES (?, ?)
`

const queryAddItemStored = `
  INSERT INTO itemStored (itemId, locationId, quantity) VALUES (?, ?, ?)
`

func (s *Store) GetItemByID(id string) (item models.Item, err error) {
	stmt, err := s.db.Prepare(queryGetItemByID)
	if err != nil {
		return item, fmt.Errorf("store get item by id statement: %w", err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	row := stmt.QueryRowContext(s.ctx, id)
	if err = row.Scan(&item.ID, &item.Name, &item.Quantity, &item.LocationID); err != nil {
		if err == sql.ErrNoRows {
			return item, ErrNotFound
		}
		return item, fmt.Errorf("store get item by id query: %v", err)
	}

	return item, nil
}

const queryGetAllItems = `
  SELECT i.id AS id, i.name AS name, s.quantity AS quantity, s.locationId AS locationId 
  FROM item AS i LEFT JOIN itemStored AS s ON i.id = s.itemId 
`

func (s *Store) GetAllItems() (items []models.Item, err error) {
	stmt, err := s.db.Prepare(queryGetAllItems)
	if err != nil {
		return nil, fmt.Errorf("store get all items statement: %w", err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	rows, err := stmt.QueryContext(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("store get all items query: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		var item models.Item
		if err = rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.LocationID); err != nil {
			return nil, fmt.Errorf("store get all items scan: %w", err)
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("store get all items rows: %w", err)
	}
	return items, nil
}

func (s *Store) AddItem(item models.Item) (err error) {
	stmt, err := s.db.Prepare(queryAddItem)
	if err != nil {
		return fmt.Errorf("store add item statement: %w", err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = stmt.ExecContext(s.ctx, item.ID, item.Name)
	if err != nil {
		return fmt.Errorf("store add item exec: %w", err)
	}
	return nil
}

func (s *Store) AddItemStored(item models.Item) (err error) {
	stmt, err := s.db.Prepare(queryAddItemStored)
	if err != nil {
		return fmt.Errorf("store add item stored statement: %w", err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = stmt.ExecContext(s.ctx, item.ID, item.LocationID, item.Quantity)
	if err != nil {
		return fmt.Errorf("store add item stored exec: %w", err)
	}
	return nil
}
