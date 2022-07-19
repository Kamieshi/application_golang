package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"app/internal/models"
)

// RepoUsersPostgres implement RepositoryUser
type RepoUsersPostgres struct {
	pool *pgxpool.Pool
}

// NewRepoUsersPostgres constructor
func NewRepoUsersPostgres(pool *pgxpool.Pool) *RepoUsersPostgres {
	return &RepoUsersPostgres{
		pool: pool,
	}
}

// Get user
func (r RepoUsersPostgres) Get(ctx context.Context, username string) (*models.User, error) {
	query := "SELECT id, username, password_hash, is_admin FROM users WHERE username=$1"
	row := r.pool.QueryRow(ctx, query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.UserName, &user.PasswordHash, &user.Admin)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Add user
func (r RepoUsersPostgres) Add(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	query := "INSERT INTO users( id,username, password_hash, is_admin) values ($1,$2,$3,$4)"
	_, err := r.pool.Exec(ctx, query, user.ID, user.UserName, user.PasswordHash, user.Admin)
	return err
}

// Delete user
func (r RepoUsersPostgres) Delete(ctx context.Context, username string) error {
	query := "DELETE FROM users WHERE username=$1"
	com, err := r.pool.Exec(ctx, query, username)
	if err != nil {
		return err
	}
	if com.String() == "DELETE 0" {
		return errors.New("no find entity for username")
	}
	return nil
}

// GetAll users
func (r RepoUsersPostgres) GetAll(ctx context.Context) ([]*models.User, error) {
	query := "SELECT id, username, password_hash, is_admin FROM users"
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users = make([]*models.User, 0, len(rows.FieldDescriptions()))

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.UserName, &user.PasswordHash, &user.Admin)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// Update user
func (r RepoUsersPostgres) Update(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET username=$1, password_hash=$2, is_admin=$3 WHERE id=$4"
	com, err := r.pool.Exec(ctx, query, user.UserName, user.PasswordHash, user.Admin, user.ID)
	if err != nil {
		return err
	}
	if com.String() == "UPDATE 0" {
		return errors.New("no find entity for ID")
	}
	return nil
}
