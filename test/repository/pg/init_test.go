package testPg

import (
	"app/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

var pgPool *pgxpool.Pool

func TestMain(t *testing.M) {
	err := godotenv.Load("/home/dmitryrusack/Work/application_golang/.env")
	if err != nil {
		panic(err)
	}
	configuration, err := config.GetConfig()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.BuildAndRun(
		"pg_for_test",
		"/home/dmitryrusack/Work/application_golang/test/repository/pg/dockerfile",
		[]string{fmt.Sprintf("POSTGRES_PASSWORD=%s", configuration.POSTGRES_PASSWORD)},
	)

	if err := pool.Retry(func() error {
		var err error
		conStr := fmt.Sprintf("postgres://postgres:%s@localhost:%s/postgres", configuration.POSTGRES_PASSWORD, resource.GetPort("5432/tcp"))
		pgPool, err = pgxpool.Connect(context.Background(), conStr)
		if err != nil {
			return err
		}
		return pgPool.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := t.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
