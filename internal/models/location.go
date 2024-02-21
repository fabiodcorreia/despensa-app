package models

type Location struct {
	Id   string
	Name string
}

func NewLocation(name string) Location {
	return Location{
		Id:   NewShortID(name),
		Name: name,
	}
}
