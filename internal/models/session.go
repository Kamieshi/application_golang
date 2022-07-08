package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID              uuid.UUID `json:"id" db:"id" bson:"_id"`
	UserId          uuid.UUID `json:"user_id" db:"user_id" bson:"user_id"`
	RfToken         string    `json:"refresh_token" db:"refresh_token" bson:"refresh_token"`
	UniqueSignature string    `json:"signature" db:"signature" bson:"signature"`
	CreatedAt       time.Time `json:"created_at" db:"created_at" bson:"created_at"`
	Disabled        bool      `json:"disabled" db:"disabled" bson:"disabled"`
}
