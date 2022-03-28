package repository

import (
	"app/internal/service/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// type ControllerEntity interface {
// 	GetAllItems(ctx context.Context) ([]models.Entity, error)
// 	GetItemForID(ctx context.Context, id int) (models.Entity, error)
// 	AddItem(ctx context.Context, obj models.Entity) (bool, error)
// 	UpdateItem(ctx context.Context, id int, obj models.Entity) error
// 	DeleteItem(ctx context.Context, id int) error
// }

type RepoEntityPostgres struct {
	Pool *pgxpool.Pool
}

func rowToEntity(row pgx.Row) (*models.Entity, error) {
	var entety models.Entity
	err := row.Scan(&entety.Id, &entety.Name, &entety.Price, &entety.IsActive)
	if err != nil {
		return &models.Entity{}, err
	}
	return &entety, nil
}

func (sp *RepoEntityPostgres) GetAll(ctx context.Context) ([]models.Entity, error) {
	rows, err := sp.Pool.Query(ctx, "select * from entity")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entitys = make([]models.Entity, 0, len(rows.FieldDescriptions()))

	for rows.Next() {
		entity, err := rowToEntity(rows)
		if err != nil {
			return nil, err
		}

		entitys = append(entitys, *entity)
	}
	return entitys, nil
}

func (sp *RepoEntityPostgres) GetForID(ctx context.Context, id int) (models.Entity, error) {

	var row pgx.Row = sp.Pool.QueryRow(ctx, "select * from entity where id=$1", id)

	ent, err := rowToEntity(row)
	if err != nil {
		return models.Entity{}, err
	}
	return *ent, nil
}

func (sp *RepoEntityPostgres) Add(ctx context.Context, obj models.Entity) error {

	_, err := sp.Pool.Exec(ctx, "insert into entity(entityname,price,isactive) values ($1,$2,$3)", obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (sp *RepoEntityPostgres) Delete(ctx context.Context, id int) error {

	com, err := sp.Pool.Exec(ctx, "delete from entity where id=$1", id)

	if err != nil {
		return err
	}
	if com.String() == "DELETE 0" {
		return errors.New("No find entity for ID")
	}

	return nil
}

func (sp *RepoEntityPostgres) Update(ctx context.Context, id int, obj models.Entity) error {

	com, err := sp.Pool.Exec(ctx, "UPDATE entity SET entityname=$2,price=$3,isactive=$4 WHERE id=$1;", id, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}
	if com.String() == "UPDATE 0" {
		return errors.New("No find entity for ID")
	}

	return nil
}
