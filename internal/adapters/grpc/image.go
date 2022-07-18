package test

import (
	"app/internal/service"
	"context"
)

type ImageServerImplement struct {
	ImageService *service.ImageService
	ImageManagerServer
}

func (i *ImageServerImplement) GetImageByEasyLink(req *GetImageByIDRequest, resp ImageManager_GetImageByEasyLinkServer) error {
	image, err := i.ImageService.Get(context.Background(), req.EasyLink)
	if err != nil {
		return err
	}
	respMessage := &GetImageByIDResponse{
		MetaData: &ImageStruct{
			FileName: image.Filename,
			Size:     int32(len(*image.Data)),
		},
		Data: *image.Data,
	}
	err = resp.Send(respMessage)
	if err != nil {
		return err
	}
	return nil
}
