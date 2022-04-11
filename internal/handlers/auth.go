package handlers

import (
	"app/internal/models"
	"app/internal/service"
	"fmt"
	echo "github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AuthHandler struct {
	AuthService service.Auth
}

func (ah AuthHandler) IsAuthentication(c echo.Context) error {
	var data usPss
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	_, isAuth, err := ah.AuthService.IsAuthentication(c.Request().Context(), data.Username, data.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, isAuth)
}

// Login godoc
// @tags Auth
// @Param formLogin body handlers.usPss true "Login form"
// @Router /auth/login [post]
func (ah AuthHandler) Login(c echo.Context) error {
	var data usPss
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	user, isAuth, err := ah.AuthService.IsAuthentication(c.Request().Context(), data.Username, data.Password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "invalid Username or password")
	}
	if !isAuth {
		return echo.ErrUnauthorized
	}
	session, err := ah.AuthService.CreateAndWriteSession(c, user)
	if err != nil {
		return echo.ErrUnauthorized
	}
	token, err := ah.AuthService.CreateToken(user.UserName, user.Admin, session.IdSession)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access":  token,
		"refresh": session.RfToken,
	})
}

func (ah AuthHandler) Info(c echo.Context) error {
	user, _ := ah.AuthService.GetUser(c)
	if (user != models.User{}) {
		return c.JSON(http.StatusAccepted, user)
	}
	return echo.ErrUnauthorized
}

// Logout godoc
// @tags Auth
// @Security ApiKeyAuth
// @Router /auth/logout [get]
func (ah AuthHandler) Logout(c echo.Context) error {
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
func (ah AuthHandler) Refresh(c echo.Context) error {
	var rt rft
	err := c.Bind(&rt)
	if err != nil {
		return err
	}
	acToken, refToken, err := ah.AuthService.RefreshAndWriteSession(c, rt.Refresh)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"access":  acToken,
		"refresh": refToken,
	})
}
