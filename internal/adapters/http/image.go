package http

import (
	"app/internal/service"
	ech "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// ImageHandler Handler for work with image service
type ImageHandler struct {
	ImageService *service.ImageService
}

// Load godoc
// Load image to app
// @tags Images
// @accept mpfd
// @Summary Load image
// @Description Load Image
// @Security ApiKeyAuth
// @Param image formData file true "File"
// @Success 200 {object} models.Image
// @Failure 500 {string} Error Parse input data
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /upload [post]
func (ih ImageHandler) Load(c ech.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	data, err := file.Open()
	if err != nil {
		logrus.WithError(err).Error()
		return err
	}
	defer func() {
		err = data.Close()
		if err != nil {
			logrus.WithError(err).Error()
		}
	}()

	loadData, err := io.ReadAll(data)
	if err != nil {
		logrus.WithError(err).Error()
		return err
	}

	img, err := ih.ImageService.Save(c.Request().Context(), file.Filename, &loadData)
	if err != nil {
		return err
	}

	img.Data = nil
	return c.JSON(http.StatusOK, img)
}

// Get godoc
// Get image from app
// @tags Images
// @Summary Get image
// @Description Get Image
// @Security ApiKeyAuth
// @Param easy_link path string true "file easy link"
// @Success 200 {file} file
// @Failure 404 {string} Not found
// @Failure 400 {string} Missing jwt token
// @Failure 401 {string} unAuthorized
// @Router /load/{easy_link} [get]
func (ih ImageHandler) Get(c ech.Context) error {
	easyLink := c.Param("easy_link")

	img, err := ih.ImageService.Get(c.Request().Context(), easyLink)
	if err != nil {
		logrus.WithError(err).Error("Handler error")
		return c.String(http.StatusNotFound, "Not found")
	}

	return c.File(img.FullPath())
}
