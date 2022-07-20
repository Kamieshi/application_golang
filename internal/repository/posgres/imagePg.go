package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"app/internal/models"
)

// RepoImagePostgres Implement RepositoryImage like Postgres
type RepoImagePostgres struct {
	pool *pgxpool.Pool
}

// NewRepoImagePostgres Constructor
func NewRepoImagePostgres(pool *pgxpool.Pool) *RepoImagePostgres {
	return &RepoImagePostgres{
		pool: pool,
	}
}

// Save image
func (r RepoImagePostgres) Save(ctx context.Context, img *models.Image) error {
	query := "INSERT INTO images(id,file_name,root_path,easy_link) values ($1,$2,$3,$4)"
	img.ID = uuid.New()
	_, err := r.pool.Exec(ctx, query, img.ID, img.Filename, img.RootPath, img.EasyLink)
	if err != nil {
		img.ID = uuid.UUID{}
		return fmt.Errorf("imagePg.go/Save : %v", err)
	}
	return err
}

// Get image
func (r RepoImagePostgres) Get(ctx context.Context, easyLink string) (*models.Image, error) {
	var image models.Image
	query := "SELECT id,file_name,root_path,easy_link FROM images WHERE easy_link=$1"
	row := r.pool.QueryRow(ctx, query, easyLink)

	err := row.Scan(&image.ID, &image.Filename, &image.RootPath, &image.EasyLink)
	if err != nil {
		return &image, fmt.Errorf("imagePg.go/Get : %v", err)
	}
	image.Data, err = image.Byte()
	if err != nil {
		return nil, fmt.Errorf("imagePg.go/Get : %v", err)
	}
	return &image, err
}

// Delete image
func (r RepoImagePostgres) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
