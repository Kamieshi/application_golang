package repository

import (
	"app/internal/models"
	"context"
	"errors"
	"fmt"
	"strconv"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RepoEntityPostgres struct {
	pool *pgxpool.Pool
}

func NewRepoEntityPostgres(pool pgxpool.Pool) RepoEntityPostgres {
	return RepoEntityPostgres{
		pool: &pool,
	}
}

func rowToEntity(row pgx.Row) (*models.Entity, error) {
	var entity models.Entity
	err := row.Scan(&entity.Id, &entity.Name, &entity.Price, &entity.IsActive)
	if err != nil {
		return &models.Entity{}, err
	}
	return &entity, nil
}

func (sp *RepoEntityPostgres) GetAll(ctx context.Context) ([]models.Entity, error) {
	rows, err := sp.pool.Query(ctx, "select * from entity")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities = make([]models.Entity, 0, len(rows.FieldDescriptions()))

	for rows.Next() {
		entity, err := rowToEntity(rows)
		if err != nil {
			return nil, err
		}

		entities = append(entities, *entity)
	}
	return entities, nil
}

func (sp *RepoEntityPostgres) GetForID(ctx context.Context, id string) (models.Entity, error) {

	var row pgx.Row = sp.pool.QueryRow(ctx, "select * from entity where id=$1", id)
	ent, err := rowToEntity(row)
	return *ent, err

}

func (sp *RepoEntityPostgres) Add(ctx context.Context, obj models.Entity) error {

	_, err := sp.pool.Exec(ctx, "insert into entity(entity_name,price,is_active) values ($1,$2,$3)", obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (sp *RepoEntityPostgres) Delete(ctx context.Context, id string) error {
	Id, err := strconv.Atoi(fmt.Sprint(id))
	if err != nil {
		return err
	}
	com, err := sp.pool.Exec(ctx, "delete from entity where id=$1", Id)

	if err != nil {
		return err
	}
	if com.String() == "DELETE 0" {
		return errors.New("No find entity for ID")
	}

	return nil
}

func (sp *RepoEntityPostgres) Update(ctx context.Context, id string, obj models.Entity) error {
	Id, err := strconv.Atoi(fmt.Sprint(id))
	if err != nil {
		return err
	}
	com, err := sp.pool.Exec(ctx, "UPDATE entity SET entity_name=$2,price=$3,is_active=$4 WHERE id=$1;", Id, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}
	if com.String() == "UPDATE 0" {
		return errors.New("No find entity for ID")
	}

	return nil
}
