package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"app/internal/models"
	"app/internal/repository"
)

// UserService work with users
type UserService struct {
	rep repository.RepoUser
}

// NewUserService constructor
func NewUserService(rep repository.RepoUser) *UserService {
	return &UserService{
		rep: rep,
	}
}

func createHash256Password(user *models.User, password string) string {
	h := sha256.New()
	h.Write([]byte(user.UserName + password + os.Getenv("SECRET_KEY")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Create user
func (us *UserService) Create(ctx context.Context, userName, password string) (*models.User, error) {
	user, err := us.rep.Get(ctx, userName)
	if user != nil {
		return user, fmt.Errorf("service user/Create : %v", errors.New("username already in use"))
	}
	if err.Error() != "no rows in result set" {
		return user, fmt.Errorf("service user/Create : %v", err)
	}
	err = nil
	user = &models.User{}
	user.UserName = userName
	passwordHash := createHash256Password(user, password)
	user.PasswordHash = passwordHash
	err = us.rep.Add(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service user/Create : %v", err)
	}
	user, _ = us.rep.Get(ctx, userName)
	return user, err
}

// Delete user
func (us *UserService) Delete(ctx context.Context, username string) error {
	return us.rep.Delete(ctx, username)
}

// Get user
func (us *UserService) Get(ctx context.Context, username string) (*models.User, error) {
	return us.rep.Get(ctx, username)
}

// GetAll users
func (us *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := us.rep.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service user/GetAll : %v", err)
	}
	return users, nil
}
