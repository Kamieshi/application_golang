package repository

import (
	"app/internal/service/models"
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ControllerEntity interface {
	GetAllItems(ctx context.Context) ([]models.Entity, error)
	GetItemForID(ctx context.Context, id int) (models.Entity, error)
	AddItem(ctx context.Context, obj models.Entity) (bool, error)
	UpdateItem(ctx context.Context, id int, obj models.Entity) error
	DeleteItem(ctx context.Context, id int) error
}

type RepoPostgres struct {
	pool *pgxpool.Pool
}

func NewRepoPostgres(urlConnection string) (*RepoPostgres, error) {
	connPool, err := pgxpool.Connect(context.Background(), urlConnection)
	if err != nil {
		log.Println("Connecting url", urlConnection)
		return nil, err
	}
	return &RepoPostgres{
		pool: connPool,
	}, nil
}

func RowToEntity(row pgx.Row) (*models.Entity, error) {
	var entety models.Entity
	err := row.Scan(&entety.Id, &entety.Name, &entety.Price, &entety.IsActive)
	if err != nil {
		return &models.Entity{}, err
	}
	return &entety, nil
}

func (sp *RepoPostgres) GetAllItems(ctx context.Context) ([]models.Entity, error) {
	conn, err := sp.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, "select * from entity")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entitys = make([]models.Entity, 0, len(rows.FieldDescriptions()))

	for rows.Next() {
		entity, err := RowToEntity(rows)
		if err != nil {
			return nil, err
		}

		entitys = append(entitys, *entity)
	}
	return entitys, nil
}

func (sp *RepoPostgres) GetItemForID(ctx context.Context, id int) (models.Entity, error) {
	conn, err := sp.pool.Acquire(ctx)
	if err != nil {
		return models.Entity{}, err
	}

	var row pgx.Row = conn.QueryRow(ctx, "select * from entity where id=$1", id)

	ent, err := RowToEntity(row)
	if err != nil {
		return models.Entity{}, err
	}
	return *ent, nil
}

func (sp *RepoPostgres) AddItem(ctx context.Context, obj models.Entity) (bool, error) {
	conn, err := sp.pool.Acquire(ctx)
	if err != nil {
		return false, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return false, err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "insert into entity(entityname,price,isactive) values ($1,$2,$3)", obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (sp *RepoPostgres) DeleteItem(ctx context.Context, id int) error {
	conn, err := sp.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "delete from entity where id=$1", id)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sp *RepoPostgres) UpdateItem(ctx context.Context, id int, obj models.Entity) error {

	conn, err := sp.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE entity SET entityname=$2,price=$3,isactive=$4 WHERE id=$1;", id, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
