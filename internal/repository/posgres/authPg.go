package repository

import (
	"app/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RepoAuthPostgres struct {
	pool *pgxpool.Pool
}

func NewRepoAuthPostgres(pool *pgxpool.Pool) RepoAuthPostgres {
	return RepoAuthPostgres{
		pool: pool,
	}
}

func rowToSession(row pgx.Row) (*models.Session, error) {
	var session models.Session
	err := row.Scan(&session.Id, &session.UserId, &session.IdSession, &session.RfToken, &session.UniqueSignature, &session.CreatedAt, &session.Disabled)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

const orderColumnsSessions string = "id, user_id, session_id, refresh_token, signature ,created_at ,disabled"

func (a RepoAuthPostgres) Create(ctx context.Context, session models.Session) error {
	query := "INSERT INTO sessions(user_id, session_id, refresh_token, signature ,created_at ,disabled) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := a.pool.Exec(ctx, query, session.UserId, session.IdSession, session.RfToken, session.UniqueSignature, session.CreatedAt, session.Disabled)
	if err != nil {
		return err
	}
	return nil
}

func (a RepoAuthPostgres) Update(ctx context.Context, session models.Session) error {
	query := "UPDATE sessions SET user_id=$1, session_id=$2, refresh_token=$3, signature=$4 ,created_at=$5 ,disabled=$6 WHERE id=$7"
	res, err := a.pool.Exec(ctx, query, session.UserId, session.IdSession, session.RfToken, session.UniqueSignature, session.CreatedAt, session.Disabled, session.Id)
	if err != nil {
		return err
	}
	if res.String() == "UPDATE 0" {
		return errors.New("no find entity for ID")
	}
	return nil
}

func (a RepoAuthPostgres) Get(ctx context.Context, SessionId string) (models.Session, error) {
	query := fmt.Sprintf("SELECT %s FROM sessions WHERE session_id=$1", orderColumnsSessions)
	var row pgx.Row = a.pool.QueryRow(ctx, query, SessionId)
	ent, err := rowToSession(row)
	return *ent, err
}

func (a RepoAuthPostgres) Delete(ctx context.Context, sessionId string) error {

	query := "DELETE FROM sessions WHERE session_id=$1"
	com, err := a.pool.Exec(ctx, query, sessionId)

	if err != nil {
		return err
	}
	if com.String() == "DELETE 0" {
		return errors.New("no find session for ID")
	}

	return nil
}
