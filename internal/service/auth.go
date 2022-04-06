package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Auth struct {
	UserRep   repository.UserRepo
	AuthRep   repository.Session
	JWTConfig middleware.JWTConfig
}

func NewAuth(userRep repository.UserRepo, sessionRep repository.Session) Auth {
	headerAuthorization := echo.HeaderAuthorization
	config := middleware.JWTConfig{
		Claims:        &CustomClaims{},
		SigningKey:    []byte(os.Getenv("SECRET_KEY")),
		AuthScheme:    "Bearer",
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   "header:" + headerAuthorization,
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return []byte(os.Getenv("SECRET_KEY")), nil
			}

			t := reflect.ValueOf(&CustomClaims{}).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err := jwt.ParseWithClaims(auth, claims, keyFunc)

			if err != nil {
				if fmt.Sprint(err)[:16] == "token is expired" && c.Path() == "/auth/refresh" {
					return token, nil
				}
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
		Skipper: func(c echo.Context) bool {
			logrus.WithFields(logrus.Fields{"path": c.Path()}).Info("Skipper JWT Auth")
			switch c.Path() {
			case "/auth/login":
				return true

			case "/user":
				if c.Request().Method == http.MethodPost {
					return true
				}
			case "/images*":
				return true
			case "/upload":
				return true
			case "/load/:easy_link":
				return true
			}

			return false
		},
	}

	return Auth{
		UserRep:   userRep,
		AuthRep:   sessionRep,
		JWTConfig: config,
	}
}

type CustomClaims struct {
	Username  string `json:"name"`
	Admin     bool   `json:"admin"`
	IdSession string `json:"id_session"`
	jwt.StandardClaims
}

func (au Auth) IsAuthentication(ctx context.Context, username string, password string) (models.User, bool, error) {
	user, err := au.UserRep.Get(ctx, username)
	if err != nil {
		logrus.WithFields(logrus.Fields{"username": username}).Warn("Unsuccessful login attempt")
		return user, false, err
	}
	inPasswordHash := createHash256Password(user, password)
	if user.PasswordHash == inPasswordHash {
		return user, true, err
	}
	return models.User{}, false, err
}

func (au Auth) CreateToken(username string, admin bool, idSession string) (string, error) {
	payLoad := &CustomClaims{
		username,
		admin,
		idSession,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(),
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
	refresh := au.createRandomOutput(user.UserName)

	session := models.Session{
		IdSession:       au.createRandomOutput("id"),
		UserId:          user.Id,
		CreatedAt:       time.Now(),
		Disabled:        false,
		RfToken:         createHashSHA256WithSalt(refresh),
		UniqueSignature: ctx.Request().UserAgent(),
	}
	err := au.AuthRep.Create(ctx.Request().Context(), session)
	if err != nil {
		return models.Session{}, err
	}
	session.RfToken = refresh
	return session, err
}

func (au Auth) RefreshAndWriteSession(ctx echo.Context, rfToken string) (string, string, error) {
	user := ctx.Get("user").(*jwt.Token)
	payLoad := user.Claims.(*CustomClaims)

	currentSession, err := au.AuthRep.Get(ctx.Request().Context(), payLoad.IdSession)
	if err != nil {
		return "", "", err
	}
	if currentSession.Disabled {
		return "", "", errors.New("Disable session")
	}
	if createHashSHA256WithSalt(rfToken) == currentSession.RfToken {
		accessToken, _ := au.CreateToken(payLoad.Username, payLoad.Admin, payLoad.IdSession)
		newRfToken := au.createRandomOutput()
		currentSession.RfToken = createHashSHA256WithSalt(newRfToken)
		au.AuthRep.Update(ctx.Request().Context(), currentSession)
		return accessToken, newRfToken, nil
	}
	return "", "", errors.New("Disable session")
}

func (au Auth) createRandomOutput(sal ...string) string {
	data := fmt.Sprint(time.Now().Unix(), os.Getenv("SECRET_KEY"), rand.Int31n(1000), sal)
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func createHashSHA256WithSalt(s string) string {
	data := fmt.Sprint(s, os.Getenv("SECRET_KEY"))
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (au Auth) GetUser(ctx echo.Context) (models.User, error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*CustomClaims)
	User, err := au.UserRep.Get(ctx.Request().Context(), claims.Username)
	return User, err
}
