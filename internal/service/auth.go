package service

import (
	"app/internal/repository"
	"app/internal/service/models"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Auth struct {
	UserRep repository.UserRepo
}

type JwtPayLoadClaims struct {
	Username string `json:"name"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

func (au Auth) JWTConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:      &JwtPayLoadClaims{},
		TokenLookup: "cookie:token",
		SigningKey:  []byte(os.Getenv("SECRET_KEY")),
	}
}

func (au Auth) IsAuthentication(ctx context.Context, usename string, password string) (models.User, bool, error) {
	user, err := au.UserRep.Get(ctx, usename)
	if err != nil {
		return user, false, err
	}
	inPasswordhash := createHash256Password(user, password)
	if user.PasswordHash == inPasswordhash {
		return user, true, err
	}
	return models.User{}, false, err
}

func (au Auth) CreateToken(user models.User) (string, error) {
	payLoad := &JwtPayLoadClaims{
		user.UserName,
		user.Admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payLoad)
	tt, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	fmt.Println(token)
	return tt, nil
}

func (au Auth) GetUser(ctx echo.Context) (models.User, error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtPayLoadClaims)
	User, err := au.UserRep.Get(ctx.Request().Context(), claims.Username)
	return User, err
}
