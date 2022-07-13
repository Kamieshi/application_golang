package main

import (
	_ "app/docs/app"
	adapters "app/internal/adapters/http"
	"app/internal/config"
	"app/internal/repository"
	repositoryMongoDB "app/internal/repository/mongodb"
	repositoryPg "app/internal/repository/posgres"
	redisRepository "app/internal/repository/redis"
	"app/internal/service"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
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

	services := &service.ServicesApp{
		AuthService:   service.NewAuthService(&repoUsers, &repoAuth),
		EntityService: service.NewEntityService(&repoEntity, repoCashEntity),
		ImageService:  service.NewImageService(&repoImages),
		UserService:   service.NewUserService(&repoUsers),
	}
	// Creating adapters
	httpAdapter := adapters.NewEchoHTTP(services)

	// Run Server
	var wg sync.WaitGroup

	wg.Add(1)
	go httpAdapter.Start(context.Background(), &wg)

	wg.Wait()
}
