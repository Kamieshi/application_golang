package handlers

import (
	"app/internal/models"
	"app/internal/service"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EntityHandler struct {
	EntityService service.EntityService
}

// List godoc
// @tags Entity
// @Security ApiKeyAuth
// @Success      200  {array}  models.Entity
// @Router /entity [get]
func (eh EntityHandler) List(c echo.Context) error {
	entities, err := eh.EntityService.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, entities)
}

// GetDetail godoc
// @tags Entity
// @Security ApiKeyAuth
// @success 200 {object} models.Entity "desc"
// @Router /entity/{id} [get]
// @Param id path string false "Id of the mongo"
func (eh EntityHandler) GetDetail(c echo.Context) error {
	id := c.Param("id")

	entity, err := eh.EntityService.GetForID(c.Request().Context(), id)

	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, entity)
}

// Update godoc
// @tags Entity
// @Security ApiKeyAuth
// @Success      200  {string} {status : 'Update successful'}
// @Param id path string false "Id of the mongo"
// @Router /entity/{id} [put]
func (eh EntityHandler) Update(c echo.Context) error {
	id := c.Param("id")

	entity := models.Entity{}

	err := c.Bind(&entity)
	if err != nil {
		return err
	}

	err = eh.EntityService.Update(c.Request().Context(), id, entity)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))

	}
	return c.JSON(http.StatusOK, "{status : 'Update successful'}")
}

// Create godoc
// @tags Entity
// @Security ApiKeyAuth
// @Param DataEntity body models.Entity true "Entity model"
// @Router /entity [post]
func (eh EntityHandler) Create(c echo.Context) error {
	entity := models.Entity{}

	err := c.Bind(&entity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("{message: %v}", err))
	}

	err = eh.EntityService.Add(c.Request().Context(), &entity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("{message: %v}", err))
	}

	return c.JSON(http.StatusOK, entity)
}

// Delete godoc
// @tags Entity
// @Security ApiKeyAuth
// @Param id path string false "Id of the mongo"
// @Router /entity/{id} [delete]
func (eh EntityHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := eh.EntityService.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}

	return c.JSON(http.StatusOK, "{message : 'Delete successful'}")
}
