package handlers

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"

	"app/internal/models"
	repository "app/internal/repository/posgres"

	"google.golang.org/grpc"
)

func TestGetEntityById(t *testing.T) {
	dataEntity := &models.Entity{
		Name:     "Testww",
		Price:    10,
		IsActive: true,
	}
	repoEntity := repository.NewRepoEntityPostgres(connPool)
	if err := repoEntity.Add(ctx, dataEntity); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := repoEntity.Delete(ctx, dataEntity.ID.String()); err != nil {
			log.WithError(err).Error()
		}
	})
	data, err := clientEntity.GetEntityById(context.Background(), &GetEntityByIdRequest{EntityId: dataEntity.ID.String()})
	assert.Nil(t, err)
	assert.NotNil(t, data)

	messageEntity, err := protojson.Marshal(data.Entity)
	if err != nil {
		t.Fatal(err)
	}

	var actualEntity models.Entity
	err = json.Unmarshal(messageEntity, &actualEntity)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &actualEntity, dataEntity)
}

func TestGetImagesByEasyLink(t *testing.T) {
	imageByIDRequest := GetImageByIDRequest{EasyLink: "14July2022_Screenshot from 2022-07-14 13-56-01.png"}
	stream, err := clientImage.GetImageByEasyLink(ctx, &imageByIDRequest)
	if err != nil {
		t.Fatal(err)
	}
	for {
		imageByResponse, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				if err = stream.(grpc.ClientStream).CloseSend(); err != nil {
					log.WithError(err).Error()
				}

			}
			assert.Equal(t, imageByResponse.GetMetaData().GetSize(), int32(len(imageByResponse.GetData())))
		}
		t.Fatal(err)
	}
}
