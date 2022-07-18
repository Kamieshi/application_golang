// Package config with config
package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
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
}

// ConnectingURLPostgres  Return connection string to Postgres
func (c *Configuration) ConnectingURLPostgres() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB)
}

// ConnectingURLMongo Return connection string to MongoDB
func (c *Configuration) ConnectingURLMongo() string {
	return fmt.Sprintf("mongodb://%v:%v", c.MongoHost, c.MongoPort)
}

// GetConfig Get instance config object
func GetConfig() (*Configuration, error) {
	conf := Configuration{}

	_, exist := os.LookupEnv("POSTGRES_PORT")
	if !exist {
		err := godotenv.Load("./localConfig.env")
		if err != nil {
			return &conf, fmt.Errorf("config.go/GetConfig Error parse data from file : %v", err)
		}
	}

	err := env.Parse(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
