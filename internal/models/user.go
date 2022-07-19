package models

import "github.com/google/uuid"

// User model
type User struct {
	ID           uuid.UUID `json:"id,omitempty" db:"id" bson:"_id" readonly:"true"`
	UserName     string    `json:"username" db:"username" bson:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash" bson:"password_hash"`
	Admin        bool      `json:"admin" db:"is_admin" bson:"admin"`
}
