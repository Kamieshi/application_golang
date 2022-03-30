package handlers

import (
	"app/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthService service.Auth
}

func (ah AuthHandler) IsAuthentication(c echo.Context) error {
	fmt.Print("as")
	var data usPss
	err := c.Bind(&data)
	if err != nil {
		return err
	}
	isAuth, err := ah.AuthService.IsAuthentication(c.Request().Context(), data.Username, data.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, isAuth)
}
