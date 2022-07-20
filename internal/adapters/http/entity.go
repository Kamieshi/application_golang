package http

import (
	"fmt"
	"net/http"

	ech "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"app/internal/models"
	"app/internal/service"
)

// EntityHandler Handler for work with Entity service
type EntityHandler struct {
	EntityService *service.EntityService
}

// List godoc
// Get all entities from DB
// @tags Entity
// @Summary Get list Entities
// @Description Get all entities from DB
// @Success      200  {array}  models.Entity
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /entity [get]
func (eh EntityHandler) List(c ech.Context) error {
	entities, err := eh.EntityService.GetAll(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Error(c.Request())
		return err
	}
	return c.JSON(http.StatusOK, entities)
}

// GetDetail godoc
// Get object models.Entity
// @Security ApiKeyAuth
// @Summary Get detail about entity
// @Description Get detail about entity
// @tags Entity
// @Param id path string false "UU ID string"
// @Success 200 {object} models.Entity "Success"
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /entity/{id} [get]
func (eh EntityHandler) GetDetail(c ech.Context) error {
	id := c.Param("id")

	entity, err := eh.EntityService.GetForID(c.Request().Context(), id)

	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, entity)
}

// Update godoc
// update models.Entity
// @tags Entity
// @Summary Update entity
// @Description Update entity
// @Security ApiKeyAuth
// @Param id path string false "UU ID string"
// @Param DataEntity body models.Entity true "Entity model"
// @Success 201 {string} Update successful
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /entity/{id} [put]
func (eh EntityHandler) Update(c ech.Context) error {
	id := c.Param("id")
	entity := models.Entity{}
	err := c.Bind(&entity)
	if err != nil {
		return err
	}
	if valErr := c.Validate(entity); valErr != nil {
		return valErr
	}
	err = eh.EntityService.Update(c.Request().Context(), id, &entity)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.String(http.StatusCreated, "Update successful")
}

// Create godoc
// create models.Entity
// Create new Entity
// @tags Entity
// @Summary Create entity
// @Description Create entity
// @Security ApiKeyAuth
// @Param DataEntity body models.Entity true "Entity model"
// @Success 201 {object} models.Entity
// @Failure 500 {string} Error Parse input data
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /entity [post]
func (eh EntityHandler) Create(c ech.Context) error {
	entity := models.Entity{}

	err := c.Bind(&entity)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if valErr := c.Validate(entity); valErr != nil {
		return valErr
	}
	err = eh.EntityService.Add(c.Request().Context(), &entity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("{message: %v}", err))
	}

	return c.JSON(http.StatusCreated, entity)
}

// Delete godoc
// delete models.Entity
// @tags Entity
// @Summary Delete entity
// @Description Delete entity
// @Security ApiKeyAuth
// @Param id path string false "ID of the mongo"
// @Success 204 {} Delete successful
// @Failure 500 {string} DeleteError
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /entity/{id} [delete]
func (eh EntityHandler) Delete(c ech.Context) error {
	id := c.Param("id")

	err := eh.EntityService.Delete(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
