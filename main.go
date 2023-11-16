package main

import (
	"database/sql"
	"log"

	"github.com/LKarrie/mdc-server/api"
	db "github.com/LKarrie/mdc-server/db/sqlc"
	"github.com/LKarrie/mdc-server/util"
	"github.com/docker/docker/client"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	stroe := db.NewStore(conn)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("cannot create docker cli:", err)
	}
	d := api.NewDocker(cli)

	server, err := api.NewServer(config, d, stroe)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannt start server:", err)
	}
}
