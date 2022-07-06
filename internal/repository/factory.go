package repository

import (
	"app/internal/config"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type RepoFactory interface {
	GetAuthRepo() RepoSession
	GetEntityRepo() RepoEntity
	GetImageRepo() RepoImage
	GetUserRepo() RepoUser
}

func GetFactory(typeFactory string, config *config.Configuration) RepoFactory {

	switch typeFactory {
	case "mongo":
		timeOutConnect, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		clientMongo, err := mongo.Connect(timeOutConnect, options.Client().ApplyURI(config.ConnectUrlMongo()))
		defer cancel()
		if err != nil {
			log.WithError(err).Panic("Error with mongo connection")
		}

		var rf *readpref.ReadPref
		err = clientMongo.Ping(timeOutConnect, rf)
		if err != nil {
			log.WithError(err).Panic("Error with mongo connection")
		}

		mongoFactory := MongoDBRepoFactory{clientMongo}

		return mongoFactory

	case "pg":
		connPool, err := pgxpool.Connect(context.Background(), config.UrlPostgres())
		if err != nil {
			log.Println("Connecting url", config.UrlPostgres())
			log.Fatal(err)
		}
		pgFactory := PgRepoFactory{connPool}
		return pgFactory
	}
	return nil
}
