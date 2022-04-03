package repository

import (
	"app/internal/models"
	"context"
)

// type ControllerEntity interface {
// 	GetAllItems(ctx context.Context) ([]models.Entity, error)
// 	GetItemForID(ctx context.Context, id int) (models.Entity, error)
// 	AddItem(ctx context.Context, obj models.Entity) (bool, error)
// 	UpdateItem(ctx context.Context, id int, obj models.Entity) error
// 	DeleteItem(ctx context.Context, id int) error
// }

type RepoEntity interface {
	GetAll(ctx context.Context) (*[]models.Entity, error)
	GetForID(ctx context.Context, id string) (models.Entity, error)
	Add(ctx context.Context, obj models.Entity) error
	Update(ctx context.Context, id string, obj models.Entity) error
	Delete(ctx context.Context, id string) error
}
