package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	POSTGRES_PASSWORD, POSTGRES_USER, POSTGRES_DB, POSTGRES_HOST, POSTGRES_PORT, MONGO_HOST, MONGO_PORT string
}

func Load() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func (c Configuration) UrlPosgres() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", c.POSTGRES_USER, c.POSTGRES_PASSWORD, c.POSTGRES_HOST, c.POSTGRES_PORT, c.POSTGRES_DB)
}

func (c Configuration) CoonnectUrlMongo() string {
	return fmt.Sprintf("mongodb://%v:%v", c.MONGO_HOST, c.MONGO_PORT)
}

func (c *Configuration) BaseInit() error {
	c.POSTGRES_DB = os.Getenv("POSTGRES_DB")
	c.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	c.POSTGRES_USER = os.Getenv("POSTGRES_USER")
	c.POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	c.POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	c.MONGO_HOST = os.Getenv("MONGO_HOST")
	c.MONGO_PORT = os.Getenv("MONGO_PORT")
	return nil
}

func GetConfig() (*Configuration, error) {
	conf := Configuration{}
	_, exists := os.LookupEnv("POSTGRES_PORT")

	if !exists {
		err := Load()
		if err != nil {
			return nil, err
		}
	}

	err := conf.BaseInit()
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
