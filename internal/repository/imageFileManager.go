package repository

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"app/internal/config"
	"app/internal/models"
)

// WriteImageInHost Write image into host machine
func WriteImageInHost(image models.Image) error {
	bytes := image.Data
	err := os.MkdirAll(image.RootPath, os.ModePerm)
	if err != nil {
		return err
	}

	if _, err = os.Stat(image.RootPath); os.IsNotExist(err) {
		err = os.Mkdir(image.RootPath, os.ModePerm)
		if err != nil {
			logrus.WithFields(logrus.Fields{"path": image.FullPath()}).Error("no permission")
			return err
		}
	}
	if err != nil {
		return err
	}

	err = os.WriteFile(image.FullPath(), *bytes, os.ModePerm)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"root_path": image.RootPath,
		}).Error("Error write files")
		return err
	}
	return err
}

// CheckImageData Check exist and access to image file
func CheckImageData(image *models.Image) error {
	file, err := os.Open(image.FullPath())
	data := make([]byte, config.Config().MaxFileSize)
	if err != nil {
		logrus.WithFields(logrus.Fields{"full_path": image.FullPath()}).Error("Error Open file")
		return err
	}
	for {
		_, bite := file.Read(data)
		if bite == io.EOF {
			break
		}
	}
	return err
}
