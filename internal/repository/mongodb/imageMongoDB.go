package repository

import (
	"app/internal/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type ImageRepoMongoDB struct {
	mongoClient     *mongo.Client
	collectionImage mongo.Collection
}

func NewImageRepoMongoDB(client mongo.Client) ImageRepoMongoDB {
	collection := client.Database(os.Getenv("APP_MONGO_DB")).Collection(os.Getenv("IMAGE_COLLECTION"))
	return ImageRepoMongoDB{
		mongoClient:     &client,
		collectionImage: *collection,
	}
}

func (repImg ImageRepoMongoDB) Save(ctx context.Context, img models.Image) (interface{}, error) {

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
		return nil, err
	}
	img.Data = &bytesArr
	logrus.Info(res)
	return res.UpsertedID, err
}

func (repImg ImageRepoMongoDB) Get(ctx context.Context, easyLink string) (*models.Image, error) {
	var image models.Image
	err := repImg.collectionImage.FindOne(ctx, bson.D{{"easy_link", easyLink}}).Decode(&image)
	if err != nil {
		logrus.WithError(err).Error("Error get image in db")
		return &image, err
	}
	return &image, err
}

func (repImg ImageRepoMongoDB) Delete(ctx context.Context, id interface{}) error {
	return nil
}
