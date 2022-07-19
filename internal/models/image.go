package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Image model
type Image struct {
	ID       uuid.UUID `json:"id,omitempty" readonly:"true"`
	Filename string    `json:"filename" bson:"filename"`
	RootPath string    `json:"-" db:"root_path" bson:"root_path" swaggerignore:"true"`
	Data     *[]byte   `json:"-" bson:"data,omitempty" swaggerignore:"true"`
	EasyLink string    `json:"easy_link" bson:"easy_link"`
}

// Byte Get bytes from Image (used only RootPath)
func (img *Image) Byte() (*[]byte, error) {
	dat, err := os.ReadFile(img.RootPath + img.Filename)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &dat, err
}

// FullPath Return full path fof image (rootPath + filename)
func (img Image) FullPath() string {
	return fmt.Sprintf("%s/%s", img.RootPath, img.Filename)
}

// NewImage Constructor Image model
func NewImage(filename string, dt *[]byte) Image {
	pwd, _ := os.Getwd()
	fName := filepath.Clean(filename)
	nowDate := time.Now()
	rootPath := fmt.Sprintf("%v/images/%d-%s-%d/", pwd, nowDate.Day(), nowDate.Month().String(), nowDate.Year())
	easyLink := fmt.Sprintf("%d%s%d_%v", nowDate.Day(), nowDate.Month(), nowDate.Year(), fName)
	return Image{
		Filename: fName,
		RootPath: rootPath,
		Data:     dt,
		EasyLink: easyLink,
	}
}
