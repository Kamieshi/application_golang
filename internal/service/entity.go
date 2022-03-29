package service

import (
	"app/internal/repository"
	"app/internal/service/models"
	"context"
)

type EntityService struct {
	rep repository.RepoEntity
}

func NewEntityService(rep repository.RepoEntity) EntityService {
	return EntityService{
		rep: rep,
	}
}

func (es EntityService) GetAll(ctx context.Context) ([]models.Entity, error) {

	entitys, err := es.rep.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return entitys, nil
}

func (es *EntityService) GetForID(ctx context.Context, id interface{}) (models.Entity, error) {
	entity, err := es.rep.GetForID(ctx, id)
	if err != nil {
		return models.Entity{}, err
	}
	return entity, nil
}

func (es EntityService) Add(ctx context.Context, obj models.Entity) error {
	err := es.rep.Add(ctx, obj)
	return err
}

func (es EntityService) Delete(ctx context.Context, id interface{}) error {
	err := es.rep.Delete(ctx, id)
	return err
}

func (es EntityService) Update(ctx context.Context, id interface{}, obj models.Entity) error {
	err := es.rep.Update(ctx, id, obj)
	return err
}
