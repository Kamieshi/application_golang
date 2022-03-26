package main

import (
	"app/internal/conf"
	"app/internal/handlers"
	"app/internal/repository"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	configuration, err := conf.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	repo, err := repository.NewRepoPostgres(configuration.UrlPosgres())
	if err != nil {
		log.Fatal(err)
	}
	handlerEntity := handlers.EntityHandler{ContrEntity: repo}
	e := echo.New()

	e.GET("/entity", handlerEntity.ListEntities)
	e.GET("/entity/:id", handlerEntity.GetDetailEntitie)
	e.PUT("/entity/:id", handlerEntity.UpdateEntitie)
	e.DELETE("/entity/:id", handlerEntity.DeleteEntitie)
	e.POST("/entity", handlerEntity.CreateEntitie)

	e.Logger.Debug(e.Start(":8000"))
}
