package handlers

import (
	"app/internal/models"
	"app/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EntityHandler struct {
	EntityService service.EntityService
}

func (eh EntityHandler) List(c echo.Context) error {
	entities, err := eh.EntityService.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entities)
}

func (eh EntityHandler) GetDetail(c echo.Context) error {
	id := c.Param("id")

	entity, err := eh.EntityService.GetForID(c.Request().Context(), id)

	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, entity)
}

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

func (eh EntityHandler) Create(c echo.Context) error {
	entity := models.Entity{}

	err := c.Bind(&entity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("{message: %v}", err))
	}

	err = eh.EntityService.Add(c.Request().Context(), entity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("{message: %v}", err))
	}

	return c.JSON(http.StatusOK, "{status : 'Update successful'}")
}

func (eh EntityHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := eh.EntityService.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}

	return c.JSON(http.StatusOK, "{message : 'Delete successful'}")
}
