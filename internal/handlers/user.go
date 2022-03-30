package handlers

import (
	"app/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	ser service.UserService
}

func (uh UserHandler) Get(c echo.Context) error {
	username := c.Param("username")
	user, err := uh.ser.Get(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusAccepted, user)
}
