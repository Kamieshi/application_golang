package test

import (
	"app/internal/service"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
)

type EntityServerImplement struct {
	EntityServ *service.EntityService
	EntityServer
}

func (e EntityServerImplement) GetEntityById(ctx context.Context, request *GetEntityByIdRequest) (*GetEntityByIdResponse, error) {
	entity, err := e.EntityServ.GetForID(ctx, request.EntityId)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	var messageEntity EntityStruct
	err = protojson.Unmarshal(data, &messageEntity)
	if err != nil {
		return nil, err
	}
	entityResponse := &GetEntityByIdResponse{
		Entity: &messageEntity,
	}
	return entityResponse, err
}
