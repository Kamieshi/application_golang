package handlers

import (
	"app/internal/service"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// AuthHandler Handler for work with AuthService
type AuthHandler struct {
	AuthService service.AuthService
}

func (a *AuthHandler) IsAuthentication(c echo.Context) error {
	var data usPss
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	_, isAuth, err := a.AuthService.IsAuthentication(c.Request().Context(), data.Username, data.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, isAuth)
}

// Login godoc
// @tags Auth
// @Param formLogin body handlers.usPss true "Login form"
// @Router /auth/login [post]
func (a *AuthHandler) Login(c echo.Context) error {
	var data usPss
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	user, isAuth, err := a.AuthService.IsAuthentication(c.Request().Context(), data.Username, data.Password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "invalid Username or password")
	}
	if !isAuth {
		return echo.ErrUnauthorized
	}
	session, err := a.AuthService.CreateAndWriteSession(c, *user)
	if err != nil {
		return echo.ErrUnauthorized
	}
	token, err := a.AuthService.CreateToken(user.UserName, user.Admin, session.ID.String())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access":  token,
		"refresh": session.RfToken,
	})
}

// Info check auth
func (a *AuthHandler) Info(c echo.Context) error {
	user, _ := a.AuthService.GetUser(c)
	if user != nil {
		return c.JSON(http.StatusAccepted, user)
	}
	return echo.ErrUnauthorized
}

// Logout godoc
// @tags Auth
// @Security ApiKeyAuth
// @Router /auth/logout [get]
func (a *AuthHandler) Logout(c echo.Context) error {
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
// @tags Auth
// @Security ApiKeyAuth
// @Param refreshData body rft true "refresh"
// @Router /auth/refresh [get]
func (a *AuthHandler) Refresh(c echo.Context) error {
	var rt rft
	err := c.Bind(&rt)
	if err != nil {
		return err
	}
	acToken, refToken, err := a.AuthService.RefreshAndWriteSession(c, rt.Refresh)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"access":  acToken,
		"refresh": refToken,
	})
}
