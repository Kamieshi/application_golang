package handlers

import (
	gr "app/internal/adapters/grpc/protocGen"
	"app/internal/service"
	"context"
)

type ImageServerImplement struct {
	ImageService *service.ImageService
	gr.ImageManagerServer
}

func (i *ImageServerImplement) GetImageByEasyLink(req *gr.GetImageByIDRequest, resp gr.ImageManager_GetImageByEasyLinkServer) error {
	image, err := i.ImageService.Get(context.Background(), req.EasyLink)
	if err != nil {
		return err
	}
	respMessage := &gr.GetImageByIDResponse{
		MetaData: &gr.ImageStruct{
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
