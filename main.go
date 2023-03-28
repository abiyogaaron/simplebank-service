package main

import (
	"database/sql"
	"log"

	"github.com/abiyogaaron/simplebank-service/api"
	db "github.com/abiyogaaron/simplebank-service/db/sqlc"
	"github.com/abiyogaaron/simplebank-service/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configuration file: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start http server: ", err)
	}
}
