package repository

import (
	"app/internal/models"
	"github.com/sirupsen/logrus"
	"os"
)

func WriteImageInHost(image models.Image) error {
	bytes := image.Data

	if _, err := os.Stat(image.RootPath); os.IsNotExist(err) {
		err = os.Mkdir(image.RootPath, os.ModePerm)
		if err != nil {
			logrus.WithFields(logrus.Fields{"path": image.FullPath()}).Error("no permission")
			return err
		}
	}
	err := os.WriteFile(image.FullPath(), *bytes, 0644)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"root_path": image.RootPath,
		}).Error("Error write files")
		return err
	}
	return err
}
