package service

import (
	"app/internal/repository"
	"app/internal/service/models"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Auth struct {
	UserRep repository.UserRepo
	AuthRep repository.Session
}

type JwtPayLoadClaims struct {
	Username  string `json:"name"`
	Admin     bool   `json:"admin"`
	IdSession string `json:"id_session"`
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

func (au Auth) CreateToken(username string, admin bool, idSession string) (string, error) {
	payLoad := &JwtPayLoadClaims{
		username,
		admin,
		idSession,
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

func (au Auth) CreateAndWriteSession(ctx echo.Context, user models.User) (models.Session, error) {
	resrash := au.createRandomOutput(user.UserName)

	session := models.Session{
		IdSession:       au.createRandomOutput("id"),
		UserId:          user.Id,
		CreatedAt:       time.Now(),
		Disabled:        false,
		RfToken:         resrash,
		UniqueSignature: ctx.Request().UserAgent(),
	}
	err := au.AuthRep.Create(ctx.Request().Context(), session)
	if err != nil {
		return models.Session{}, err
	}
	return session, err
}

func (au Auth) RefrashAndWriteSession(ctx echo.Context, rfToken string) (string, string, error) {
	user := ctx.Get("user").(*jwt.Token)
	payLoad := user.Claims.(*JwtPayLoadClaims)
	currentSession, err := au.AuthRep.Get(ctx.Request().Context(), payLoad.IdSession)
	if err != nil {
		return "", "", err
	}
	if currentSession.Disabled {
		return "", "", errors.New("Disable session")
	}
	if rfToken == currentSession.RfToken {
		accessToken, _ := au.CreateToken(payLoad.Username, payLoad.Admin, payLoad.IdSession)
		currentSession.RfToken = au.createRandomOutput()
		au.AuthRep.Update(ctx.Request().Context(), currentSession)
		return accessToken, currentSession.RfToken, nil
	}
	return "", "", errors.New("Disable session")
}

func (au Auth) createRandomOutput(sal ...string) string {
	data := fmt.Sprint(time.Now().Unix(), os.Getenv("SECRET_KEY"), rand.Int31n(1000), sal)
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (au Auth) GetUser(ctx echo.Context) (models.User, error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtPayLoadClaims)
	User, err := au.UserRep.Get(ctx.Request().Context(), claims.Username)
	return User, err
}
