package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"github.com/sirupsen/logrus"
)

// EntityService it's structure for work with cache and entity repository
type EntityService struct {
	rep     repository.RepoEntity
	cashRep repository.CacheEntityRepository
}

// NewEntityService return
func NewEntityService(rep *repository.RepoEntity, cahRep repository.CacheEntityRepository) *EntityService {

	return &EntityService{
		rep:     *rep,
		cashRep: cahRep,
	}
}

func (e EntityService) GetAll(ctx context.Context) ([]*models.Entity, error) {

	entities, err := e.rep.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (e *EntityService) GetForID(ctx context.Context, id string) (*models.Entity, error) {

	entity, is_exist := e.cashRep.Get(ctx, id)
	if is_exist {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Info("Use cache")

		return entity, nil
	}
	entity, err := e.rep.GetForID(ctx, id)
	if err != nil {
		return entity, err
	}
	errCache := e.cashRep.Set(ctx, entity)
	if err != nil {
		logrus.Error(errCache)
	}
	return entity, err
}

func (e EntityService) Add(ctx context.Context, obj *models.Entity) error {
	err := e.rep.Add(ctx, obj)
	if err != nil {
		return err
	}

	errSet := e.cashRep.Set(ctx, obj)
	if errSet != nil {
		logrus.WithError(errSet)
	}

	logrus.Info("Value successful get in cash repository ")

	return err
}

func (e EntityService) Delete(ctx context.Context, id string) error {
	err := e.rep.Delete(ctx, id)
	e.cashRep.Delete(ctx, id)
	return err
}

func (e EntityService) Update(ctx context.Context, id string, obj *models.Entity) error {
	err := e.rep.Update(ctx, id, obj)
	if err != nil {
		return err
	}
	e.cashRep.Delete(ctx, id)
	return err
}
