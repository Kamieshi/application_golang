package repository

import (
	"application/config"
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type ShardPostgres struct {
	host, port, username, password, default_database string
	conn                                             *pgx.Conn
}

func CreateShardPostgress(host, port, username, password, defaultDatabase string) (*ShardPostgres, error) {

	connection, err := pgx.Connect(context.Background(), CreateUrlConnect(password, username, host, port, defaultDatabase))
	if err != nil {
		log.Fatal(err)
	}

	var sp ShardPostgres = ShardPostgres{
		conn:             connection,
		host:             host,
		port:             port,
		username:         username,
		password:         password,
		default_database: defaultDatabase,
	}
	return &sp, nil
}

func (sp *ShardPostgres) Close() error {
	err := sp.conn.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetPg() (*ShardPostgres, error) {

	conf := config.GetConfig()

	pgCon, err := CreateShardPostgress(conf.POSTGRESS_HOST, conf.POSTGRES_PORT, conf.POSTGRES_USER, conf.POSTGRES_PASSWORD, conf.POSTGRES_DB)
	if err != nil {
		log.Fatal(err)
	}
	return pgCon, err
}
