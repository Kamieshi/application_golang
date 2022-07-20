package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"

	"app/internal/config"
)

const (
	ContextDirForApp = "/home/dmitryrusack/Work/application_golang"
	secretKey        = "123"
)

var (
	addrAPI          string          //nolint:gochecknoglobals
	connPullDB       *pgxpool.Pool   //nolint:gochecknoglobals
	ctx              context.Context //nolint:gochecknoglobals
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

	appRedis, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "redis",
		Name:       "redis",
		Repository: "redis",
		Tag:        "latest",
	})
	if err != nil {
		log.Fatal(err)
	}

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
				fmt.Sprintf("ECHO_PORT=%s", configuration.EchoPort),
				fmt.Sprintf("SECRET_KEY=%s", secretKey),
				fmt.Sprintf("REDIS_URL=%s:6379", appRedis.Container.NetworkSettings.IPAddress),
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	addrAPI = fmt.Sprintf("http://127.0.0.1:%s", appAPI.GetPort(fmt.Sprintf("%s/tcp", configuration.EchoPort)))
	ctx = context.Background()

	// Wait start api
	if err = pool.Retry(func() error {
		if resp, errResp := http.Get(addrAPI + "/ping"); errResp != nil {
			if resp != nil {
				defer func() {
					if err = resp.Body.Close(); err != nil {
						log.WithError(err).Error()
					}
				}()
				return err
			}
			return errors.New("cannot send to ping")
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
