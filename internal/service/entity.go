package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EntityService struct {
	rep      repository.RepoEntity
	cashRep  repository.CacheEntityRepository
	UseCache bool
}

func NewEntityService(rep repository.RepoEntity, cahRep repository.CacheEntityRepository) EntityService {
	if cahRep != nil {
		return EntityService{
			rep:      rep,
			cashRep:  cahRep,
			UseCache: true,
		}
	}
	return EntityService{
		rep:      rep,
		UseCache: false,
	}
}

func (es EntityService) GetAll(ctx context.Context) (*[]models.Entity, error) {

	entities, err := es.rep.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (es *EntityService) GetForID(ctx context.Context, id string) (models.Entity, error) {
	if es.UseCache {
		entity, err := es.cashRep.Get(ctx, id)
		if err != nil {

			logrus.WithFields(logrus.Fields{
				"id": id,
			}).Info("entity not found")

			entity, err = es.rep.GetForID(ctx, id)
			if err != nil {
				return models.Entity{}, err
			}

			err = es.cashRep.Set(ctx, &entity)
			if err != nil {
				logrus.WithError(err)
			}
			logrus.Info("Value successful get in cash repository ")
			logrus.Info("use db")
			return entity, err
		}
		logrus.Info("use cache")
		return entity, err
	}
	entity, err := es.rep.GetForID(ctx, id)
	if err != nil {
		return models.Entity{}, err
	}
	logrus.Info("use db")
	return entity, nil

}

func (es EntityService) Add(ctx context.Context, obj *models.Entity) error {
	err := es.rep.Add(ctx, obj)
	if err != nil {
		return err
	}

	if es.UseCache {

		errSet := es.cashRep.Set(ctx, obj)
		if errSet != nil {
			logrus.WithError(errSet)
		}
		logrus.Info("Value successful get in cash repository ")
	}

	return err
}

func (es EntityService) Delete(ctx context.Context, id string) error {
	err := es.rep.Delete(ctx, id)
	if err != nil {
		return err
	}

	if es.UseCache {

		errDel := es.cashRep.Delete(ctx, id)
		if errDel != nil {
			logrus.WithError(errDel).Error("Delete from cache")
		}
	}
	return err
}

func (es EntityService) Update(ctx context.Context, id string, obj models.Entity) error {
	err := es.rep.Update(ctx, id, obj)
	if err != nil {
		return err
	}
	if es.UseCache {
		objId, _ := primitive.ObjectIDFromHex(id)
		obj.Id = objId
		errUpdate := es.cashRep.Set(ctx, &obj)
		if errUpdate != nil {
			logrus.WithError(errUpdate).Error("Update entity")
		}
	}
	return err
}
