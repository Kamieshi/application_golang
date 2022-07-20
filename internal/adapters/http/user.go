package http

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"app/internal/service"

	"github.com/labstack/echo/v4"
)

// UserHandler  Handler for work with User service
type UserHandler struct {
	Ser *service.UserService
}

type usPss struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Get godoc
// get models.User
// @tags User Control
// @Summary Get entity
// @Description Get entity
// @Security ApiKeyAuth
// @Param username path string true "username"
// @Success 200 {object} models.User
// @Failure 404 {string} User not found
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /user/{username} [get]
func (uh UserHandler) Get(c echo.Context) error {
	username := c.Param("username")
	user, err := uh.Ser.Get(c.Request().Context(), username)
	if err != nil {
		log.WithError(err).Error()
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// Create godoc
// @tags User Control
// @Summary Create User
// @Description Create new user
// @Param userData body usPss true "Create new User"
// @Success 201 {object} models.User
// @Failure 404 {string} User not found
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /user [post]
func (uh UserHandler) Create(c echo.Context) error {
	var dat usPss
	err := c.Bind(&dat)
	if err != nil {
		log.WithError(err).Error()
		return err
	}
	user, err := uh.Ser.Create(c.Request().Context(), dat.Username, dat.Password)
	if err != nil {
		log.WithError(err).Error()
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusCreated, user)
}

// Delete godoc
// Delete user
// @tags User Control
// @Summary Delete User
// @Description Delete new user
// @Security ApiKeyAuth
// @Param username path string true "username for delete"
// @Success 201 {object} models.User
// @Failure 404 {string} User not found
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /user/{username} [delete]
func (uh UserHandler) Delete(c echo.Context) error {
	username := c.Param("username")
	err := uh.Ser.Delete(c.Request().Context(), username)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// GetAll godoc
// @tags User Control
// @Summary Get all users
// @Description Get all users
// @Security ApiKeyAuth
// @Success 200 {array} models.User
// @Failure 500 {string} Error from service
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /user [get]
func (uh UserHandler) GetAll(c echo.Context) error {
	users, err := uh.Ser.GetAll(c.Request().Context())
	if err != nil {
		log.WithError(err).Error()
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}
