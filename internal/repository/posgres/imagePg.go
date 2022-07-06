package repository

import (
	"app/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RepoImagePostgres struct {
	pool *pgxpool.Pool
}

func NewRepoImagePostgres(pool *pgxpool.Pool) RepoImagePostgres {
	return RepoImagePostgres{
		pool: pool,
	}
}

func (r RepoImagePostgres) Save(ctx context.Context, img models.Image) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepoImagePostgres) Get(ctx context.Context, easyLink string) (*models.Image, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepoImagePostgres) Delete(ctx context.Context, id interface{}) error {
	//TODO implement me
	panic("implement me")
}
