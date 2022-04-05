package repository

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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

func (repImg ImageRepoMongoDB) Save(ctx context.Context, img models.Image) error {

	err := repository.WriteImageInHost(img)
	if err != nil {
		logrus.WithFields(logrus.Fields{"full_path": img.FullPath()}).Error("Unsuccessful write in host")
		return err
	}
	img.Data = nil
	res, err := repImg.collectionImage.InsertOne(ctx, img)
	if err != nil {
		logrus.WithFields(logrus.Fields{"full_path": img.FullPath()}).Error("Unsuccessful write in mongodb")
		return err
	}
	logrus.Info(res)
	return err
}

func (repImg ImageRepoMongoDB) Get(ctx context.Context, easyLink string) (models.Image, error) {
	return models.Image{}, nil
}

func (repImg ImageRepoMongoDB) Delete(ctx context.Context, img models.Image) error {
	return nil
}
