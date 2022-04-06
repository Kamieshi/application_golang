package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"github.com/sirupsen/logrus"
)

type ImageService struct {
	ImageRepository repository.ImageRepository
}

func (ims ImageService) Save(ctx context.Context, fileName string, data *[]byte) (*models.Image, error) {
	image := models.NewImage(fileName, data)

	id, err := ims.ImageRepository.Save(ctx, image)
	if err != nil {
		logrus.Error("Error write image in db")
		return nil, err
	}
	logrus.Info("Successful write image in db")
	err = repository.WriteImageInHost(image)
	if err != nil {
		logrus.Error("Error write image in host")
		err = ims.ImageRepository.Delete(ctx, id)
		if err != nil {
			logrus.Error("Error delete image in db")
			return &image, err
		}
	}
	logrus.Info("Successful SAVE image")
	return &image, err
}

func (ims ImageService) Get(ctx context.Context, easyLink string) (models.Image, error) {
	image, err := ims.ImageRepository.Get(ctx, easyLink)
	return image, err
}

func (ims ImageService) Delete(ctx context.Context, image models.Image) error {
	err := ims.ImageRepository.Delete(ctx, image)
	return err
}
