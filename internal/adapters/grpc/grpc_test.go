package test

import (
	"app/internal/config"
	"app/internal/models"
	repository "app/internal/repository/posgres"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"io"

	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"testing"
)

var opts []grpc.DialOption
var clientEntity EntityClient
var clientImage ImageManagerClient
var ctx context.Context
var connPool *pgxpool.Pool

func TestMain(m *testing.M) {
	godotenv.Load("/home/dmitryrusack/Work/application_golang/.env")
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	opts = []grpc.DialOption{
		grpc.WithInsecure(),
	}
	grpcAddress := fmt.Sprintf("%s:%s", conf.GrpcHost, conf.GrpcPort)
	conn, err := grpc.Dial(grpcAddress, opts...)
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
	imageByIdRequest := GetImageByIDRequest{EasyLink: "14July2022_Screenshot from 2022-07-14 13-56-01.png"}
	stream, err := clientImage.GetImageByEasyLink(ctx, &imageByIdRequest)
	if err != nil {
		t.Fatal(err)
	}
	for {
		imageByResponse, err := stream.Recv()
		if err == io.EOF {
			if err = stream.CloseSend(); err != nil {
				log.WithError(err).Error()
			}
		}
		assert.Equal(t, imageByResponse.GetMetaData().GetSize(), int32(len(imageByResponse.GetData())))
		break
	}
}
