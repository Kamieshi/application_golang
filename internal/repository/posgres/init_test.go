package repository

import (
	"app/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var pgPool *pgxpool.Pool
var ctx = context.Background()

func TestMain(t *testing.M) {
	configuration, err := config.GetConfig("/home/rusak/application_golang/localConfig.env")
	if err != nil {
		log.WithError(err).Panic()
	}

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
		log.Fatal(err)
	}

	appFlyWay, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "flyWay",
		Name:       "flyWay",
		Repository: "flyway/flyway",
		Tag:        "latest",
		Env:        nil,
		Entrypoint: nil,
		Cmd:        []string{fmt.Sprintf("-url=jdbc:postgresql://%s:5432/postgres -schemas=public -user=postgres -password=postgres -connectRetries=10 migrate", appPostgres.Container.NetworkSettings.IPAddress)},
		Mounts:     []string{fmt.Sprintf("%s:/flyway/sql", configuration.PathToMigration)},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.WithError(err).Panic()
	}

	if err := pool.Retry(func() error {
		var err error
		conStr := fmt.Sprintf("postgres://postgres:%s@localhost:%s/postgres", configuration.PostgresPassword, appPostgres.GetPort("5432/tcp"))
		pgPool, err = pgxpool.Connect(context.Background(), conStr)
		if err != nil {
			return err
		}
		tContext, tContextCancel := context.WithTimeout(ctx, time.Duration(1*time.Second))
		defer tContextCancel()
		return pgPool.Ping(tContext)
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := t.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(appPostgres); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appFlyWay); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
