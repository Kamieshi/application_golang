package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
	"github.com/sirupsen/logrus"
)

type ImageService struct {
	ImageRepository repository.RepoImage
}

func NewImageService(imagerRepository *repository.RepoImage) *ImageService {
	return &ImageService{
		ImageRepository: *imagerRepository,
	}

}

func (ims ImageService) Save(ctx context.Context, fileName string, data *[]byte) (*models.Image, error) {
	image := models.NewImage(fileName, data)

	err := ims.ImageRepository.Save(ctx, &image)
	if err != nil {
		logrus.Error("Error write image in db")
		return nil, err
	}
	logrus.Info("Successful write image in db")
	err = repository.WriteImageInHost(image)
	if err != nil {
		logrus.Error("Error write image in host")
		err = ims.ImageRepository.Delete(ctx, image.Id)
		if err != nil {
			logrus.Error("Error delete image in db")
			return &image, err
		}
	}
	logrus.Info("Successful SAVE image")
	return &image, err
}

func (ims ImageService) Get(ctx context.Context, easyLink string) (*models.Image, error) {
	image, err := ims.ImageRepository.Get(ctx, easyLink)
	if err != nil {
		logrus.WithError(err).Error("Error Get in repository image")
		return nil, err
	}
	err = repository.CheckImageData(image)
	if err != nil {
		logrus.WithError(err).Error("Error Init image data")
		return nil, err
	}
	return image, err
}

func (ims ImageService) Delete(ctx context.Context, image models.Image) error {
	err := ims.ImageRepository.Delete(ctx, image.Id)
	return err
}
