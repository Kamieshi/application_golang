package models

type User struct {
	Id           interface{} `json:"id" db:"id" bson:"_id"`
	UserName     string      `json:"username" db:"username" bson:"username"`
	PasswordHash string			 `json:"pasword_hash" db:"pasword_hash" bson:"pasword_hash"`
}