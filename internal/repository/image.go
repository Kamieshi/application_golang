package repository

import (
	"app/internal/models"
	"context"
	"github.com/google/uuid"
)

type RepoImage interface {
	Save(ctx context.Context, img *models.Image) error
	Get(ctx context.Context, easyLink string) (*models.Image, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
