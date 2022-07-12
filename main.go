package main

import (
	_ "app/docs/app"
	"app/internal/config"
	"app/internal/handlers"
	"app/internal/repository"
	repositoryMongoDB "app/internal/repository/mongodb"
	repositoryPg "app/internal/repository/posgres"
	redisRepository "app/internal/repository/redis"
	"app/internal/service"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4" //nolint:typecheck
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"time"
)

// @title Golang Application Swagger
// @version 0.1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configuration, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})

	//Create repository
	var repoEntity repository.RepoEntity
	var repoUsers repository.RepoUser
	var repoImages repository.RepoImage
	var repoAuth repository.RepoSession

	//Init repository
	typeDB := "pg"
	if typeDB == "pg" {
		connPool, err := pgxpool.Connect(context.Background(), configuration.UrlPostgres())
		if err != nil {
			log.Println("Connecting url", configuration.UrlPostgres())
			log.Fatal(err)
		}
		repoEntity = repositoryPg.NewRepoEntityPostgres(connPool)
		repoAuth = repositoryPg.NewRepoAuthPostgres(connPool)
		repoUsers = repositoryPg.NewRepoUsersPostgres(connPool)
		repoImages = repositoryPg.NewRepoImagePostgres(connPool)
	} else {
		timeOutConnect, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		clientMongo, err := mongo.Connect(timeOutConnect, options.Client().ApplyURI(configuration.ConnectUrlMongo()))
		defer cancel()
		if err != nil {
			log.WithError(err).Panic("Error with mongo connection")
		}

		var rf *readpref.ReadPref
		err = clientMongo.Ping(timeOutConnect, rf)
		if err != nil {
			log.WithError(err).Panic("Error with mongo connection")
		}
		repoEntity = repositoryMongoDB.NewRepoEntityMongoDB(*clientMongo)
		repoAuth = repositoryMongoDB.NewAuthRepoMongoDB(*clientMongo)
		repoUsers = repositoryMongoDB.NewUserRepoMongoDB(*clientMongo)
		repoImages = repositoryMongoDB.NewImageRepoMongoDB(*clientMongo)
	}

	//Cash repo
	repoCashEntity := redisRepository.NewCashSteamEntityRep(configuration.REDIS_URL, &repoEntity)
	if err != nil {
		log.Fatal(err)
	}

	go repoCashEntity.Listener(context.Background())

	//Creating services Postgres
	entService := service.NewEntityService(repoEntity, repoCashEntity)
	userService := service.NewUserService(repoUsers)
	imageService := service.NewImageService(repoImages)
	authService := service.NewAuthService(repoUsers, repoAuth)

	// Creating handlers
	handlerEntity := handlers.EntityHandler{EntityService: entService}
	userHandler := handlers.UserHandler{Ser: userService}
	authHandler := handlers.AuthHandler{AuthService: authService}
	imageHandler := handlers.ImageHandler{ImageService: imageService}

	e := echo.New() //nolint:typecheck
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}   ${uri}  ${status}    ${latency_human}\n",
	}))

	// Entity Routing
	entityGr := e.Group("/entity")
	entityGr.GET("", handlerEntity.List)
	entityGr.GET("/:id", handlerEntity.GetDetail)
	entityGr.PUT("/:id", handlerEntity.Update)
	entityGr.DELETE("/:id", handlerEntity.Delete)
	entityGr.POST("", handlerEntity.Create)

	// User Routing
	userGr := e.Group("/user")
	userGr.GET("/:username", userHandler.Get)
	userGr.POST("", userHandler.Create)
	userGr.DELETE("", userHandler.Delete)
	userGr.GET("", userHandler.GetAll)

	// Auth Routing
	e.Use(middleware.JWTWithConfig(authService.JWTConfig))
	authGr := e.Group("/auth")
	authGr.POST("/login", authHandler.Login)
	authGr.GET("/info", authHandler.Info)
	authGr.GET("/logout", authHandler.Logout)
	authGr.POST("/refresh", authHandler.Refresh)

	// static
	e.Static("/images", "./static/images")

	// Image Routing
	e.POST("/upload", imageHandler.Load)
	e.GET("/load/:easy_link", imageHandler.Get)

	//Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//Ping
	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

	// Run Server
	e.Logger.Debug(e.Start(":8005"))
}
