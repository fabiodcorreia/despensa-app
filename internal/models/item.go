package models

type Item struct {
	Id         string
	Name       string
	Quantity   int8
	LocationId string
}

func NewItem(id, name string) Item {
	return Item{
		Id:   id,
		Name: name,
	}
}

func NewItemWitoutId(name string) Item {
	return NewItem(NewShortID(name), name)
}
