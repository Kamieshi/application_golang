package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"app/internal/models"
	"app/internal/repository"
)

// EntityService it's structure for work with cache and entity repository
type EntityService struct {
	rep     repository.RepoEntity
	cashRep repository.CacheEntityRepository
}

// NewEntityService Constructor EntityService struct
func NewEntityService(rep repository.RepoEntity, cahRep repository.CacheEntityRepository) *EntityService {
	return &EntityService{
		rep:     rep,
		cashRep: cahRep,
	}
}

// GetAll returns all entity
func (e *EntityService) GetAll(ctx context.Context) ([]*models.Entity, error) {
	entities, err := e.rep.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service entity/GetAll : %v", err)
	}
	return entities, nil
}

// GetForID Get entity by id
func (e *EntityService) GetForID(ctx context.Context, id string) (*models.Entity, error) {
	entity, exist := e.cashRep.Get(ctx, id)
	if exist {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Info("Use cache")

		return entity, nil
	}
	entity, err := e.rep.GetForID(ctx, id)
	if err != nil {
		return entity, fmt.Errorf("service entity/GetForID : %v", err)
	}
	errCache := e.cashRep.Set(ctx, entity)
	if err != nil {
		logrus.Error(errCache)
	}
	return entity, nil
}

// Add new entity in db
func (e *EntityService) Add(ctx context.Context, obj *models.Entity) error {
	err := e.rep.Add(ctx, obj)
	if err != nil {
		return fmt.Errorf("service entity/Add : %v", err)
	}
	errSet := e.cashRep.Set(ctx, obj)
	if errSet != nil {
		logrus.WithError(errSet).Error()
	}

	return nil
}

// Delete entity from db
func (e *EntityService) Delete(ctx context.Context, id string) error {
	err := e.rep.Delete(ctx, id)
	e.cashRep.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service entity/Delete : %v", err)
	}
	return err
}

// Update entity (post)
func (e *EntityService) Update(ctx context.Context, id string, obj *models.Entity) error {
	err := e.rep.Update(ctx, id, obj)
	e.cashRep.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service entity/Update : %v", err)
	}
	return err
}
