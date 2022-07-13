package adapters

import (
	"app/internal/adapters"
	"app/internal/adapters/http/handlers"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"sync"
)

type EchoHTTP struct {
	appServices  *adapters.ServicesApp
	echoInstance *echo.Echo
}

func (e *EchoHTTP) Start(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	err := e.echoInstance.Start(":8005")
	return err
}

func NewEchoHTTP(services *adapters.ServicesApp) *EchoHTTP {
	// Creating handlers
	e := &EchoHTTP{
		echoInstance: echo.New(),
		appServices:  services,
	}
	handlerEntity := handlers.EntityHandler{EntityService: e.appServices.EntityService}
	userHandler := handlers.UserHandler{Ser: e.appServices.UserService}
	authHandler := handlers.AuthHandler{AuthService: e.appServices.AuthService}
	imageHandler := handlers.ImageHandler{ImageService: e.appServices.ImageService}

	e.echoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}   ${uri}  ${status}    ${latency_human}\n",
	}))

	// Entity Routing
	entityGr := e.echoInstance.Group("/entity")
	entityGr.GET("", handlerEntity.List)
	entityGr.GET("/:id", handlerEntity.GetDetail)
	entityGr.PUT("/:id", handlerEntity.Update)
	entityGr.DELETE("/:id", handlerEntity.Delete)
	entityGr.POST("", handlerEntity.Create)

	// User Routing
	userGr := e.echoInstance.Group("/user")
	userGr.GET("/:username", userHandler.Get)
	userGr.POST("", userHandler.Create)
	userGr.DELETE("", userHandler.Delete)
	userGr.GET("", userHandler.GetAll)

	// Auth Routing
	e.echoInstance.Use(middleware.JWTWithConfig(*services.AuthService.JWTConfig))
	authGr := e.echoInstance.Group("/auth")
	authGr.POST("/login", authHandler.Login)
	authGr.GET("/info", authHandler.Info)
	authGr.GET("/logout", authHandler.Logout)
	authGr.POST("/refresh", authHandler.Refresh)

	// static
	e.echoInstance.Static("/images", "./static/images")

	// Image Routing
	e.echoInstance.POST("/upload", imageHandler.Load)
	e.echoInstance.GET("/load/:easy_link", imageHandler.Get)

	//Swagger
	e.echoInstance.GET("/swagger/*", echoSwagger.WrapHandler)

	//Ping
	e.echoInstance.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

	return e
}
