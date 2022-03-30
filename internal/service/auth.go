package service

import (
	"app/internal/repository"
	"context"
)

type Auth struct {
	UserRep repository.UserRepo
}

func (au Auth) IsAuthentication(ctx context.Context, usename string, password string) (bool, error) {
	user, err := au.UserRep.Get(ctx, usename)
	if err != nil {
		return false, err
	}
	inPasswordhash := createHash256Password(user, password)
	if user.PasswordHash == inPasswordhash {
		return true, err
	}
	return false, err
}
