package repository

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"app/internal/models"
)

// ImageRepoMongoDB implement interface RepositoryImage like mongoDB
type ImageRepoMongoDB struct {
	mongoClient     *mongo.Client
	collectionImage *mongo.Collection
}

// NewImageRepoMongoDB Constructor
func NewImageRepoMongoDB(client *mongo.Client) *ImageRepoMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("IMAGE_COLLECTION"))
	return &ImageRepoMongoDB{
		mongoClient:     client,
		collectionImage: collection,
	}
}

// Save image
func (repImg *ImageRepoMongoDB) Save(ctx context.Context, img *models.Image) error {
	bytesArr := *img.Data
	img.Data = nil
	filter := bson.M{
		"root_path": img.RootPath,
		"filename":  img.Filename,
	}
	bsonObj := bson.M{
		"$set": img,
	}
	opts := options.Update().SetUpsert(true)
	res, err := repImg.collectionImage.UpdateOne(ctx, filter, bsonObj, opts)
	if err != nil {
		logrus.WithFields(logrus.Fields{"full_path": img.FullPath()}).Error("Unsuccessful write in mongodb")
		logrus.WithError(err).Error()
		return err
	}
	img.Data = &bytesArr
	logrus.Info(res)
	return err
}

// Get image
func (repImg *ImageRepoMongoDB) Get(ctx context.Context, easyLink string) (*models.Image, error) {
	var image models.Image
	query := bson.D{primitive.E{Key: "easy_link", Value: easyLink}}
	err := repImg.collectionImage.FindOne(ctx, query).Decode(&image)
	if err != nil {
		logrus.WithError(err).Error("Error get image in db")
		return &image, err
	}
	return &image, err
}

// Delete image
func (repImg *ImageRepoMongoDB) Delete(ctx context.Context, id uuid.UUID) error {
	query := bson.D{primitive.E{Key: "_id", Value: id}}
	res := repImg.collectionImage.FindOneAndDelete(ctx, query)
	return res.Err()
}
