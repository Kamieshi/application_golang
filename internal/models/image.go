package models

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type Image struct {
	Filename string  `json:"filename" bson:"filename"`
	RootPath string  `json:"root_path" bson:"root_path"`
	Data     *[]byte `json:"data,omitempty" bson:"data,omitempty"`
	EasyLink string  `json:"easy_link" bson:"easy_link"`
}

func (im *Image) Byte() (*[]byte, error) {
	dat, err := os.ReadFile(im.RootPath)
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
	easylink := fmt.Sprintf("%d%s%d_%v", nowDate.Day(), nowDate.Month(), nowDate.Year(), fName)
	return Image{
		Filename: fName,
		RootPath: rootPath,
		Data:     dt,
		EasyLink: easylink,
	}
}
