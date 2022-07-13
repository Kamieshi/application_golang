package main

import (
	_ "app/docs/app"
	"app/internal/adapters/http/handlers"
	"app/internal/config"
	"app/internal/repository"
	repositoryMongoDB "app/internal/repository/mongodb"
	repositoryPg "app/internal/repository/posgres"
	redisRepository "app/internal/repository/redis"
	"app/internal/service"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"sync"

	//nolint:typecheck
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
		repoEntity = repositoryMongoDB.NewRepoEntityMongoDB(clientMongo)
		repoAuth = repositoryMongoDB.NewAuthRepoMongoDB(clientMongo)
		repoUsers = repositoryMongoDB.NewUserRepoMongoDB(clientMongo)
		repoImages = repositoryMongoDB.NewImageRepoMongoDB(clientMongo)
	}

	//Cash repo
	repoCashEntity := redisRepository.NewCashSteamEntityRep(configuration.REDIS_URL, &repoEntity)
	if err != nil {
		log.Fatal(err)
	}

	go repoCashEntity.Listener(context.Background())

	//Creating services Postgres

	AuthService := service.NewAuthService(&repoUsers, &repoAuth)
	EntityService := service.NewEntityService(&repoEntity, repoCashEntity)
	ImageService := service.NewImageService(&repoImages)
	UserService := service.NewUserService(&repoUsers)

	// Creating adapters
	// Echo HTTP
	e := echo.New()
	handlerEntity := handlers.EntityHandler{EntityService: EntityService}
	userHandler := handlers.UserHandler{Ser: UserService}
	authHandler := handlers.AuthHandler{AuthService: AuthService}
	imageHandler := handlers.ImageHandler{ImageService: ImageService}

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
	e.Use(middleware.JWTWithConfig(*AuthService.JWTConfig))
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
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info(e.Start(":8005"))
	}()

	wg.Wait()
}
