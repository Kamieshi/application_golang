package main

import (
	"app/internal/config"
	"app/internal/handlers"
	mongoRepository "app/internal/repository/mongodb"
	redisRepository "app/internal/repository/redis"
	"app/internal/service"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"

	"context"

	"github.com/labstack/echo/v4" //nolint:typecheck
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "app/docs/app"
)

// @title Golang Application Swagger
// @version 0.1
// @host localhost:8000

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configuration, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.InitLogger()

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

	e := echo.New() //nolint:typecheck
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}   ${uri}  ${status}    ${latency_human}\n",
	}))

	// Entity struct
	repoEntityMongo := mongoRepository.NewRepoEntityMongoDB(*clientMongo)
	repoCashEntity := redisRepository.NewCashSteamEntityRep(configuration.REDIS_URL)
	entService := service.NewEntityService(&repoEntityMongo, repoCashEntity)
	//entService := service.NewEntityService(&repoEntityMongo, nil)
	handlerEntity := handlers.EntityHandler{EntityService: entService}

	entityGr := e.Group("/entity")
	entityGr.GET("", handlerEntity.List)
	entityGr.GET("/:id", handlerEntity.GetDetail)
	entityGr.PUT("/:id", handlerEntity.Update)
	entityGr.DELETE("/:id", handlerEntity.Delete)
	entityGr.POST("", handlerEntity.Create)

	// User struct
	userRepoMongo := mongoRepository.NewUserRepoMongoDB(*clientMongo)
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
	authRep := mongoRepository.NewAuthRepoMongoDB(*clientMongo)

	authService := service.NewAuth(userRepoMongo, authRep)

	authHandler := handlers.AuthHandler{
		AuthService: authService,
	}

	e.Use(middleware.JWTWithConfig(authService.JWTConfig))

	authGr := e.Group("/auth")

	authGr.POST("/login", authHandler.Login)
	authGr.GET("/info", authHandler.Info)
	authGr.GET("/logout", authHandler.Logout)
	authGr.POST("/refresh", authHandler.Refresh)

	// Work with static
	e.Static("/images", "./static/images")

	// Work with upload images
	imageRepo := mongoRepository.NewImageRepoMongoDB(*clientMongo)
	imageService := service.ImageService{
		ImageRepository: imageRepo,
	}
	imageHandler := handlers.ImageHandler{ImageService: imageService}
	e.POST("/upload", imageHandler.Load)
	e.GET("/load/:easy_link", imageHandler.Get)

	//Swagger

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Run Server
	e.Logger.Debug(e.Start(":8000"))
}
