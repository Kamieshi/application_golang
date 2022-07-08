package repository

import (
	"app/internal/models"
	"context"
	"github.com/google/uuid"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepoMongoDB struct {
	mongoClient *mongo.Client
	collection  mongo.Collection
}

func NewUserRepoMongoDB(client mongo.Client) UserRepoMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("USER_COLLECTION"))
	return UserRepoMongoDB{
		mongoClient: &client,
		collection:  *collection,
	}
}

func (ur UserRepoMongoDB) Get(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := ur.collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&user)
	return &user, err
}

func (ur UserRepoMongoDB) Add(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	_, err := ur.collection.InsertOne(ctx, user)
	return err
}

func (ur UserRepoMongoDB) Delete(ctx context.Context, username string) error {
	res := ur.collection.FindOneAndDelete(ctx, bson.D{{"username", username}})
	return res.Err()
}

func (ur UserRepoMongoDB) GetAll(ctx context.Context) ([]models.User, error) {
	cursor, err := ur.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	var resUser []models.User

	err = cursor.All(ctx, &resUser)
	if err != nil {
		return nil, err
	}
	return resUser, nil
}

func (ur UserRepoMongoDB) Update(ctx context.Context, user *models.User) error {
	//TODO
	return nil
}
