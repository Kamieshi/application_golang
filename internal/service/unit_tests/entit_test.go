package unit_tests

import (
	"app/internal/models"
	"app/internal/service"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAll(t *testing.T) {
	respMock := new(respEntityMock)
	cachMock := new(cashMock)
	ctx := context.Background()
	assert := assert.New(t)
	dataForAllQuery := []*models.Entity{
		&models.Entity{
			ID:       uuid.New(),
			Name:     "1",
			Price:    0,
			IsActive: false,
		},
		&models.Entity{
			ID:       uuid.New(),
			Name:     "2",
			Price:    0,
			IsActive: true,
		},
	}
	respMock.On("GetAll").Return(dataForAllQuery, nil)

	entityService := service.NewEntityService(respMock, cachMock)
	entitiesActual, err := entityService.GetAll(ctx)
	assert.Nil(err)
	assert.Equal(dataForAllQuery, entitiesActual)
}

func TestGetForID(t *testing.T) {
	respMock := new(respEntityMock)
	cachMock := new(cashMock)
	entityService := service.NewEntityService(respMock, cachMock)
	ctx := context.Background()
	assert := assert.New(t)
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
	respMock.On("GetForID", expectEntity1.ID.String()).Return(expectEntity1, nil)
	respMock.On("GetForID", expectEntity2.ID.String()).Return(expectEntity2, nil)

	cachMock.On("Get", expectEntity1.ID.String()).Return(expectEntity1, true)

	var emptyEntity *models.Entity
	cachMock.On("Get", expectEntity2.ID.String()).Return(emptyEntity, false)

	actualEntity1, err := entityService.GetForID(ctx, expectEntity1.ID.String())
	assert.Nil(err)
	assert.Equal(expectEntity1, actualEntity1)

	actualEntity2, err := entityService.GetForID(ctx, expectEntity2.ID.String())
	assert.Nil(err)
	assert.Equal(expectEntity2, actualEntity2)
}
