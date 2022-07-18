package handlers

import (
	gr "app/internal/adapters/grpc/protocGen"
	"app/internal/service"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
)

type EntityServerImplement struct {
	EntityServ *service.EntityService
	gr.EntityServer
}

func (e EntityServerImplement) GetEntityById(ctx context.Context, request *gr.GetEntityByIdRequest) (*gr.GetEntityByIdResponse, error) {
	entity, err := e.EntityServ.GetForID(ctx, request.EntityId)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	var messageEntity gr.EntityStruct
	err = protojson.Unmarshal(data, &messageEntity)
	if err != nil {
		return nil, err
	}
	entityResponse := &gr.GetEntityByIdResponse{
		Entity: &messageEntity,
	}
	return entityResponse, err
}
