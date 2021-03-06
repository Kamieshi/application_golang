// Package service Work with services
package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"

	"app/internal/models"
	"app/internal/repository"

	"github.com/golang-jwt/jwt"
	ech "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const maxRand = 10000

// AuthService service for work with Auth end Session
type AuthService struct {
	UserRep   repository.RepoUser
	AuthRep   repository.RepoSession
	JWTConfig *middleware.JWTConfig
}

// NewAuthService Constructor
func NewAuthService(userRep repository.RepoUser, sessionRep repository.RepoSession) *AuthService {
	headerAuthorization := ech.HeaderAuthorization
	config := &middleware.JWTConfig{
		Claims:        &CustomClaims{},
		SigningKey:    []byte(os.Getenv("SECRET_KEY")),
		AuthScheme:    "Bearer",
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   "header:" + headerAuthorization,
		ParseTokenFunc: func(auth string, c ech.Context) (interface{}, error) {
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
				if err.Error()[:16] == "token is expired" && c.Path() == "/auth/refresh" {
					return token, nil
				}
				return nil, fmt.Errorf("service auth/NewAuthService : %v", err)
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
		Skipper: func(c ech.Context) bool {
			switch c.Path() {
			case "/auth/login":
				return true
			case "/user":
				if c.Request().Method == http.MethodPost {
					return true
				}
			case "/swagger/*":
				return true
			case "/ping":
				return true
			}
			return false
		},
	}

	return &AuthService{
		UserRep:   userRep,
		AuthRep:   sessionRep,
		JWTConfig: config,
	}
}

// CustomClaims Payload from access token
type CustomClaims struct {
	Username  string    `json:"name"`
	Admin     bool      `json:"admin"`
	IDSession uuid.UUID `json:"id_session"`
	jwt.StandardClaims
}

// IsAuthentication Check JWT of session and
func (a *AuthService) IsAuthentication(ctx context.Context, username, password string) (*models.User, error) {
	user, err := a.UserRep.Get(ctx, username)

	if err != nil {
		return user, fmt.Errorf("service auth/IsAuthentication : %v", err)
	}
	inPasswordHash := createHash256Password(user, password)
	if user.PasswordHash == inPasswordHash {
		return user, nil
	}
	return nil, fmt.Errorf("service auth/IsAuthentication : %v", err)
}

// CreateToken Create access token
func (a *AuthService) CreateToken(username string, admin bool, idSession uuid.UUID) (string, error) {
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
		return "", fmt.Errorf("service auth/CreateToken : %v", err)
	}
	fmt.Println(token)
	return tt, nil
}

// CreateAndWriteSession Create new session and write into repository
func (a *AuthService) CreateAndWriteSession(ctx ech.Context, user models.User) (models.Session, error) {
	refresh := a.createRandomOutput(user.UserName)

	session := models.Session{
		ID:              uuid.New(),
		UserID:          user.ID,
		CreatedAt:       time.Now(),
		Disabled:        false,
		RfToken:         createHashSHA256WithSalt(refresh),
		UniqueSignature: ctx.Request().UserAgent(),
	}
	err := a.AuthRep.Create(ctx.Request().Context(), &session)
	if err != nil {
		return models.Session{}, fmt.Errorf("service auth/CreateAndWriteSession : %v", err)
	}
	return session, err
}

// RefreshAndWriteSession Update RF token
func (a *AuthService) RefreshAndWriteSession(ctx ech.Context, rfToken string) (AccessToken, RfToken string, err error) {
	user := ctx.Get("user").(*jwt.Token)
	payLoad := user.Claims.(*CustomClaims)
	currentSession, err := a.AuthRep.Get(ctx.Request().Context(), payLoad.IDSession)
	if currentSession.Disabled {
		return "", "", fmt.Errorf("service auth/RefreshAndWriteSession : %v", errors.New("this session was disabled"))
	}
	if err != nil {
		return "", "", fmt.Errorf("service auth/RefreshAndWriteSession : %v", err)
	}
	if rfToken == currentSession.RfToken {
		AccessToken, _ = a.CreateToken(payLoad.Username, payLoad.Admin, payLoad.IDSession)
		RfToken = createHashSHA256WithSalt(a.createRandomOutput())
		currentSession.RfToken = RfToken
		err = a.AuthRep.Update(ctx.Request().Context(), currentSession)
		if err != nil {
			return "", "", fmt.Errorf("service auth/RefreshAndWriteSession : %v", err)
		}
		return
	}
	return "", "", fmt.Errorf("service auth/RefreshAndWriteSession : %v", errors.New("disable session"))
}

func (a *AuthService) createRandomOutput(sal ...string) string {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(maxRand))
	data := fmt.Sprint(time.Now().Unix(), os.Getenv("SECRET_KEY"), nBig.Int64(), sal)
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

// DisableSession Disable session (disable false -> true in repository)
func (a *AuthService) DisableSession(ctx ech.Context, id uuid.UUID) error {
	err := a.AuthRep.Disable(ctx.Request().Context(), id)
	if err != nil {
		return fmt.Errorf("service auth/DisableSession : %v", err)
	}
	return err
}

// GetUser Check access token and return user
func (a *AuthService) GetUser(ctx ech.Context) (*models.User, error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*CustomClaims)
	User, err := a.UserRep.Get(ctx.Request().Context(), claims.Username)
	if err != nil {
		return nil, fmt.Errorf("service auth/GetUser : %v", err)
	}
	return User, err
}
