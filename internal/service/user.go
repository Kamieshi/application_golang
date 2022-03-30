package service

import (
	"app/internal/repository"
	"app/internal/service/models"
	"context"
	"crypto/sha256"
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
	return fmt.Sprintf("%x\n", h.Sum(nil))
}

func (us UserService) Create(ctx context.Context, userName string, password string) (models.User, error) {
	user, err := us.rep.Get(ctx, userName)
	if err != nil {
		return user, err
	}
	if (user != models.User{}) {
		return models.User{}, err
	}
	user = models.User{
		UserName: userName,
	}
	passwordHash := createHash256Password(user, password)
	user.PasswordHash = passwordHash
	err = us.rep.Create(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return models.User{}, err

}

func (us UserService) Drop(ctx context.Context, username string) error {
	return us.rep.Delete(ctx, username)
}

func (us UserService) Get(ctx context.Context, username string) (models.User, error) {
	return us.rep.Get(ctx, username)
}

