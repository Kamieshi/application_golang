package models

type Entity struct {
	Id       interface{} `json:"id" db:"id" bson:"_id,omitempty"`
	Name     string      `json:"name" db:"entity_name" bson:"name"`
	Price    int64       `db:"price" json:"price" bson:"price"`
	IsActive bool        `db:"is_active" json:"is_active" bson:"is_active"`
}
