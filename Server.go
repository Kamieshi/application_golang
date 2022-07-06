package main

import (
	_ "app/docs/app"
	"app/internal/config"
	"app/internal/handlers"
	repository "app/internal/repository"
	redisRepository "app/internal/repository/redis"
	"app/internal/service"
	"github.com/labstack/echo/v4" //nolint:typecheck
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
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
	config.InitLogger()

	//Cash repo
	repoCashEntity := redisRepository.NewCashSteamEntityRep(configuration.REDIS_URL)

	repoFactory := repository.GetFactory("pg", configuration)

	//Creating repositories Postgres
	repoEntity := repoFactory.GetEntityRepo()
	if err != nil {
		log.Fatal(err)
	}

	repoUsersPg := repoFactory.GetUserRepo()
	if err != nil {
		log.Fatal(err)
	}

	repoImagesPg := repoFactory.GetImageRepo()
	if err != nil {
		log.Fatal(err)
	}

	repoAuthPg := repoFactory.GetAuthRepo()
	if err != nil {
		log.Fatal(err)
	}

	//Creating services Postgres
	entService := service.NewEntityService(repoEntity, repoCashEntity)
	userService := service.NewUserService(repoUsersPg)
	imageService := service.NewImageService(repoImagesPg)
	authService := service.NewAuthService(repoUsersPg, repoAuthPg)

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

	// Run Server
	e.Logger.Debug(e.Start(":8000"))
}
