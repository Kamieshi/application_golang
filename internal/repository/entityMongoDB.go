package repository

import (
	"app/internal/service/models"
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoEntityMongoDB struct {
	mongoClient      *mongo.Client
	collectionEntity mongo.Collection
}

func NewRepoEntityMongoDB(client mongo.Client) RepoEntityMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("ENTITY_COOLLECTION"))
	return RepoEntityMongoDB{
		mongoClient:      &client,
		collectionEntity: *collection,
	}
}

func (rm *RepoEntityMongoDB) GetAll(ctx context.Context) ([]models.Entity, error) {
	cursor, err := rm.collectionEntity.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	var buferEntity models.Entity
	var resEntitys []models.Entity
	for cursor.Next(ctx) {
		err = cursor.Decode(&buferEntity)
		if err != nil {
			return nil, err
		}
		resEntitys = append(resEntitys, buferEntity)
	}
	return resEntitys, nil
}

func (rm *RepoEntityMongoDB) GetForID(ctx context.Context, id string) (models.Entity, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Entity{}, err
	}
	var entity models.Entity

	err = rm.collectionEntity.FindOne(ctx, bson.D{{"_id", Id}}).Decode(&entity)
	if err != nil {
		return entity, err
	}
	fmt.Println(Id)
	return entity, nil
}

func (rm *RepoEntityMongoDB) Add(ctx context.Context, obj models.Entity) error {
	obj.Id = primitive.NewObjectID()
	_, err := rm.collectionEntity.InsertOne(ctx, obj)
	return err
}

func (rm *RepoEntityMongoDB) Update(ctx context.Context, id string, obj models.Entity) error {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if obj.Id != Id {
		return errors.New("Different ObjectId in query and Body")
	}
	res := rm.collectionEntity.FindOneAndUpdate(ctx, bson.D{{"_id", Id}}, obj)
	err = res.Err()
	return err
}

func (rm *RepoEntityMongoDB) Delete(ctx context.Context, id string) error {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res := rm.collectionEntity.FindOneAndDelete(ctx, bson.D{{"_id", Id}})
	err = res.Err()
	return err

}
