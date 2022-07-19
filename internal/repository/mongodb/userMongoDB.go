package repository

import (
	"context"
	"os"

	"github.com/google/uuid"

	"app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepoMongoDB implement RepositoryUser like mongoDB
type UserRepoMongoDB struct {
	mongoClient *mongo.Client
	collection  mongo.Collection
}

// NewUserRepoMongoDB constructor
func NewUserRepoMongoDB(client *mongo.Client) *UserRepoMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("USER_COLLECTION"))
	return &UserRepoMongoDB{
		mongoClient: client,
		collection:  *collection,
	}
}

// Get user
func (ur UserRepoMongoDB) Get(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := ur.collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&user)
	return &user, err
}

// Add user
func (ur UserRepoMongoDB) Add(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	_, err := ur.collection.InsertOne(ctx, user)
	return err
}

// Delete user
func (ur UserRepoMongoDB) Delete(ctx context.Context, username string) error {
	res := ur.collection.FindOneAndDelete(ctx, bson.D{{"username", username}})
	return res.Err()
}

// GetAll users
func (ur UserRepoMongoDB) GetAll(ctx context.Context) ([]*models.User, error) {
	cursor, err := ur.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	var resUser []*models.User

	err = cursor.All(ctx, &resUser)
	if err != nil {
		return nil, err
	}
	return resUser, nil
}

// Update user
func (ur UserRepoMongoDB) Update(ctx context.Context, user *models.User) error {
	return nil
}
