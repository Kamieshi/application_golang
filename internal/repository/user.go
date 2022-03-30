package repository

import (
	"app/internal/service/models"
	"context"
)

type UserRepo interface {
	Get(ctx context.Context, username string) (models.User, error)
	Create(ctx context.Context, user models.User) error
	Delete(ctx context.Context, username string) error
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context) error
}
