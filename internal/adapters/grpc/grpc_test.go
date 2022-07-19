package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"app/internal/config"
	"app/internal/models"
	repository "app/internal/repository/posgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	clientEntity EntityClient
	clientImage  ImageManagerClient
	ctx          context.Context
	connPool     *pgxpool.Pool
)

func TestMain(m *testing.M) {
	err := godotenv.Load("/home/dmitryrusack/Work/application_golang/localConf.env")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	grpcAddress := fmt.Sprintf("%s:%s", conf.GrpcHost, conf.GrpcPort)
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	clientEntity = NewEntityClient(conn)
	clientImage = NewImageManagerClient(conn)
	defer func() {
		if err = conn.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()

	ctx = context.Background()
	connPool, _ = pgxpool.Connect(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.PostgresUser, conf.PostgresPassword, conf.PostgresHost, conf.PostgresPort, conf.PostgresDB),
	)
	code := m.Run()
	os.Exit(code)
}

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
		if err == io.EOF {
			if err = stream.(grpc.ClientStream).CloseSend(); err != nil {
				log.WithError(err).Error()
			}
			break
		}
		assert.Equal(t, imageByResponse.GetMetaData().GetSize(), int32(len(imageByResponse.GetData())))
	}
}
