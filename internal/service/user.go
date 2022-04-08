package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
)

type UserService struct {
	rep repository.UserRepo
}

func NewUserService(rep repository.UserRepo) *UserService {
	return &UserService{
		rep: rep,
	}
}

func createHash256Password(user models.User, password string) string {
	// TODO: Don't Use username

	h := sha256.New()
	h.Write([]byte(user.UserName + password + os.Getenv("SECRET_KEY")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (us UserService) Create(ctx context.Context, userName string, password string) (models.User, error) {
	user, _ := us.rep.Get(ctx, userName)
	if (user != models.User{}) {
		return user, errors.New("username already in use")
	}

	user = models.User{
		UserName: userName,
	}
	passwordHash := createHash256Password(user, password)
	user.PasswordHash = passwordHash
	err := us.rep.Add(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	user, _ = us.rep.Get(ctx, userName)

	return user, err

}

func (us UserService) Delete(ctx context.Context, username string) error {
	return us.rep.Delete(ctx, username)
}

func (us UserService) Get(ctx context.Context, username string) (models.User, error) {
	return us.rep.Get(ctx, username)
}

func (us UserService) GetAll(ctx context.Context) (*[]models.User, error) {

	users, err := us.rep.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
