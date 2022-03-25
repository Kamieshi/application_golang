package views

import (
	"application/repository"
	"application/service/models"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ListEntities(c echo.Context) error {
	shard, err := repository.GetPg()

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}
	var entities []models.Entity = make([]models.Entity, 0)
	entities, err = repository.GetAllItems(shard)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}
	fmt.Println(entities)
	return c.String(http.StatusOK, "GetEntities")
}

func GetDetailEntitie(c echo.Context) error {

	return c.String(http.StatusOK, "GetDetailAntitie")
}

func UpdateEntitie(c echo.Context) error {

	return c.String(http.StatusOK, "GetDetailAntitie")
}

func CreateEntitie(c echo.Context) error {

	return c.String(http.StatusOK, "GetDetailAntitie")
}

func DeleteEntitie(c echo.Context) error {

	return c.String(http.StatusOK, "GetDetailAntitie")
}
