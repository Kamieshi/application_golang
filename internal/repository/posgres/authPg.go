package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"app/internal/models"
)

// RepoAuthPostgres implement RepositorySession
type RepoAuthPostgres struct {
	pool *pgxpool.Pool
}

// NewRepoAuthPostgres Constructor
func NewRepoAuthPostgres(pool *pgxpool.Pool) *RepoAuthPostgres {
	return &RepoAuthPostgres{
		pool: pool,
	}
}

// Create session
func (r RepoAuthPostgres) Create(ctx context.Context, session *models.Session) error {
	session.ID = uuid.New()
	query := "INSERT INTO sessions(id,user_id, refresh_token, signature ,created_at ,disabled) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := r.pool.Exec(ctx, query, session.ID, session.UserID, session.RfToken, session.UniqueSignature, session.CreatedAt, session.Disabled)
	if err != nil {
		return err
	}
	return nil
}

// Update session
func (r RepoAuthPostgres) Update(ctx context.Context, session *models.Session) error {
	query := "UPDATE sessions SET user_id=$1, refresh_token=$2, signature=$3 ,created_at=$4 ,disabled=$5 WHERE id=$6"
	res, err := r.pool.Exec(ctx, query, session.UserID, session.RfToken, session.UniqueSignature, session.CreatedAt, session.Disabled, session.ID)
	if err != nil {
		return err
	}
	if !res.Update() {
		return errors.New("no find entity for ID")
	}
	return nil
}

// Get session
func (r RepoAuthPostgres) Get(ctx context.Context, SessionID uuid.UUID) (*models.Session, error) {
	query := "SELECT id, user_id,  refresh_token, signature ,created_at ,disabled FROM sessions WHERE id=$1"
	var row = r.pool.QueryRow(ctx, query, SessionID)
	var session models.Session
	err := row.Scan(&session.ID, &session.UserID, &session.RfToken, &session.UniqueSignature, &session.CreatedAt, &session.Disabled)
	if err != nil {
		return nil, err
	}
	return &session, err
}

// Delete session
func (r RepoAuthPostgres) Delete(ctx context.Context, sessionID uuid.UUID) error {
	query := "DELETE FROM sessions WHERE id=$1"
	com, err := r.pool.Exec(ctx, query, sessionID)
	if err != nil {
		return err
	}
	if !com.Delete() {
		return errors.New("no find session for ID")
	}
	return nil
}

// Disable Session
func (r RepoAuthPostgres) Disable(ctx context.Context, sessionID uuid.UUID) error {
	query := "UPDATE sessions SET disabled=$1 WHERE id=$2"
	com, err := r.pool.Exec(ctx, query, true, sessionID)
	log.Info(com)
	return err
}
