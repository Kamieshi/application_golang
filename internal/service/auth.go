package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AuthService struct {
	UserRep   repository.RepoUser
	AuthRep   repository.RepoSession
	JWTConfig middleware.JWTConfig
}

func NewAuthService(userRep repository.RepoUser, sessionRep repository.RepoSession) AuthService {
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
			case "/swagger/*":
				return true
			}

			return false
		},
	}

	return AuthService{
		UserRep:   userRep,
		AuthRep:   sessionRep,
		JWTConfig: config,
	}
}

type CustomClaims struct {
	Username  string    `json:"name"`
	Admin     bool      `json:"admin"`
	IdSession uuid.UUID `json:"id_session"`
	jwt.StandardClaims
}

func (a AuthService) IsAuthentication(ctx context.Context, username string, password string) (*models.User, bool, error) {
	user, err := a.UserRep.Get(ctx, username)
	if err != nil {
		logrus.WithFields(logrus.Fields{"username": username}).Warn("Unsuccessful login attempt")
		return user, false, err
	}
	inPasswordHash := createHash256Password(user, password)
	if user.PasswordHash == inPasswordHash {
		return user, true, err
	}
	return nil, false, err
}

func (a AuthService) CreateToken(username string, admin bool, idSession uuid.UUID) (string, error) {

	timeLive, _ := strconv.Atoi(os.Getenv("TIME_LIVE_MINUTE_JWT"))

	payLoad := &CustomClaims{
		username,
		admin,
		idSession,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(timeLive)).Unix(),
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

func (a AuthService) CreateAndWriteSession(ctx echo.Context, user models.User) (models.Session, error) {
	refresh := a.createRandomOutput(user.UserName)

	session := models.Session{
		ID:              uuid.New(),
		UserId:          user.ID,
		CreatedAt:       time.Now(),
		Disabled:        false,
		RfToken:         createHashSHA256WithSalt(refresh),
		UniqueSignature: ctx.Request().UserAgent(),
	}
	err := a.AuthRep.Create(ctx.Request().Context(), &session)
	if err != nil {
		return models.Session{}, err
	}
	return session, err
}

func (a AuthService) RefreshAndWriteSession(ctx echo.Context, rfToken string) (string, string, error) {
	user := ctx.Get("user").(*jwt.Token)
	payLoad := user.Claims.(*CustomClaims)

	currentSession, err := a.AuthRep.Get(ctx.Request().Context(), payLoad.IdSession)
	if currentSession.Disabled {
		return "", "", errors.New("This session was disabled")
	}
	if err != nil {
		return "", "", err
	}
	if currentSession.Disabled {
		return "", "", errors.New("disable session")
	}
	if createHashSHA256WithSalt(rfToken) == currentSession.RfToken {
		accessToken, _ := a.CreateToken(payLoad.Username, payLoad.Admin, payLoad.IdSession)
		newRfToken := a.createRandomOutput()
		currentSession.RfToken = createHashSHA256WithSalt(newRfToken)
		a.AuthRep.Update(ctx.Request().Context(), currentSession)
		return accessToken, newRfToken, nil
	}
	return "", "", errors.New("disable session")
}

func (a AuthService) createRandomOutput(sal ...string) string {
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

func (a AuthService) DisableSession(ctx echo.Context, id uuid.UUID) error {
	err := a.AuthRep.Disable(ctx.Request().Context(), id)
	return err
}

func (a AuthService) GetUser(ctx echo.Context) (*models.User, error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*CustomClaims)
	User, err := a.UserRep.Get(ctx.Request().Context(), claims.Username)
	return User, err
}
