package repository

import (
	"app/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepoFactory struct {
	MongoClient *mongo.Client
}

func (m MongoDBRepoFactory) GetAuthRepo() RepoSession {
	return repository.NewAuthRepoMongoDB(*m.MongoClient)
}

func (m MongoDBRepoFactory) GetEntityRepo() RepoEntity {
	return repository.NewRepoEntityMongoDB(*m.MongoClient)
}

func (m MongoDBRepoFactory) GetImageRepo() RepoImage {
	return repository.NewImageRepoMongoDB(*m.MongoClient)
}

func (m MongoDBRepoFactory) GetUserRepo() RepoUser {
	return repository.NewUserRepoMongoDB(*m.MongoClient)
}
