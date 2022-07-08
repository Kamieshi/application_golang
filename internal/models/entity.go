package models

import (
	"github.com/google/uuid"
)

type Entity struct {
	ID       uuid.UUID `json:"id" db:"id" bson:"_id,omitempty" swaggerignore:"x-nullable,x-abc=def,!x-omitempty"`
	Name     string    `json:"name" db:"entity_name" bson:"name"`
	Price    int64     `db:"price" json:"price" bson:"price" validate:"min=1,max=100"`
	IsActive bool      `db:"is_active" json:"is_active" bson:"is_active"`
}
