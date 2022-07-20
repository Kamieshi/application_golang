package handlers

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
	"google.golang.org/grpc"

	"app/internal/config"
)

const (
	ContextDirForApp = "/home/dmitryrusack/Work/application_golang"
	secretKey        = "123"
)

var (
	clientEntity EntityClient       //nolint:gochecknoglobals
	clientImage  ImageManagerClient //nolint:gochecknoglobals
	ctx          context.Context    //nolint:gochecknoglobals
	connPool     *pgxpool.Pool      //nolint:gochecknoglobals
	addrAPIEcho  string             //nolint:gochecknoglobals
	addrRPC      string             //nolint:gochecknoglobals
)

func TestMain(m *testing.M) {
	ctx = context.Background()
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

	if err = pool.Retry(func() error {
		conStr := fmt.Sprintf("postgres://postgres:%s@localhost:%s/postgres", "postgres", appPostgres.GetPort("5432/tcp"))
		connPool, err = pgxpool.Connect(ctx, conStr)
		if err != nil {
			return err
		}

		return connPool.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
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

				fmt.Sprintf("GRPC_PORT=%s", configuration.GrpcPort),
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

	addrRPC = fmt.Sprintf("127.0.0.1:%s", appAPI.GetPort(fmt.Sprintf("%s/tcp", configuration.GrpcPort)))
	addrAPIEcho = fmt.Sprintf("http://127.0.0.1:%s", appAPI.GetPort(fmt.Sprintf("%s/tcp", configuration.EchoPort)))
	ctx = context.Background()

	// Wait start api
	if err = pool.Retry(func() error {
		if resp, errResp := http.Get(addrAPIEcho + "/ping"); errResp != nil {
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

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(addrRPC, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	clientEntity = NewEntityClient(conn)
	clientImage = NewImageManagerClient(conn)
	defer func() {
		if err = conn.Close(); err != nil {
			log.WithError(err).Error()
		}
	}()

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
