package test

import (
	"app/internal/adapters/grpc/protocGen"
	"context"
	"github.com/stretchr/testify/assert"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"testing"
)

var opts []grpc.DialOption
var clientEntity protocGen.EntityClient

func TestMain(m *testing.M) {
	opts = []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:5300", opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	clientEntity = protocGen.NewEntityClient(conn)
	defer conn.Close()

	code := m.Run()
	os.Exit(code)
}

func TestGet(t *testing.T) {
	id_row := "8145b222-b19c-42ff-835b-d4a5048213b5"

	data, err := clientEntity.GetEntityById(context.Background(), &protocGen.GetEntityByIdRequest{EntityId: id_row})
	assert.Nil(t, err)
	assert.NotNil(t, data)
}
