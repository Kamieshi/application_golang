package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
)

const (
	pathToMigrations = "/home/dmitryrusack/Work/application_golang/migrations"
	ContextDirForApp = "/home/dmitryrusack/Work/application_golang"
)

var (
	addrAPI          string          //nolint:gochecknoglobals
	connPullDB       *pgxpool.Pool   //nolint:gochecknoglobals
	ctx              context.Context //nolint:gochecknoglobals
	secretKey        = "123"         //nolint:gochecknoglobals
	URLCreateUser    string          //nolint:gochecknoglobals
	urlLogin         string          //nolint:gochecknoglobals
	urlCheckAuth     string          //nolint:gochecknoglobals
	urlLogOut        string          //nolint:gochecknoglobals
	urlRefresh       string          //nolint:gochecknoglobals
	urlCreateEntity  string          //nolint:gochecknoglobals
	urlGetByIdEntity string          //nolint:gochecknoglobals
	urlGetAllEntity  string          //nolint:gochecknoglobals
	urlDeleteEntity  string          //nolint:gochecknoglobals
)

func TestMain(m *testing.M) { //nolint:funlen
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	closer := func(resource *dockertest.Resource) {
		if resource != nil {
			if err = resource.Close(); err != nil {
				log.WithError(err).Error()
			}
		}
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
	defer closer(appPostgres)
	appFlyWay, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "flyWay",
		Name:       "flyWay",
		Repository: "flyway/flyway",
		Tag:        "latest",
		Env:        nil,
		Entrypoint: nil,
		Cmd: []string{
			fmt.Sprintf(
				"-url=jdbc:postgresql://%s:5432/postgres -schemas=public -user=postgres -password=postgres -connectRetries=10 migrate",
				appPostgres.Container.NetworkSettings.IPAddress)},
		Mounts: []string{fmt.Sprintf("%s:/flyway/sql", pathToMigrations)},
	})

	defer closer(appFlyWay)
	if err != nil {
		log.Fatal(err)
	}
	appRedis, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "redis",
		Name:       "redis",
		Repository: "redis",
		Tag:        "latest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer closer(appRedis)

	appAPI, err := pool.BuildAndRunWithBuildOptions(
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
				fmt.Sprintf("SECRET_KEY=%s", secretKey),
				fmt.Sprintf("REDIS_URL=%s:6379", appRedis.Container.NetworkSettings.IPAddress),
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	defer closer(appAPI)

	addrAPI = fmt.Sprintf("http://127.0.0.1:%s", appAPI.GetPort("8005/tcp"))
	ctx = context.Background()

	// Wait start api
	if err = pool.Retry(func() error {
		if resp, err := http.Get(addrAPI + "/ping"); err != nil {
			if err = resp.Body.Close(); err != nil {
				return err
			}
		}
		return err
	}); err != nil {
		log.Fatalf("Could not connect to API: %s", err)
	}

	URLCreateUser = addrAPI + "/user"
	urlLogin = addrAPI + "/auth/login"
	urlCheckAuth = addrAPI + "/auth/info"
	urlLogOut = addrAPI + "/auth/logout"
	urlRefresh = addrAPI + "/auth/refresh"

	urlCreateEntity = addrAPI + "/entity"
	urlGetAllEntity = addrAPI + "/entity"
	urlGetByIdEntity = addrAPI + "/entity/"
	urlDeleteEntity = addrAPI + "/entity/"

	// Init connectionPull
	if err = pool.Retry(func() error {
		conStr := fmt.Sprintf("postgres://postgres:%s@%s:5432/postgres", "postgres", appPostgres.Container.NetworkSettings.IPAddress)
		connPullDB, err = pgxpool.Connect(ctx, conStr)
		if err != nil {
			return err
		}

		return connPullDB.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	if err = os.Setenv("SECRET_KEY", secretKey); err != nil {
		log.WithError(err).Panic()
	}

	code := m.Run()

	if err := pool.Purge(appPostgres); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appFlyWay); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appAPI); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appRedis); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
