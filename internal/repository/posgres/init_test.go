package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"

	"app/internal/config"
)

var (
	pgPool *pgxpool.Pool          //nolint:gochecknoglobals
	ctx    = context.Background() //nolint:gochecknoglobals
)

func TestMain(t *testing.M) {
	pathToConfig, err := filepath.Abs("../../../localConf.env")
	if err != nil {
		log.WithError(err).Fatal()
	}
	configuration, err := config.GetConfig(pathToConfig)
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
		Env: []string{
			fmt.Sprintf("FLYWAY_CREATE_SCHEMAS=%s", "true"),
			fmt.Sprintf("FLYWAY_CONNECT_RETRIES_INTERVAL=%d", 2),
			fmt.Sprintf("FLYWAY_CONNECT_RETRIES=%d", 5),
			fmt.Sprintf("FLYWAY_PASSWORD=%s", configuration.PostgresPassword),
			fmt.Sprintf("FLYWAY_USER=%s", configuration.PostgresUser),
			fmt.Sprintf("FLYWAY_SCHEMAS=%s", configuration.PostgresDB),
			fmt.Sprintf("FLYWAY_URL=%s",
				fmt.Sprintf("jdbc:postgresql://%s:5432/%s", appPostgres.Container.NetworkSettings.IPAddress, configuration.PostgresDB)),
			fmt.Sprintf("FLYWAY_BASELINE_ON_MIGRATE=%s", "true"),
		},
		Entrypoint: nil,
		Cmd:        []string{"migrate"},
		Mounts:     []string{fmt.Sprintf("%s:/flyway/sql", configuration.PathToMigration)},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Retry(func() error {
		var err error
		conStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/postgres", configuration.PostgresUser, configuration.PostgresPassword, appPostgres.GetPort("5432/tcp"))
		pgPool, err = pgxpool.Connect(context.Background(), conStr)
		if err != nil {
			return err
		}
		tContext, tContextCancel := context.WithTimeout(ctx, time.Second)
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
