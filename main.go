package main

import (
	"app/internal/config"
	"app/internal/handlers"
	"app/internal/repository"
	"app/internal/service"
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	configuration, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// connPool, err := pgxpool.Connect(context.Background(), configuration.UrlPosgres())
	// if err != nil {
	// 	log.Println("Connecting url", configuration.UrlPosgres())
	// 	log.Fatal(err)
	// }

	// repoPg := repository.NewRepoEntityPostgres(*connPool)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// entService := service.NewEntityService(&repoPg)

	clientMongo, err := mongo.Connect(context.Background(), options.Client().ApplyURI(configuration.CoonnectUrlMongo()))
	if err != nil {
		log.Fatalln(err)
	}

	repoMongo := repository.NewRepoEntityMongoDB(*clientMongo)

	entService := service.NewEntityService(&repoMongo)

	handlerEntity := handlers.EntityHandler{EntityService: entService}

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}   ${uri}  ${status}    ${latency_human}\n",
	}))

	e.GET("/entity", handlerEntity.List)
	e.GET("/entity/:id", handlerEntity.GetDetail)
	e.PUT("/entity/:id", handlerEntity.Update)
	e.DELETE("/entity/:id", handlerEntity.Delete)
	e.POST("/entity", handlerEntity.Create)

	e.Logger.Debug(e.Start(":8000"))
}
