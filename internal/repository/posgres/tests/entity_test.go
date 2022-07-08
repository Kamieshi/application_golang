package tests

import (
	"app/internal/models"
	repository "app/internal/repository/posgres"
	"reflect"
	"testing"
)

func TestRepositoryEntityAdd(t *testing.T) {
	repEntity := repository.NewRepoEntityPostgres(pgPool)
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
		repEntity.Delete(ctx, entity.ID.String())
	})

	entityFromDb, err := repEntity.GetForID(ctx, entity.ID.String())
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(*entityFromDb, entity) == false {
		t.Error("Error")
	}
}

func TestRepositoryEntityGet(t *testing.T) {
	repEntity := repository.NewRepoEntityPostgres(pgPool)
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
		repEntity.Delete(ctx, entity.ID.String())
	})

	entityFromDb, err := repEntity.GetForID(ctx, entity.ID.String())
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(*entityFromDb, entity) == false {
		t.Error("No equal")
	}
}

func TestRepositoryEntityDelete(t *testing.T) {
	repEntity := repository.NewRepoEntityPostgres(pgPool)
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
	if checkEntity != nil {
		t.Error("Error delete")
	}
}

func TestRepositoryEntityUpdate(t *testing.T) {
	repEntity := repository.NewRepoEntityPostgres(pgPool)
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
		repEntity.Delete(ctx, entity.ID.String())
	})

	entity.Name = "New name"
	err = repEntity.Update(ctx, entity.ID.String(), entity)
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
	repEntity := repository.NewRepoEntityPostgres(pgPool)
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
		repEntity.Delete(ctx, entity1.ID.String())
		repEntity.Delete(ctx, entity2.ID.String())
	})

	entitiesFromDB, err := repEntity.GetAll(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(entitiesFromDB) != 2 {
		t.Error("Uncorrected count rows")
	}
}
