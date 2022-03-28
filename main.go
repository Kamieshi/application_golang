package main

import (
	"app/internal/config"
	"app/internal/handlers"
	"app/internal/repository"
	"app/internal/service"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	configuration, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	connPool, err := pgxpool.Connect(context.Background(), configuration.UrlPosgres())
	if err != nil {
		log.Println("Connecting url", configuration.UrlPosgres())
		log.Fatal(err)
	}

	repo := repository.RepoEntityPostgres{Pool: connPool}
	if err != nil {
		log.Fatal(err)
	}

	entService := service.EntityService{Rep: repo}

	handlerEntity := handlers.EntityHandler{EntityService: entService}

	e := echo.New()
	e.GET("/entity", handlerEntity.List)
	e.GET("/entity/:id", handlerEntity.GetDetail)
	e.PUT("/entity/:id", handlerEntity.Update)
	e.DELETE("/entity/:id", handlerEntity.Delete)
	e.POST("/entity", handlerEntity.Create)

	e.Logger.Debug(e.Start(":8000"))
}
