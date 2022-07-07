package testPg

import (
	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	logrus.Fatalf("Could not connect to docker: %s", err)

}
