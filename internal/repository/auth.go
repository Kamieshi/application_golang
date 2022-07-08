package repository

import (
	"app/internal/models"
	"context"
)

type RepoSession interface {
	Create(ctx context.Context, session *models.Session) error
	Update(ctx context.Context, session *models.Session) error
	Get(ctx context.Context, SessionId string) (*models.Session, error)
	Delete(ctx context.Context, sessionId string) error
}
