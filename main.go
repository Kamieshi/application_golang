package main

import (
	ctx "context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"

	_ "app/docs/app"
	adaptergRPC "app/internal/adapters/grpc"
	adapterHTTP "app/internal/adapters/http"
	"app/internal/config"
	"app/internal/repository"
	repositoryMongoDB "app/internal/repository/mongodb"
	repositoryPg "app/internal/repository/posgres"
	redisRepository "app/internal/repository/redis"
	"app/internal/service"
)

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
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

	formatter := runtime.Formatter{ChildFormatter: &log.TextFormatter{
		FullTimestamp: true,
	}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Create repository
	var repoEntity repository.RepoEntity
	var repoUsers repository.RepoUser
	var repoImages repository.RepoImage
	var repoAuth repository.RepoSession

	// Init repository
	if configuration.UsedDB == "pg" {
		var connPool *pgxpool.Pool
		connPool, err = pgxpool.Connect(ctx.Background(), configuration.ConnectingURLPostgres())
		if err != nil {
			log.Println("Connecting url", configuration.ConnectingURLPostgres())
			log.Fatal(err)
		}
		repoEntity = repositoryPg.NewRepoEntityPostgres(connPool)
		repoAuth = repositoryPg.NewRepoAuthPostgres(connPool)
		repoUsers = repositoryPg.NewRepoUsersPostgres(connPool)
		repoImages = repositoryPg.NewRepoImagePostgres(connPool)
	} else {
		timeOutConnect, cancel := ctx.WithTimeout(ctx.Background(), time.Duration(configuration.TTLBackground)*time.Second)
		var clientMongo *mongo.Client
		clientMongo, err = mongo.Connect(timeOutConnect, options.Client().ApplyURI(configuration.ConnectingURLMongo()))
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

	// Cash repo
	repoCashEntity := redisRepository.NewCashSteamEntityRep(configuration.RedisURL, repoEntity)
	if err != nil {
		log.Error(err)
	}

	go repoCashEntity.Listener(ctx.Background())

	// Creating services Postgres

	AuthService := service.NewAuthService(repoUsers, repoAuth)
	EntityService := service.NewEntityService(repoEntity, repoCashEntity)
	ImageService := service.NewImageService(repoImages)
	UserService := service.NewUserService(repoUsers)

	// Creating adapters
	// Echo HTTP
	e := echo.New()
	e.Validator = adapterHTTP.NewCustomValidator()
	handlerEntity := adapterHTTP.EntityHandler{EntityService: EntityService}
	userHandler := adapterHTTP.UserHandler{Ser: UserService}
	authHandler := adapterHTTP.AuthHandler{AuthService: AuthService}
	imageHandler := adapterHTTP.ImageHandler{ImageService: ImageService}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}   ${uri}  ${status}    ${latency_human}\n",
	}))
	entityGr := e.Group("/entity")
	entityGr.GET("", handlerEntity.List)
	entityGr.GET("/:id", handlerEntity.GetDetail)
	entityGr.PUT("/:id", handlerEntity.Update)
	entityGr.DELETE("/:id", handlerEntity.Delete)
	entityGr.POST("", handlerEntity.Create)
	userGr := e.Group("/user")
	userGr.GET("/:username", userHandler.Get)
	userGr.POST("", userHandler.Create)
	userGr.DELETE("", userHandler.Delete)
	userGr.GET("", userHandler.GetAll)
	e.Use(middleware.JWTWithConfig(*AuthService.JWTConfig))
	authGr := e.Group("/auth")
	authGr.POST("/login", authHandler.Login)
	authGr.GET("/info", authHandler.Info)
	authGr.GET("/logout", authHandler.Logout)
	authGr.POST("/refresh", authHandler.Refresh)
	e.Static("/images", "./static/images")
	e.POST("/upload", imageHandler.Load)
	e.GET("/load/:easy_link", imageHandler.Get)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

	// Creating gRPC adapter
	listener, err := net.Listen(configuration.GrpcProtocol, ":"+configuration.GrpcPort)
	if err != nil {
		log.Errorf("main.go/main Start Lisening protocol %s port %s : %v", configuration.GrpcProtocol, configuration.GrpcPort, err)
	}

	grpcServer := grpc.NewServer()

	serverUser := struct {
		adaptergRPC.UserServer
	}{}
	adaptergRPC.RegisterEntityServer(grpcServer, &adaptergRPC.EntityServerImplement{EntityServ: EntityService})
	adaptergRPC.RegisterUserServer(grpcServer, &serverUser)
	adaptergRPC.RegisterImageManagerServer(grpcServer, &adaptergRPC.ImageServerImplement{ImageService: ImageService})
	// Run Server
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("HTTP ECHO Start work")
		log.Info(e.Start(":" + configuration.EchoPort))
		log.Info("HTTP ECHO End work")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("gRPC server Start work")
		log.Info(grpcServer.Serve(listener))
		log.Info("gRPC End start work")
	}()
	wg.Wait()
}
