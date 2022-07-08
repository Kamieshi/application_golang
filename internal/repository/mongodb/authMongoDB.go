package repository

import (
	"app/internal/models"
	"context"
	"github.com/google/uuid"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepoMongoDB struct {
	mongoClient *mongo.Client
	collection  mongo.Collection
}

func NewAuthRepoMongoDB(client mongo.Client) AuthRepoMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("SESSION_COLLECTION"))
	return AuthRepoMongoDB{
		mongoClient: &client,
		collection:  *collection,
	}
}

func (ar AuthRepoMongoDB) Create(ctx context.Context, session *models.Session) error {
	session.ID = uuid.New()
	_, err := ar.collection.InsertOne(ctx, session)
	return err
}

func (ar AuthRepoMongoDB) Update(ctx context.Context, session *models.Session) error {
	_, err := ar.collection.ReplaceOne(ctx, bson.D{{"_id", session.ID}}, session)
	if err != nil {
		return err
	}
	return nil
}

func (ar AuthRepoMongoDB) Get(ctx context.Context, SessionId string) (*models.Session, error) {
	var session models.Session
	err := ar.collection.FindOne(ctx, bson.D{{"session_id", SessionId}}).Decode(&session)
	return &session, err
}
func (ar AuthRepoMongoDB) Delete(ctx context.Context, sessionId string) error {
	res := ar.collection.FindOneAndDelete(ctx, bson.D{{"session_id", sessionId}})
	err := res.Err()
	return err
}
