package handlers

import (
	"app/internal/service"
	"app/internal/service/models"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

func (ah AuthHandler) Login(c echo.Context) error {
	var data usPss
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	user, isAuth, err := ah.AuthService.IsAuthentication(c.Request().Context(), data.Username, data.Password)
	if !isAuth {
		return echo.ErrUnauthorized
	}
	token, err := ah.AuthService.CreateToken(user)
	if err != nil {
		return echo.ErrUnauthorized
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (ah AuthHandler) Info(c echo.Context) error {
	user, _ := ah.AuthService.GetUser(c)
	if (user != models.User{}) {
		return c.JSON(http.StatusAccepted, user)
	}
	return echo.ErrUnauthorized
}

func (ah AuthHandler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)
	return c.String(http.StatusAccepted, "Logout")
}

func (ah AuthHandler) refrash(c echo.Context) error {
	return nil
}
