package services

import (
	"errors"

	"github.com/fabiodcorreia/despensa-app/internal/models"
	"github.com/fabiodcorreia/despensa-app/internal/storage"
)

var ErrItemNotFound = errors.New("item not found")
var ErrItemExists = errors.New("item already exists")

type ItemService struct {
	store *storage.Store
}

func NewItemService(store *storage.Store) ItemService {
	return ItemService{
		store,
	}
}

func (s ItemService) GetItem(itemID string) (models.Item, error) {
	item, err := s.store.GetItemByID(itemID)
	if errors.Is(err, storage.ErrNotFound) {
		return item, ErrItemNotFound
	}
	return item, err
}

func (s ItemService) GetItems() ([]models.Item, error) {
	return s.store.GetAllItems()
}

func (s ItemService) AddItem(item models.Item) error {
	if err := item.Validate(); err != nil {
		return err
	}
	// check if the item exists
	_, err := s.store.GetItemByID(item.ID)
	if err == nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return err
		}
		// if not exists insert
		err = s.store.AddItem(item)
		if err != nil {
			// if insert fails return error
			return err
		}
	}
	// insert storage and return result
	return s.store.AddItemStored(item)
}

//
// func (s ItemService) UpdateItem(item models.Item) error {
// 	_, err := s.store.GetItemById(item.Id)
// 	if err != nil {
// 		if err.is(storage.ErrNotFound) {
// 			return ErrItemNotFound
// 		}
// 		return err
// 	}
//
// 	return s.store.UpdateItem(item)
// }
