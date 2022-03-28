package models

type Entity struct {
	Id       int64  `json:"id" db:"id" `
	Name     string `json:"name" db:"entityname"`
	Price    int64  `db:"price" json:"price"`
	IsActive bool   `db:"isactive" json:"is_active"`
}
