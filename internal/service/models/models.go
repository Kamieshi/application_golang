package models

type Entity struct {
	Id       int64  `db: "id"`
	Name     string `db: "entityname"`
	Price    int64  `db: "price"`
	IsActive bool   `db: "isactive"`
}
