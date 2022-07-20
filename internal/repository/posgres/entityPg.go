// Package repository repositories from work with models
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"app/internal/models"
)

// RepoEntityPostgres implementation RepositoryEntity like Postgres
type RepoEntityPostgres struct {
	pool *pgxpool.Pool
}

// NewRepoEntityPostgres Constructor
func NewRepoEntityPostgres(pool *pgxpool.Pool) *RepoEntityPostgres {
	return &RepoEntityPostgres{
		pool: pool,
	}
}

// GetAll return all entities from table entity
func (sp RepoEntityPostgres) GetAll(ctx context.Context) ([]*models.Entity, error) {
	query := "SELECT id,entity_name,price,is_active FROM entity"
	rows, err := sp.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities = make([]*models.Entity, 0, len(rows.FieldDescriptions()))
	var entity models.Entity
	for rows.Next() {
		err = rows.Scan(&entity.ID, &entity.Name, &entity.Price, &entity.IsActive)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}
	return entities, nil
}

// GetForID Get entity by id
func (sp RepoEntityPostgres) GetForID(ctx context.Context, id string) (*models.Entity, error) {
	query := "SELECT id,entity_name,price,is_active FROM entity WHERE id=$1"
	marshalUUID, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return nil, fmt.Errorf("entityPg.go/Update : %v", err)
	}
	var row = sp.pool.QueryRow(ctx, query, marshalUUID)
	var entity models.Entity
	err = row.Scan(&entity.ID, &entity.Name, &entity.Price, &entity.IsActive)
	if err != nil {
		return nil, fmt.Errorf("entityPg.go/Update : %v", err)
	}
	return &entity, fmt.Errorf("entityPg.go/Update : %v", err)
}

// Add Write new entity into DB
func (sp RepoEntityPostgres) Add(ctx context.Context, obj *models.Entity) error {
	idRow := uuid.New()
	query := "INSERT INTO entity(id,entity_name,price,is_active) values ($1,$2,$3,$4)"
	_, err := sp.pool.Exec(ctx, query, idRow, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return fmt.Errorf("entityPg.go/Add : %v", err)
	}
	obj.ID = idRow
	return nil
}

// Delete Delete entity form table entities
func (sp RepoEntityPostgres) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM entity WHERE id=$1"
	com, err := sp.pool.Exec(ctx, query, id)
	if com.Delete() || err == pgx.ErrNoRows {
		return nil
	}
	return fmt.Errorf("entityPg.go/Delete : %v", err)
}

// Update entity
func (sp RepoEntityPostgres) Update(ctx context.Context, id string, obj *models.Entity) error {
	marshalUUID, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return fmt.Errorf("entityPg.go/Update : %v", err)
	}
	query := "UPDATE entity SET entity_name=$2,price=$3,is_active=$4 WHERE id=$1"
	com, err := sp.pool.Exec(ctx, query, marshalUUID, obj.Name, obj.Price, obj.IsActive)
	if err != nil {
		return fmt.Errorf("entityPg.go/Update : %v", err)
	}
	if !com.Update() {
		return fmt.Errorf("entityPg.go/Update : %v", errors.New("no find entity for ID"))
	}
	return nil
}
