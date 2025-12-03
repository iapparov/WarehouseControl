package item

import (
	"errors"
	"github.com/google/uuid"
)

type Item struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Price float64   `json:"price"`
}

func NewItem(name string, count int, price float64) (*Item, error) {
	if name == "" {
		return nil, errors.New("name required")
	}
	if price <= 0 {
		return nil, errors.New("price must be > 0")
	}
	return &Item{
		ID:    uuid.New(),
		Name:  name,
		Count: count,
		Price: price,
	}, nil
}

func (i *Item) ChangeItem(name string, count int, price float64) error {
	if name == "" {
		return errors.New("name cant be empty")
	}
	if price <= 0 {
		return errors.New("price must be >0")
	}
	i.Name = name
	i.Count = count
	i.Price = price
	return nil
}

type FieldDiff struct {
	Old interface{} `json:"old"`
	New interface{} `json:"new"`
}

type ItemDiff map[string]FieldDiff

func (i *Item) Diff(other Item) ItemDiff {
	diff := ItemDiff{}

	if i.Name != other.Name {
		diff["name"] = FieldDiff{Old: i.Name, New: other.Name}
	}
	if i.Count != other.Count {
		diff["count"] = FieldDiff{Old: i.Count, New: other.Count}
	}
	if i.Price != other.Price {
		diff["price"] = FieldDiff{Old: i.Price, New: other.Price}
	}
	return diff
}
