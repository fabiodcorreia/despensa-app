package storage

import (
	"database/sql"
	"fmt"

	"github.com/fabiodcorreia/despensa-app/internal/models"
)

func (s *Store) GetItemById(id string) (models.Item, error) {
	var item models.Item

	row := s.db.QueryRow("SELECT id, name FROM item WHERE id = ?", id)
	if err := row.Scan(&item.Id, &item.Name); err != nil {
		if err == sql.ErrNoRows {
			return item, fmt.Errorf("get item by id: item not found for id %s", id)
		}
		return item, fmt.Errorf("get item by id: %v", err)
	}

	return item, nil
}

func (s *Store) GetAllItems() ([]models.Item, error) {
	var items []models.Item
	rows, err := s.db.Query("SELECT id, name FROM item ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("get all items: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Name); err != nil {
			return nil, fmt.Errorf("get all items: %v", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all items: %v", err)
	}

	return items, nil
}

func (s *Store) CreateItem(item models.Item) error {
	result, err := s.db.Exec("INSERT INTO item VALUES (?,?)", item.Id, item.Name)
	if err != nil {
		return fmt.Errorf("create item: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("create item: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("create item: affected %d rows instead of 1", count)
	}

	return nil
}

func (s *Store) UpdateItem(item models.Item) error {
	result, err := s.db.Exec("UPDATE item SET name = ? WHERE id = ?", item.Name, item.Id)
	if err != nil {
		return fmt.Errorf("update item: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update item: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("update item: affected %d rows instead of 1", count)
	}

	return nil
}

func (s *Store) DeleteItem(id string) error {
	result, err := s.db.Exec("DELETE item WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete item: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete item: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("delete item: affected %d rows instead of 1", count)
	}

	return nil
}
