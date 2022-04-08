package handlers

import (
	"app/internal/models"
	"app/internal/service"
	"fmt"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
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
