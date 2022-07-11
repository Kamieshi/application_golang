package tests

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
	"time"
)

var pgPool *pgxpool.Pool
var ctx = context.Background()

func TestMain(t *testing.M) {
	pwdPath, _ := os.Getwd()
	err := godotenv.Load(pwdPath + "/.env")
	if err != nil {
		panic(err)
	}
	configuration, err := config.GetConfig()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	//pwdDir, _ := os.Getwd()
	resource, err := pool.BuildAndRun(
		"pg_for_test",
		pwdPath+"/dockerfile",
		[]string{fmt.Sprintf("POSTGRES_PASSWORD=%s", configuration.POSTGRES_PASSWORD)},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Retry(func() error {
		var err error
		conStr := fmt.Sprintf("postgres://postgres:%s@localhost:%s/postgres", configuration.POSTGRES_PASSWORD, resource.GetPort("5432/tcp"))
		pgPool, err = pgxpool.Connect(context.Background(), conStr)
		if err != nil {
			return err
		}
		tContext, _ := context.WithTimeout(ctx, time.Duration(1*time.Second))
		return pgPool.Ping(tContext)
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
