package repository

import (
	"app/internal/models"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RepoImagePostgres struct {
	pool *pgxpool.Pool
}

func NewRepoImagePostgres(pool *pgxpool.Pool) *RepoImagePostgres {
	return &RepoImagePostgres{
		pool: pool,
	}
}

func (r RepoImagePostgres) Save(ctx context.Context, img *models.Image) error {
	query := "INSERT INTO images(id,file_name,root_path,easy_link) values ($1,$2,$3,$4)"
	img.ID = uuid.New()
	_, err := r.pool.Exec(ctx, query, img.ID, img.Filename, img.RootPath, img.EasyLink)
	if err != nil {
		img.ID = uuid.UUID{}
	}
	return err
}

func (r RepoImagePostgres) Get(ctx context.Context, easyLink string) (*models.Image, error) {
	var image models.Image
	query := "SELECT id,file_name,root_path,easy_link FROM images WHERE easy_link=$1"
	row := r.pool.QueryRow(ctx, query, easyLink)

	err := row.Scan(&image.ID, &image.Filename, &image.RootPath, &image.EasyLink)
	if err != nil {
		return &image, err
	}
	image.Data, err = image.Byte()
	return &image, err
}

func (r RepoImagePostgres) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
