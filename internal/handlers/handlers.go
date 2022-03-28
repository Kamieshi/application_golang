package handlers

import (
	"app/internal/service"
	"app/internal/service/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type EntityHandler struct {
	EntityService service.EntityService
}

func (eh EntityHandler) List(c echo.Context) error {
	entitys, err := eh.EntityService.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entitys)
}

func (eh EntityHandler) GetDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	entity, err := eh.EntityService.GetForID(c.Request().Context(), id)

	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, entity)
}

func (eh EntityHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	entity := models.Entity{}

	err = c.Bind(&entity)
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
		return err
	}
	err = eh.EntityService.Add(c.Request().Context(), entity)

	return c.JSON(http.StatusOK, "{status : 'Update successful'}")
}

func (eh EntityHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	err = eh.EntityService.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("{message: %v}", err))
	}
	return c.JSON(http.StatusOK, "{message : 'Delete successful'}")
}
