package handlers

import (
	"app/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type EntityHandler struct {
	ContrEntity repository.ControllerEntity
}

func (eh *EntityHandler) ListEntities(c echo.Context) error {

	return c.String(http.StatusOK, "GetEntities")
}

func (eh *EntityHandler) GetDetailEntitie(c echo.Context) error {
	pId := c.Param("id")
	id, err := strconv.Atoi(pId)
	if err != nil {
		return err
	}
	entity, err := eh.ContrEntity.GetItemForID(c.Request().Context(), id)

	if err != nil {
		return c.String(http.StatusNotFound, "Not Found")
	}

	marshaldata, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(marshaldata))
}

func (eh *EntityHandler) UpdateEntitie(c echo.Context) error {

	return c.String(http.StatusOK, "GetDetailAntitie")
}

func (eh *EntityHandler) CreateEntitie(c echo.Context) error {
	return c.String(http.StatusOK, "GetDetailAntitie")
}

func (eh *EntityHandler) DeleteEntitie(c echo.Context) error {

	return c.String(http.StatusOK, "GetDetailAntitie")
}
