package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type Image struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Filename string    `json:"filename" bson:"filename" json:"filename,omitempty"`
	RootPath string    `json:"root_path" bson:"root_path" json:"root_path,omitempty"`
	Data     *[]byte   `json:"data,omitempty" bson:"data,omitempty" json:"data,omitempty"`
	EasyLink string    `json:"easy_link" bson:"easy_link" json:"easy_link,omitempty"`
}

func (img *Image) Byte() (*[]byte, error) {
	dat, err := os.ReadFile(img.RootPath + img.Filename)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &dat, err
}

func (img Image) FullPath() string {
	return fmt.Sprintf("%s/%s", img.RootPath, img.Filename)
}

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
