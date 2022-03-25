package repository

import (
	"fmt"
	"log"
)

func CreateUrlConnect(password, username, host, port, defaulDatabase string) string {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, defaulDatabase)
	log.Println(connectionString)
	return connectionString
}
