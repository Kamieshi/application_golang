// Package models include all models for application
package models

import (
	"github.com/google/uuid"
)

// Entity model
type Entity struct {
	ID       uuid.UUID `json:"id" db:"id" bson:"_id,omitempty" readonly:"true"`
	Name     string    `json:"name" db:"name" bson:"name"`
	Price    int32     `db:"price" json:"price" bson:"price" validate:"min=1,max=1000"`
	IsActive bool      `db:"is_active" json:"isActive" bson:"is_active" protobuf:"varint,4,opt,name=isActive,proto3"`
}
