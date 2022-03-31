package repository

import (
	"app/internal/service/models"
	"context"
)

type Session interface {
	Create(ctx context.Context, session models.Session) error
	Update(ctx context.Context, session models.Session) error
	Get(ctx context.Context, SessionId string) (models.Session, error)
	Delete(ctx context.Context, session string) error
}
