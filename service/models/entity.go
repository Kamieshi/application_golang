package models

import "fmt"

type Entity struct {
	Id       int64  `db: "id"`
	Name     string `db: "entityname"`
	Price    int64  `db: "price"`
	IsActive bool   `db: "isactive"`
}

func CreateEntity(name string, price int64, isActive bool) *Entity {
	return &Entity{
		Name:     name,
		Price:    price,
		IsActive: isActive,
	}
}

func (e Entity) Show() string {
	return fmt.Sprintf("%v %v %v %v", e.Id, e.Name, e.Price, e.IsActive)
}
