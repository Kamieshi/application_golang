package repository

import (
	"application/service/models"
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

func RowToEntity(row pgx.Row) (*models.Entity, error) {
	var entety models.Entity
	err := row.Scan(&entety.Id, &entety.Name, &entety.Price, &entety.IsActive)
	if err != nil {
		return &models.Entity{}, err
	}
	return &entety, nil
}

func (sp *ShardPostgres) GetAllItems() ([]models.Entity, error) {
	rows, err := sp.conn.Query(context.Background(), "select * from entity")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var entitys = make([]models.Entity, 0, len(rows.FieldDescriptions()))

	for rows.Next() {
		entity, err := RowToEntity(rows)
		if err != nil {
			log.Fatal(err)
		}
	
		entitys = append(entitys, *entity)
	}
	return entitys, nil
}

func (sp *ShardPostgres) GetItemForID(id int) (models.Entity, error) {
	var row pgx.Row = sp.conn.QueryRow(context.Background(), "select * from entity where id=$1", id)

	ent, err := RowToEntity(row)
	if err != nil {
		log.Fatal(err)
	}
	return *ent, nil
}

func (sp *ShardPostgres) AddItem(obj models.Entity) (bool, error) {
	tx, err := sp.conn.Begin(context.Background())
	if err != nil {
		return false, err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "insert into entity(entityname,price,isactive) values ($1,$2,$3)", obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return false, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return false, err
	}
	return true, nil
}

func (sp *ShardPostgres) DeleteItem(id int) error {
	tx, err := sp.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "delete from entity where id=$1", id)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (sp *ShardPostgres) UpdateItem(id int, obj models.Entity) error {
	tx, err := sp.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "UPDATE entity SET entityname=$2,price=$3,isactive=$4 WHERE id=$1;", id, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
