package repository

import (
	"app/internal/repository/posgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgRepoFactory struct {
	Pool *pgxpool.Pool
}

func (p PgRepoFactory) GetAuthRepo() RepoSession {
	return repository.NewRepoAuthPostgres(p.Pool)
}

func (p PgRepoFactory) GetEntityRepo() RepoEntity {
	return repository.NewRepoEntityPostgres(p.Pool)
}

func (p PgRepoFactory) GetImageRepo() RepoImage {
	return repository.NewRepoImagePostgres(p.Pool)
}

func (p PgRepoFactory) GetUserRepo() RepoUser {
	return repository.NewRepoUsersPostgres(p.Pool)
}
