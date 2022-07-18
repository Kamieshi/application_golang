package repository

import (
	"app/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type RepoAuthPostgres struct {
	pool *pgxpool.Pool
}

func NewRepoAuthPostgres(pool *pgxpool.Pool) *RepoAuthPostgres {
	return &RepoAuthPostgres{
		pool: pool,
	}
}

func rowToSession(row pgx.Row) (*models.Session, error) {
	var session models.Session
	err := row.Scan(&session.ID, &session.UserID, &session.RfToken, &session.UniqueSignature, &session.CreatedAt, &session.Disabled)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

const orderColumnsSessions string = "id, user_id,  refresh_token, signature ,created_at ,disabled"

func (r RepoAuthPostgres) Create(ctx context.Context, session *models.Session) error {
	session.ID = uuid.New()
	query := "INSERT INTO sessions(id,user_id, refresh_token, signature ,created_at ,disabled) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := r.pool.Exec(ctx, query, session.ID, session.UserID, session.RfToken, session.UniqueSignature, session.CreatedAt, session.Disabled)
	if err != nil {
		return err
	}
	return nil
}

func (r RepoAuthPostgres) Update(ctx context.Context, session *models.Session) error {
	query := "UPDATE sessions SET user_id=$1, refresh_token=$2, signature=$3 ,created_at=$4 ,disabled=$5 WHERE id=$6"
	res, err := r.pool.Exec(ctx, query, session.UserID, session.RfToken, session.UniqueSignature, session.CreatedAt, session.Disabled, session.ID)
	if err != nil {
		return err
	}
	if res.String() == "UPDATE 0" {
		return errors.New("no find entity for ID")
	}
	return nil
}

func (r RepoAuthPostgres) Get(ctx context.Context, SessionId uuid.UUID) (*models.Session, error) {
	query := fmt.Sprintf("SELECT %s FROM sessions WHERE id=$1", orderColumnsSessions)
	var row pgx.Row = r.pool.QueryRow(ctx, query, SessionId)
	ent, err := rowToSession(row)
	return ent, err
}

func (r RepoAuthPostgres) Delete(ctx context.Context, sessionId uuid.UUID) error {

	query := "DELETE FROM sessions WHERE id=$1"
	com, err := r.pool.Exec(ctx, query, sessionId)

	if err != nil {
		return err
	}
	if com.String() == "DELETE 0" {
		return errors.New("no find session for ID")
	}

	return nil
}

func (r RepoAuthPostgres) Disable(ctx context.Context, sessionId uuid.UUID) error {
	query := "UPDATE sessions SET disabled=$1 WHERE id=$2"
	com, err := r.pool.Exec(ctx, query, true, sessionId)
	log.Info(com)
	return err
}
