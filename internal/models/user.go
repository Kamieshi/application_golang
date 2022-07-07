package models

type User struct {
	Id           interface{} `json:"id,omitempty" db:"id" bson:"_id"`
	UserName     string      `json:"username" db:"username" bson:"username"`
	PasswordHash string      `json:"password_hash" db:"password_hash" bson:"password_hash"`
	Admin        bool        `json:"admin" db:"is_admin" bson:"admin"`
}
