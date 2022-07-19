// Package repository work with repository
package repository

import (
	"context"
	"os"

	"app/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRepoMongoDB Instance mongo rep
type AuthRepoMongoDB struct {
	mongoClient *mongo.Client
	collection  mongo.Collection
}

// NewAuthRepoMongoDB Constructor method
func NewAuthRepoMongoDB(client *mongo.Client) *AuthRepoMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("SESSION_COLLECTION"))
	return &AuthRepoMongoDB{
		mongoClient: client,
		collection:  *collection,
	}
}

// Create new session
func (ar *AuthRepoMongoDB) Create(ctx context.Context, session *models.Session) error {
	session.ID = uuid.New()
	_, err := ar.collection.InsertOne(ctx, session)
	return err
}

// Update session
func (ar *AuthRepoMongoDB) Update(ctx context.Context, session *models.Session) error {
	_, err := ar.collection.ReplaceOne(ctx, bson.D{{"_id", session.ID}}, session)
	if err != nil {
		return err
	}
	return nil
}

// Get session
func (ar *AuthRepoMongoDB) Get(ctx context.Context, SessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := ar.collection.FindOne(ctx, bson.D{{"session_id", SessionID}}).Decode(&session)
	return &session, err
}

// Delete session
func (ar *AuthRepoMongoDB) Delete(ctx context.Context, sessionID uuid.UUID) error {
	res := ar.collection.FindOneAndDelete(ctx, bson.D{{"session_id", sessionID}})
	err := res.Err()
	return err
}

// Disable session(session.Disabled false->true)
func (ar *AuthRepoMongoDB) Disable(ctx context.Context, sessionID uuid.UUID) error {
	return nil
}
