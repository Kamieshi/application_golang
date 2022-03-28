package service

import (
	"app/internal/repository"
	"app/internal/service/models"
	"context"
)

type EntityService struct {
	Rep repository.RepoEntityPostgres
}

func (es EntityService) GetAll(ctx context.Context) ([]models.Entity, error) {

	entitys, err := es.Rep.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return entitys, nil
}

func (es *EntityService) GetForID(ctx context.Context, id int) (models.Entity, error) {
	entity, err := es.Rep.GetForID(ctx, id)
	if err != nil {
		return models.Entity{}, err
	}
	return entity, nil
}

func (es EntityService) Add(ctx context.Context, obj models.Entity)  error {
	err := es.Rep.Add(ctx, obj)
	return err
}

func (es EntityService) Delete(ctx context.Context, id int) error {
	err := es.Rep.Delete(ctx, id)
	return err
}

func (es EntityService) Update(ctx context.Context, id int, obj models.Entity) error {
	err := es.Rep.Update(ctx, id, obj)
	return err
}
