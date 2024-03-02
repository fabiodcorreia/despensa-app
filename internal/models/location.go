package models

import (
	"errors"
)

type Location struct {
	ID   string
	Name string
}

var ErrLocationIDTooLong = errors.New("location id cannot be longer than 10 characters")
var ErrLocationIDEmpty = errors.New("location id cannot be empty")
var ErrLocationNameEmpty = errors.New("location name cannot be empty")
var ErrLocationNameTooLong = errors.New("location name cannot be longer than 24 characters")

func NewLocationWithID(id, name string) Location {
	return Location{
		ID:   id,
		Name: name,
	}
}

func NewLocation(name string) Location {
	return NewLocationWithID(NewShortID(name), name)
}

func (l Location) Validate() error {
	if l.ID == "" {
		return ErrLocationIDEmpty
	}
	if l.Name == "" {
		return ErrLocationNameEmpty
	}
	if len(l.ID) > 10 {
		return ErrLocationIDTooLong
	}
	if len(l.Name) > 24 {
		return ErrLocationNameTooLong
	}
	return nil
}
