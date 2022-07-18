package repository

import (
	"app/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RepoEntityPostgres struct {
	pool *pgxpool.Pool
}

func NewRepoEntityPostgres(pool *pgxpool.Pool) *RepoEntityPostgres {
	return &RepoEntityPostgres{
		pool: pool,
	}
}

const orderColumnsEntity string = "id,entity_name,price,is_active"

func rowToEntity(row pgx.Row) (*models.Entity, error) {
	var entity models.Entity
	err := row.Scan(&entity.ID, &entity.Name, &entity.Price, &entity.IsActive)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (sp RepoEntityPostgres) GetAll(ctx context.Context) ([]*models.Entity, error) {
	query := fmt.Sprintf("SELECT %s FROM entity", orderColumnsEntity)
	rows, err := sp.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities = make([]*models.Entity, 0, len(rows.FieldDescriptions()))

	for rows.Next() {
		entity, err := rowToEntity(rows)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}
	return entities, nil
}

func (sp RepoEntityPostgres) GetForID(ctx context.Context, id string) (*models.Entity, error) {
	query := fmt.Sprintf("SELECT %s FROM entity WHERE id=$1", orderColumnsEntity)
	marshalUUID, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return nil, err
	}
	var row = sp.pool.QueryRow(ctx, query, marshalUUID)
	ent, err := rowToEntity(row)
	return ent, err

}

func (sp RepoEntityPostgres) Add(ctx context.Context, obj *models.Entity) error {
	idRow := uuid.New()
	query := "INSERT INTO entity(id,entity_name,price,is_active) values ($1,$2,$3,$4)"
	_, err := sp.pool.Exec(ctx, query, idRow, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}
	obj.ID = idRow
	return nil
}

func (sp RepoEntityPostgres) Delete(ctx context.Context, id string) error {
	marshalUUID, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	query := "DELETE FROM entity WHERE id=$1"
	_, err = sp.pool.Exec(ctx, query, marshalUUID)

	return fmt.Errorf("Error Delete Postgres rep : %v", err)
}

func (sp RepoEntityPostgres) Update(ctx context.Context, id string, obj *models.Entity) error {
	marshalUUID, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return err
	}
	query := "UPDATE entity SET entity_name=$2,price=$3,is_active=$4 WHERE id=$1"
	com, err := sp.pool.Exec(ctx, query, marshalUUID, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}
	if com.String() == "UPDATE 0" {
		return errors.New("no find entity for ID")
	}
	return nil
}
