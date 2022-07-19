package repository

import (
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"app/internal/models"
)

func TestRepositoryEntityAdd(t *testing.T) {
	repEntity := NewRepoEntityPostgres(pgPool)
	var entity = models.Entity{
		Name:     "TestEntity1",
		Price:    250,
		IsActive: true,
	}
	err := repEntity.Add(ctx, &entity)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		err = repEntity.Delete(ctx, entity.ID.String())
		if err != nil {
			log.WithError(err).Error()
		}
	})

	entityFromDB, err := repEntity.GetForID(ctx, entity.ID.String())
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(*entityFromDB, entity) == false {
		t.Error("Error")
	}
}

func TestRepositoryEntityGet(t *testing.T) {
	repEntity := NewRepoEntityPostgres(pgPool)
	var entity = models.Entity{
		Name:     "TestEntity1",
		Price:    250,
		IsActive: true,
	}
	err := repEntity.Add(ctx, &entity)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		err = repEntity.Delete(ctx, entity.ID.String())
		if err != nil {
			log.WithError(err).Error()
		}
	})

	entityFromDB, err := repEntity.GetForID(ctx, entity.ID.String())
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(*entityFromDB, entity) == false {
		t.Error("No equal")
	}
}

func TestRepositoryEntityDelete(t *testing.T) {
	repEntity := NewRepoEntityPostgres(pgPool)
	var entity = models.Entity{
		Name:     "TestEntity1",
		Price:    250,
		IsActive: true,
	}
	err := repEntity.Add(ctx, &entity)
	if err != nil {
		t.Error(err)
	}

	err = repEntity.Delete(ctx, entity.ID.String())
	if err != nil {
		t.Error(err)
	}

	checkEntity, err := repEntity.GetForID(ctx, entity.ID.String())
	if err != pgx.ErrNoRows {
		t.Error(err)
	}
	if checkEntity != nil {
		t.Error("Error delete")
	}
}

func TestRepositoryEntityUpdate(t *testing.T) {
	repEntity := NewRepoEntityPostgres(pgPool)
	var entity = models.Entity{
		Name:     "TestEntity1",
		Price:    250,
		IsActive: true,
	}
	err := repEntity.Add(ctx, &entity)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		err = repEntity.Delete(ctx, entity.ID.String())
		if err != nil {
			log.WithError(err).Error()
		}
	})
	newName := "Name name"
	entity.Name = newName
	err = repEntity.Update(ctx, entity.ID.String(), &entity)
	if err != nil {
		t.Error(err)
	}
	checkEntity, err := repEntity.GetForID(ctx, entity.ID.String())
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(*checkEntity, entity) == false {
		t.Error("No equal")
	}
}

func TestRepositoryEntityGetAll(t *testing.T) {
	repEntity := NewRepoEntityPostgres(pgPool)
	var entity1 = models.Entity{
		Name:     "TestEntity1",
		Price:    250,
		IsActive: true,
	}
	var entity2 = models.Entity{
		Name:     "TestEntity1",
		Price:    250,
		IsActive: true,
	}
	err := repEntity.Add(ctx, &entity1)
	if err != nil {
		t.Error(err)
	}
	err = repEntity.Add(ctx, &entity2)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		err = repEntity.Delete(ctx, entity1.ID.String())
		if err != nil {
			log.WithError(err).Error()
		}
		err = repEntity.Delete(ctx, entity2.ID.String())
		if err != nil {
			log.WithError(err).Error()
		}
	})

	entitiesFromDB, err := repEntity.GetAll(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(entitiesFromDB) != 2 {
		t.Error("Uncorrected count rows")
	}
}
