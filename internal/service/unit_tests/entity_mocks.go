package unit_tests

import (
	"app/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type respEntityMock struct {
	mock.Mock
}

func (r respEntityMock) GetAll(ctx context.Context) ([]*models.Entity, error) {
	args := r.Called()
	return args.Get(0).([]*models.Entity), args.Error(1)
}

func (r respEntityMock) GetForID(ctx context.Context, id string) (*models.Entity, error) {
	args := r.Called(id)
	return args.Get(0).(*models.Entity), args.Error(1)
}

func (r respEntityMock) Add(ctx context.Context, obj *models.Entity) error {
	return nil
}

func (r respEntityMock) Update(ctx context.Context, id string, obj *models.Entity) error {
	return nil
}

func (r respEntityMock) Delete(ctx context.Context, id string) error {
	return nil
}

type cashMock struct {
	mock.Mock
}

func (c cashMock) Set(ctx context.Context, entity *models.Entity) error {
	return nil
}

func (c cashMock) Get(ctx context.Context, idEntity string) (*models.Entity, bool) {
	args := c.Called(idEntity)
	return args.Get(0).(*models.Entity), args.Bool(1)
}

func (c cashMock) Delete(ctx context.Context, idEntity string) {
}
