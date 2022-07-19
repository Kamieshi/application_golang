package repository

import (
	"context"

	"app/internal/models"
)

// RepoUser Common User interface
type RepoUser interface {
	Get(ctx context.Context, username string) (*models.User, error)
	Add(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, username string) error
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
}
