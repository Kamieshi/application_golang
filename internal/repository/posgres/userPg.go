package repository

import (
	"app/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RepoUsersPostgres struct {
	pool *pgxpool.Pool
}

func NewRepoUsersPostgres(pool *pgxpool.Pool) RepoUsersPostgres {
	return RepoUsersPostgres{
		pool: pool,
	}
}

func (r RepoUsersPostgres) Get(ctx context.Context, username string) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepoUsersPostgres) Add(ctx context.Context, user models.User) error {
	//TODO implement me
	panic("implement me")
}

func (r RepoUsersPostgres) Delete(ctx context.Context, username string) error {
	//TODO implement me
	panic("implement me")
}

func (r RepoUsersPostgres) GetAll(ctx context.Context) (*[]models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepoUsersPostgres) Update(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
