package handlers

import (
	"app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type ImageHandler struct {
	ImageService service.ImageService
}

func (ih ImageHandler) Load(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	data, err := file.Open()
	if err != nil {
		return err
	}
	defer data.Close()
	loadData, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	img, err := ih.ImageService.Save(c.Request().Context(), file.Filename, &loadData)
	if err != nil {
		return err
	}
	img.Data = nil
	return c.JSON(http.StatusAccepted, img)
}

func (ih ImageHandler) Get(c echo.Context) error {
	easy_link := c.Param("easy_link")
	img, err := ih.ImageService.Get(c.Request().Context(), easy_link)
	if err != nil {
		logrus.WithError(err).Error("Handler error")
		return err
	}
	return c.File(img.FullPath())
}
