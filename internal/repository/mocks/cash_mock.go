// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "app/internal/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CacheEntityRepository is an autogenerated mock type for the CacheEntityRepository type
type CacheEntityRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, idEntity
func (_m *CacheEntityRepository) Delete(ctx context.Context, idEntity string) {
	_m.Called(ctx, idEntity)
}

// Get provides a mock function with given fields: ctx, idEntity
func (_m *CacheEntityRepository) Get(ctx context.Context, idEntity string) (*models.Entity, bool) {
	ret := _m.Called(ctx, idEntity)

	var r0 *models.Entity
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Entity); ok {
		r0 = rf(ctx, idEntity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Entity)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, idEntity)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Set provides a mock function with given fields: ctx, entity
func (_m *CacheEntityRepository) Set(ctx context.Context, entity *models.Entity) error {
	ret := _m.Called(ctx, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Entity) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCacheEntityRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCacheEntityRepository creates a new instance of CacheEntityRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCacheEntityRepository(t mockConstructorTestingTNewCacheEntityRepository) *CacheEntityRepository {
	mock := &CacheEntityRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}