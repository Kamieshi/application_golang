package repository

import (
	"app/internal/models"
	"context"
)

type ImageRepository interface {
	Save(ctx context.Context, img models.Image) error
	Get(ctx context.Context, easyLink string) (models.Image, error)
	Delete(ctx context.Context, img models.Image) error
}
