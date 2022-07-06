package repository

import (
	"app/internal/models"
	"context"
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

func (a RepoAuthPostgres) Create(ctx context.Context, session models.Session) error {
	//TODO implement me
	panic("implement me")
}

func (a RepoAuthPostgres) Update(ctx context.Context, session models.Session) error {
	//TODO implement me
	panic("implement me")
}

func (a RepoAuthPostgres) Get(ctx context.Context, SessionId string) (models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (a RepoAuthPostgres) Delete(ctx context.Context, session string) error {
	//TODO implement me
	panic("implement me")
}
