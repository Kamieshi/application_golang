package service

import (
	"app/internal/models"
	repoMock "app/internal/repository/mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var ctx context.Context
var errorFromRepo error
var errorFormCache error

func TestMain(t *testing.M) {
	errorFromRepo = errors.New("Error from repositpory")
	errorFormCache = errors.New("Error form cache")
	ctx = context.Background()
	code := t.Run()
	os.Exit(code)
}

func TestGetAllPositive(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	dataForAllQuery := []*models.Entity{
		{
			ID:       uuid.New(),
			Name:     "1",
			Price:    0,
			IsActive: false,
		},
		{
			ID:       uuid.New(),
			Name:     "2",
			Price:    0,
			IsActive: true,
		},
	}

	repoEntityMock.On("GetAll", ctx).Return(dataForAllQuery, nil)

	entitiesActual, err := entityService.GetAll(ctx)
	_assert.Nil(err)
	_assert.Equal(dataForAllQuery, entitiesActual)
}

func TestGetForIDPositive(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		ID:       uuid.New(),
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	expectEntity2 := &models.Entity{
		ID:       uuid.New(),
		Name:     "Entity 2",
		Price:    0,
		IsActive: false,
	}
	var emptyEntity *models.Entity
	repoEntityMock.On("GetForID", ctx, expectEntity1.ID.String()).Return(expectEntity1, nil)
	repoEntityMock.On("GetForID", ctx, expectEntity2.ID.String()).Return(expectEntity2, nil)

	repoCacheMock.On("Get", ctx, expectEntity1.ID.String()).Return(expectEntity1, true)
	repoCacheMock.On("Get", ctx, expectEntity2.ID.String()).Return(emptyEntity, false)
	repoCacheMock.On("Set", ctx, expectEntity2).Return(nil)

	actualEntity1, err := entityService.GetForID(ctx, expectEntity1.ID.String())
	_assert.Nil(err)
	_assert.Equal(expectEntity1, actualEntity1)

	actualEntity2, err := entityService.GetForID(ctx, expectEntity2.ID.String())
	_assert.Nil(err)
	_assert.Equal(expectEntity2, actualEntity2)
}

func TestUpdatePositive(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		ID:       uuid.New(),
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	repoEntityMock.On("Update", ctx, expectEntity1.ID.String(), expectEntity1).Return(nil)
	repoCacheMock.On("Delete", ctx, expectEntity1.ID.String()).Return()
	err := entityService.Update(ctx, expectEntity1.ID.String(), expectEntity1)
	_assert.Nil(err)
}

func TestAddPositive(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	repoEntityMock.On("Add", ctx, expectEntity1).Return(nil)
	repoCacheMock.On("Set", ctx, expectEntity1).Return(nil)
	err := entityService.Add(ctx, expectEntity1)
	_assert.Nil(err)
}

func TestDeletePositive(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		ID:       uuid.New(),
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	repoEntityMock.On("Delete", ctx, expectEntity1.ID.String()).Return(nil)
	repoCacheMock.On("Delete", ctx, expectEntity1.ID.String()).Return()
	err := entityService.Delete(ctx, expectEntity1.ID.String())

	_assert.Nil(err)
}

func TestGetAllNegative(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	dataForAllQuery := []*models.Entity{
		{
			ID:       uuid.New(),
			Name:     "1",
			Price:    0,
			IsActive: false,
		},
		{
			ID:       uuid.New(),
			Name:     "2",
			Price:    0,
			IsActive: true,
		},
	}

	repoEntityMock.On("GetAll", ctx).Return(dataForAllQuery, errorFromRepo)

	entitiesActual, err := entityService.GetAll(ctx)
	_assert.Error(err)
	_assert.NotEqual(dataForAllQuery, entitiesActual)
}

func TestGetForIDNegative(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		ID:       uuid.New(),
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	var emptyEntity *models.Entity
	repoEntityMock.On("GetForID", ctx, expectEntity1.ID.String()).Return(emptyEntity, errorFromRepo)
	repoCacheMock.On("Get", ctx, expectEntity1.ID.String()).Return(emptyEntity, false)

	actualEntity1, err := entityService.GetForID(ctx, expectEntity1.ID.String())
	_assert.Error(err)
	_assert.NotEqual(expectEntity1, actualEntity1)

}

func TestUpdateNegative(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		ID:       uuid.New(),
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	repoEntityMock.On("Update", ctx, expectEntity1.ID.String(), expectEntity1).Return(errorFromRepo)
	repoCacheMock.On("Delete", ctx, expectEntity1.ID.String()).Return()
	err := entityService.Update(ctx, expectEntity1.ID.String(), expectEntity1)
	_assert.Error(err)
}

func TestAddNegative(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	repoEntityMock.On("Add", ctx, expectEntity1).Return(errorFromRepo)
	repoCacheMock.On("Set", ctx, expectEntity1).Return(errors.New("Error form cache"))
	err := entityService.Add(ctx, expectEntity1)
	_assert.Error(err)
}

func TestDeleteNegative(t *testing.T) {
	repoEntityMock := new(repoMock.RepoEntity)
	repoCacheMock := new(repoMock.CacheEntityRepository)
	entityService := NewEntityService(repoEntityMock, repoCacheMock)
	_assert := assert.New(t)
	expectEntity1 := &models.Entity{
		ID:       uuid.New(),
		Name:     "Entity 1",
		Price:    0,
		IsActive: false,
	}
	repoEntityMock.On("Delete", ctx, expectEntity1.ID.String()).Return(errorFromRepo)
	repoCacheMock.On("Delete", ctx, expectEntity1.ID.String()).Return(errorFormCache)
	err := entityService.Delete(ctx, expectEntity1.ID.String())

	_assert.Error(err)
}
