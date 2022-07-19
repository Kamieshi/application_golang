// Package config with config
package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Configuration configuration for application
type Configuration struct {
	UsedDB           string `env:"USED_DB"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresDB       string `env:"POSTGRES_DB"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	MongoHost        string `env:"MONGO_HOST"`
	MongoPort        string `env:"MONGO_PORT"`
	RedisURL         string `env:"REDIS_URL"`
	GrpcHost         string `env:"GRPC_HOST"`
	GrpcPort         string `env:"GRPC_PORT"`
	GrpcProtocol     string `env:"GRPC_PROTOCOL"`
	EchoPort         string `env:"ECHO_PORT"`
	PathToMigration  string `env:"PATH_TO_MIGRATIONS"`
	MaxFileSize      int64  `env:"MAX_FILE_SIZE" envDefault:"1000"`
}

var _singleConfig *Configuration //nolint:gochecknoglobals

// ConnectingURLPostgres  Return connection string to Postgres
func (c *Configuration) ConnectingURLPostgres() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB)
}

// ConnectingURLMongo Return connection string to MongoDB
func (c *Configuration) ConnectingURLMongo() string {
	return fmt.Sprintf("mongodb://%v:%v", c.MongoHost, c.MongoPort)
}

// Config Get link to simpleConfig
func Config() *Configuration {
	return _singleConfig
}

// GetConfig Get instance config object
func GetConfig(args ...string) (*Configuration, error) {
	if _singleConfig != nil {
		return _singleConfig, nil
	}
	conf := Configuration{}
	if len(args) != 0 {
		absPath := args[0]
		err := godotenv.Load(absPath)
		if err != nil {
			return &conf, fmt.Errorf("config.go/GetConfig Error parse data from file : %v", err)
		}
	}
	_, exist := os.LookupEnv("POSTGRES_PORT")
	if !exist {
		err := godotenv.Load(os.Getenv("PATH_TO_CONFIG_FILE"))
		if err != nil {
			return &conf, fmt.Errorf("config.go/GetConfig Error parse data from file : %v", err)
		}
	}

	err := env.Parse(&conf)
	if err != nil {
		return nil, err
	}
	_singleConfig = &conf
	return _singleConfig, nil
}
