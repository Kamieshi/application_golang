package models

import "time"

type Session struct {
	Id             interface{} `json:"id" db:"id" bson:"_id"`
	UserId         interface{} `json:"user_id" db:"user_id" bson:"user_id`
	idSesison      string      `json:"session_id" db:"session_id" bson:"session_id"`
	rfToken        string      `json:refrash_token" db:"refrash_token" bson:"refrash_token"`
	uniqueSignatur string      `json:"signature" db:"signature" bson:"signature"`
	createdAt      time.Time   `json:"created_at" db:"created_at" bson:"created_at"`
}
