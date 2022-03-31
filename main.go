package main

import (
	"app/internal/config"
	"app/internal/handlers"
	repository "app/internal/repository/mongodb"
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

	// connPool, err := pgxpool.Connect(context.Background(), configuration.UrlPostgres())
	// if err != nil {
	// 	log.Println("Connecting url", configuration.UrlPostgres())
	// 	log.Fatal(err)
	// }

	// repoPg := repository.NewRepoEntityPostgres(*connPool)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// entService := service.NewEntityService(&repoPg)

	clientMongo, err := mongo.Connect(context.Background(), options.Client().ApplyURI(configuration.ConnectUrlMongo()))
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}   ${uri}  ${status}    ${latency_human}\n",
	}))

	// Entity struct
	repoEntityMongo := repository.NewRepoEntityMongoDB(*clientMongo)
	entService := service.NewEntityService(&repoEntityMongo)
	handlerEntity := handlers.EntityHandler{EntityService: entService}

	entityGr := e.Group("/entity")
	entityGr.GET("", handlerEntity.List)
	entityGr.GET("/:id", handlerEntity.GetDetail)
	entityGr.PUT("/:id", handlerEntity.Update)
	entityGr.DELETE("/:id", handlerEntity.Delete)
	entityGr.POST("", handlerEntity.Create)

	// User struct
	userRepoMongo := repository.NewUserRepoMongoDB(*clientMongo)
	userService := service.NewUserService(userRepoMongo)
	userHandler := handlers.UserHandler{
		Ser: *userService,
	}

	userGr := e.Group("/user")
	userGr.GET("/:username", userHandler.Get)
	userGr.POST("", userHandler.Create)
	userGr.DELETE("", userHandler.Delete)
	userGr.GET("", userHandler.GetAll)

	// Auth struct
	authService := service.Auth{UserRep: userRepoMongo}
	authHandler := handlers.AuthHandler{
		AuthService: authService,
	}

	jwtConf := authService.JWTConfig()

	authGr := e.Group("/auth")

	authGr.POST("", authHandler.IsAuthentication)
	authGr.POST("/login", authHandler.Login)
	authGr.GET("/info", authHandler.Info, middleware.JWTWithConfig(jwtConf))
	authGr.GET("/logout", authHandler.Logout, middleware.JWTWithConfig(jwtConf))
	// Run Server
	e.Logger.Debug(e.Start(":8000"))
}
