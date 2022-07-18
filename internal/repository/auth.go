// Package repository repositories app
package repository

import (
	"app/internal/models"
	"context"
	"github.com/google/uuid"
)

// RepoSession Common interface for repository Session
type RepoSession interface {
	Create(ctx context.Context, session *models.Session) error
	Update(ctx context.Context, session *models.Session) error
	Get(ctx context.Context, SessionID uuid.UUID) (*models.Session, error)
	Delete(ctx context.Context, sessionID uuid.UUID) error
	Disable(ctx context.Context, sessionID uuid.UUID) error
}
