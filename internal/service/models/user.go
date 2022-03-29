package models

type User struct{
	Id       interface{} `json:"id" db:"id" bson:"_id"`
	FirstName     string      `json:"name" db:"entityname" bson:"name"`
	LastName    int64       `db:"price" json:"price" bson:"price"`
	banned bool        `db:"banned" json:"is_active" bson:"is_active"`
}