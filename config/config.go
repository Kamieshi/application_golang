package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	POSTGRES_PASSWORD, POSTGRES_USER, POSTGRES_DB, POSTGRESS_HOST, POSTGRES_PORT string
}

func Load() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unsuccessful attempt to initialize configuration")
	}
	return nil
}

func (c *Configuration) BaseInit() error {
	c.POSTGRES_DB = os.Getenv("POSTGRES_DB")
	c.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	c.POSTGRES_USER = os.Getenv("POSTGRES_USER")
	c.POSTGRESS_HOST = os.Getenv("POSTGRESS_HOST")
	c.POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	return nil
}

func GetConfig() *Configuration {
	conf := Configuration{}
	err := conf.BaseInit()
	if err != nil {
		log.Fatalln(err)
	}
	return &conf
}
