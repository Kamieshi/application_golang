package models

import (
	"fmt"
	"strconv"
)

type Entity struct {
	Id       string `json:"id" db:"id" bson:"_id,omitempty" swaggerignore:"x-nullable,x-abc=def,!x-omitempty"`
	Name     string `json:"name" db:"entity_name" bson:"name"`
	Price    int64  `db:"price" json:"price" bson:"price" validate:"min=1,max=100"`
	IsActive bool   `db:"is_active" json:"is_active" bson:"is_active"`
}

func (ent *Entity) InitForMap(obj map[string]interface{}) error {
	fmt.Println('s')
	name := obj["Name"].(string)
	id := obj["Id"].(string)
	price, _ := strconv.Atoi(obj["Price"].(string))
	isActive := false
	if obj["IsActive"] == "1" {
		isActive = true
	}
	ent.Id = id
	ent.IsActive = isActive
	ent.Price = int64(price)
	ent.Name = name
	return nil
}
