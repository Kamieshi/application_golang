package handlers

import (
	"app/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Ser *service.UserService
}

type usPss struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (uh UserHandler) Get(c echo.Context) error {
	username := c.Param("username")
	user, err := uh.Ser.Get(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusAccepted, user)
}

// Create godoc
// @tags User Control
// @Param userData body usPss true "Create new User"
// @Router /user [post]
func (uh UserHandler) Create(c echo.Context) error {
	var dat usPss
	err := c.Bind(&dat)
	if err != nil {
		return err
	}
	user, err := uh.Ser.Create(c.Request().Context(), dat.Username, dat.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusAccepted, user)
}

func (uh UserHandler) Delete(c echo.Context) error {

	username := c.Param("username")
	err := uh.Ser.Delete(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, "{message : 'Delete successful'}")
}

func (uh UserHandler) GetAll(c echo.Context) error {
	users, err := uh.Ser.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}
