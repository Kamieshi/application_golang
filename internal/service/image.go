package service

import (
	"app/internal/models"
	"app/internal/repository"
	"context"
)

type ImageService struct {
	ImageRepository repository.ImageRepository
}

func (ims ImageService) Save(ctx context.Context, fileName string, data *[]byte) (*models.Image, error) {
	image := models.NewImage(fileName, data)
	err := ims.ImageRepository.Save(ctx, image)
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
