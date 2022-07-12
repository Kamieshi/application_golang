package repository

import (
	"app/internal/models"
	"context"
	"github.com/google/uuid"
)

type RepoSession interface {
	Create(ctx context.Context, session *models.Session) error
	Update(ctx context.Context, session *models.Session) error
	Get(ctx context.Context, SessionId uuid.UUID) (*models.Session, error)
	Delete(ctx context.Context, sessionId uuid.UUID) error
	Disable(ctx context.Context, sessionId uuid.UUID) error
}
