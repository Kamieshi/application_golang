package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"

	"app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RepoEntityMongoDB Implement interface Repository entity like mongodb
type RepoEntityMongoDB struct {
	mongoClient      *mongo.Client
	collectionEntity *mongo.Collection
}

// NewRepoEntityMongoDB Constructor
func NewRepoEntityMongoDB(client *mongo.Client) *RepoEntityMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("ENTITY_COLLECTION"))
	return &RepoEntityMongoDB{
		mongoClient:      client,
		collectionEntity: collection,
	}
}

// GetAll Objects from collections entity
func (rm RepoEntityMongoDB) GetAll(ctx context.Context) ([]*models.Entity, error) {
	cursor, err := rm.collectionEntity.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	var resEntities []*models.Entity

	err = cursor.All(ctx, &resEntities)
	if err != nil {
		return nil, err
	}

	return resEntities, nil
}

// GetForID Get object by id
func (rm RepoEntityMongoDB) GetForID(ctx context.Context, id string) (*models.Entity, error) {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entity models.Entity

	err = rm.collectionEntity.FindOne(ctx, bson.D{{"_id", ID}}).Decode(&entity)
	if err != nil {
		return nil, err
	}
	fmt.Println(ID)
	return &entity, nil
}

// Add Write new object in mongoDB
func (rm RepoEntityMongoDB) Add(ctx context.Context, obj *models.Entity) error {
	obj.ID = uuid.New()
	_, err := rm.collectionEntity.InsertOne(ctx, obj)
	return err
}

// Update entity file
func (rm RepoEntityMongoDB) Update(ctx context.Context, id string, obj *models.Entity) error {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	obj.ID = uuid.New()
	updateData := bson.M{
		"$set": obj,
	}
	res := rm.collectionEntity.FindOneAndUpdate(ctx, bson.D{{"_id", ID}}, updateData)
	err = res.Err()
	return err
}

// Delete entity from collection
func (rm RepoEntityMongoDB) Delete(ctx context.Context, id string) error {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res := rm.collectionEntity.FindOneAndDelete(ctx, bson.D{{"_id", ID}})
	err = res.Err()
	return err
}
