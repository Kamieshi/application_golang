package handlers

import (
	"application/handlers/views"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/entity", views.ListEntities)
	e.GET("/entity/:id", views.GetDetailEntitie)
	e.PUT("/entity/:id", views.UpdateEntitie)
	e.DELETE("/entity/:id", views.DeleteEntitie)
	e.POST("/entity", views.CreateEntitie)
}
