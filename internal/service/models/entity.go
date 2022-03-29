package models

type Entity struct {
	Id       interface{} `json:"id" db:"id" bson:"_id"`
	Name     string      `json:"name" db:"entityname" bson:"name"`
	Price    int64       `db:"price" json:"price" bson:"price"`
	IsActive bool        `db:"isactive" json:"is_active" bson:"is_active"`
}

