package repository

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"app/internal/models"
)

const maxFileSize = 200

// WriteImageInHost Write image into host machine
func WriteImageInHost(image models.Image) error {
	bytes := image.Data
	err := os.MkdirAll(image.RootPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("imageFileManager.go/WriteImageInHost : %v", err)
	}
	if _, err = os.Stat(image.RootPath); os.IsNotExist(err) {
		err = os.Mkdir(image.RootPath, os.ModePerm)
		if err != nil {
			logrus.WithFields(logrus.Fields{"path": image.FullPath()}).Error("no permission")
			return fmt.Errorf("imageFileManager.go/WriteImageInHost : %v", err)
		}
	}
	if err != nil {
		return fmt.Errorf("imageFileManager.go/WriteImageInHost : %v", err)
	}

	err = os.WriteFile(image.FullPath(), *bytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("imageFileManager.go/WriteImageInHost : %v", err)
	}
	return nil
}

// CheckImageData Check exist and access to image file
func CheckImageData(image *models.Image) error {
	file, err := os.Open(image.FullPath())
	if err != nil {
		return fmt.Errorf("imageFileManager.go/CheckImageData : %v", err)
	}
	data := make([]byte, maxFileSize)
	if err != nil {
		return fmt.Errorf("imageFileManager.go/CheckImageData : %v", err)
	}
	for {
		_, bite := file.Read(data)
		if bite == io.EOF {
			break
		}
	}
	return nil
}
