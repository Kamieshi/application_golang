package http

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"testing"
)

var URLCreateUser, urlLogin, urlCheckAuth, urlLogOut, urlRefresh string
var urlCreateEntity, urlGetByIdEntity, urlGetAllEntity, urlUpdateEntity, urlDeleteEntity string

const (
	//TODO reformat
	pathToMigrations = "/home/dmitryrusack/Work/application_golang/migrations"
	ContextDirForApp = "/home/dmitryrusack/Work/application_golang"
)

var addrAPI string

var connPullDB *pgxpool.Pool
var ctx context.Context

var secretKey = "123"

func TestMain(m *testing.M) {
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
		Cmd:        []string{fmt.Sprintf("-url=jdbc:postgresql://%s:5432/postgres -schemas=public -user=postgres -password=postgres -connectRetries=10 migrate", appPostgres.Container.NetworkSettings.IPAddress)},
		Mounts:     []string{fmt.Sprintf("%s:/flyway/sql", pathToMigrations)},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer closer(appFlyWay)

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
				fmt.Sprintf("SECRET_KEY=%s", secretKey),
				fmt.Sprintf("REDIS_URL=%s:6379", appRedis.Container.NetworkSettings.IPAddress),
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	defer closer(appApi)

	addrAPI = fmt.Sprintf("http://127.0.0.1:%s", appApi.GetPort("8005/tcp"))
	ctx = context.Background()
	//Wait start api
	if err = pool.Retry(func() error {
		_, err = http.Get(addrAPI + "/ping")
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
	urlUpdateEntity = addrAPI + "/entity/"
	urlDeleteEntity = addrAPI + "/entity/"

	//Init connectionPull
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
	if err := pool.Purge(appApi); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(appRedis); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

//func TestMain(m *testing.M) {
//	godotenv.Load("/home/dmitryrusack/Work/application_golang/.env")
//	secretKey = os.Getenv("SECRET_KEY")
//	addrApi = "http://127.0.0.1:8005"
//	URLCreateUser = addrApi + "/user"
//	urlLogin = addrApi + "/auth/login"
//	urlCheckAuth = addrApi + "/auth/info"
//	urlLogOut = addrApi + "/auth/logout"
//	urlRefresh = addrApi + "/auth/refresh"
//	urlCreateEntity = addrApi + "/entity"
//	urlGetAllEntity = addrApi + "/entity"
//	urlGetByIdEntity = addrApi + "/entity/"
//	urlUpdateEntity = addrApi + "/entity/"
//	urlDeleteEntity = addrApi + "/entity/"
//	ctx = context.Background()
//	connPullDB, _ = pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5433/postgres")
//	code := m.Run()
//	os.Exit(code)
//}
