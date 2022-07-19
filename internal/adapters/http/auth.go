// Package http Work with http adapter from Echo
package http

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"app/internal/service"
)

// AuthHandler Handler for work with AuthService
type AuthHandler struct {
	AuthService *service.AuthService
}

// Login godoc
// login user
// @Summary Login User
// @Description Login user
// @Tags Auth
// @Param formLogin body usPss true "Login form"
// @Success 200 {object} responseSuccessLogin
// @Failure 401 {string} Invalid username or password
// @Failure 502 {string} Error create token
// @Router /auth/login [post]
func (a *AuthHandler) Login(c echo.Context) error {
	var data usPss
	err := c.Bind(&data)
	if err != nil {
		logrus.WithError(err).Error()
		return err
	}
	if valError := c.Validate(data); valError != nil {
		return valError
	}
	user, err := a.AuthService.IsAuthentication(c.Request().Context(), data.Username, data.Password)
	if err != nil {
		logrus.WithError(err).Error()
		return c.String(http.StatusUnauthorized, "invalid Username or password")
	}
	if user == nil {
		return echo.ErrUnauthorized
	}
	session, err := a.AuthService.CreateAndWriteSession(c, *user)
	if err != nil {
		logrus.WithError(err).Error()
		return echo.ErrUnauthorized
	}
	token, err := a.AuthService.CreateToken(user.UserName, user.Admin, session.ID)
	if err != nil {
		logrus.WithError(err).Error()
		return c.String(http.StatusBadGateway, err.Error())
	}

	return c.JSON(http.StatusOK, responseSuccessLogin{
		Access:  token,
		Refresh: session.RfToken,
	})
}

// Info godoc
// Get info about current login user
// @tags Auth
// @Summary Get info about current user
// @Description Info about current user
// @Security ApiKeyAuth
// @Success 202 {object} models.User
// @Failure 401 {string} User unauthorized
// @Router /auth/info [get]
func (a *AuthHandler) Info(c echo.Context) error {
	user, _ := a.AuthService.GetUser(c)
	if user != nil {
		return c.JSON(http.StatusAccepted, user)
	}
	return echo.ErrUnauthorized
}

// Logout godoc
// @Summary Logout route
// @Description Logout current active user
// @tags Auth
// @Security ApiKeyAuth
// @Success 202 {string} Logout complete successful
// @Failure 400 {string} User unauthorized
// @Router /auth/logout [get]
func (a *AuthHandler) Logout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.CustomClaims)
	err := a.AuthService.DisableSession(c, claims.IDSession)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)
	return c.String(http.StatusAccepted, "Logout")
}

type rft struct {
	Refresh string `json:"refresh" `
}

// Refresh godoc
// refresh current access token
// @tags Auth
// @Summary Refresh session
// @Description Refresh session for current user
// @Security ApiKeyAuth
// @Param refreshData body rft true "refresh"
// @Access 202 {object} responseSuccessLogin
// @Failure 401 {string} Uncorrected tokens
// @Failure 400 {string} Error parsing input values
// @Router /auth/refresh [get]
func (a *AuthHandler) Refresh(c echo.Context) error {
	var rt rft
	err := c.Bind(&rt)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	acToken, refToken, err := a.AuthService.RefreshAndWriteSession(c, rt.Refresh)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusAccepted, responseSuccessLogin{
		Access:  acToken,
		Refresh: refToken,
	})
}
