package tests

import (
	"fmt"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

const (
	//TODO reformat
	pathToMigrations = "/home/dmitryrusack/Work/application_golang/migrations"
	ContextDirForApp = "/home/dmitryrusack/Work/application_golang"
)

var addrApi string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	appPostgres, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "postgres",
		Name:       "postgres",
		Repository: "postgres",
		Tag:        "latest",
		Env:        []string{"POSTGRES_PASSWORD=postgres"},
	})
	if err != nil {
		log.Error(err)
	}

	appFlyWay, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "flyWay",
		Name:       "flyWay",
		Repository: "flyway/flyway",
		Tag:        "latest",
		Env:        nil,
		Entrypoint: nil,
		Cmd:        []string{fmt.Sprintf("-url=jdbc:postgresql://%s:5432/postgres -schemas=public -user=postgres -password=postgres -connectRetries=10 migrate", appPostgres.Container.NetworkSettings.IPAddress)},
		Mounts:     []string{fmt.Sprintf("%s:/flyway/sql", pathToMigrations)},
	})
	if err != nil {
		log.Error(err)
	}

	appRedis, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "redis",
		Name:       "redis",
		Repository: "redis",
		Tag:        "latest",
	})
	if err != nil {
		log.Error(err)
	}

	appApi, err := pool.BuildAndRunWithBuildOptions(
		&dockertest.BuildOptions{
			Dockerfile: "dockerfile",
			ContextDir: ContextDirForApp,
		},
		&dockertest.RunOptions{
			Hostname:   "application",
			Name:       "test_app",
			Repository: "",
			Tag:        "",
			Env: []string{
				"POSTGRES_USER=postgres",
				fmt.Sprintf("POSTGRES_HOST=%s", appPostgres.Container.NetworkSettings.IPAddress),
				"POSTGRES_PORT=5432",
				"POSTGRES_PASSWORD=postgres",
				"POSTGRES_DB=postgres",
				fmt.Sprintf("REDIS_URL=%s:6379", appRedis.Container.NetworkSettings.IPAddress),
			},
		})
	if err != nil {
		panic("Can't start api")
	}
	addrApi = fmt.Sprintf("http://127.0.0.1:%s", appApi.GetPort("8005/tcp"))

	code := m.Run()

	if err := pool.Purge(appPostgres); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appFlyWay); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appApi); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appRedis); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
