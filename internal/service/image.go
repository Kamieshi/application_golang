package service

import (
	"context"
	"fmt"

	"app/internal/models"
	"app/internal/repository"
)

// ImageService Service for work with image
type ImageService struct {
	ImageRepository repository.RepoImage
}

// NewImageService Constructor
func NewImageService(imagerRepository repository.RepoImage) *ImageService {
	return &ImageService{
		ImageRepository: imagerRepository,
	}
}

// Save image
func (ims ImageService) Save(ctx context.Context, fileName string, data *[]byte) (*models.Image, error) {
	image := models.NewImage(fileName, data)

	if err := ims.ImageRepository.Save(ctx, &image); err != nil {
		return nil, fmt.Errorf("service image/Save : %v", err)
	}
	if err := repository.WriteImageInHost(image); err != nil {
		errRep := ims.ImageRepository.Delete(ctx, image.ID)
		if errRep != nil {
			return nil, fmt.Errorf("service image/Save : %v , service image/Save : %v", err, errRep)
		}
		return nil, fmt.Errorf("service image/Save : %v", err)
	}
	return &image, nil
}

// Get image
func (ims ImageService) Get(ctx context.Context, easyLink string) (*models.Image, error) {
	image, err := ims.ImageRepository.Get(ctx, easyLink)
	if err != nil {
		return nil, fmt.Errorf("service image/Get : %v", err)
	}
	err = repository.CheckImageData(image)
	if err != nil {
		return nil, fmt.Errorf("service image/Get : %v", err)
	}
	return image, err
}

// Delete Image
func (ims ImageService) Delete(ctx context.Context, image models.Image) error {
	err := ims.ImageRepository.Delete(ctx, image.ID)
	if err != nil {
		return fmt.Errorf("service image/Delete : %v", err)
	}
	return err
}
