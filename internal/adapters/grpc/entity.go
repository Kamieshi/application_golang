// Package handlers handlers RPC
package handlers

import (
	"context"
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"

	"app/internal/service"
)

// EntityServerImplement implement method from proto-gen
type EntityServerImplement struct {
	EntityServ *service.EntityService
	EntityServer
}

// GetEntityByID get by ID
func (e EntityServerImplement) GetEntityByID(ctx context.Context, request *GetEntityByIdRequest) (*GetEntityByIdResponse, error) {
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
