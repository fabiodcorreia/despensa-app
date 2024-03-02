package models

import "fmt"

type Item struct {
	ID         string
	Name       string
	Quantity   int8
	LocationID string
}

func NewItemWithID(id, name string) Item {
	return Item{
		ID:   id,
		Name: name,
	}
}

func NewItem(name string) Item {
	return NewItemWithID(NewShortID(name), name)
}

// TODO: move this errors to const

func (i Item) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("item id cannot be empty")
	}
	if i.Name == "" {
		return fmt.Errorf("item name cannot be empty")
	}
	if i.LocationID == "" {
		return fmt.Errorf("item location cannot be empty")
	}
	if i.Quantity < 0 {
		return fmt.Errorf("item quantity cannot be negative")
	}
	return nil
}
