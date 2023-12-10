package main

import (
	"database/sql"
	"log"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/server"
	"github.com/aryyawijaya/simple-bank/util"
	_ "github.com/lib/pq"
)

// const (
// 	driverName     = "postgres"
// 	dataSourceName = "postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable"
// 	serverAddress  = "localhost:8080"
// )

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v\n", err)
	}
	// conn, err := sql.Open(driverName, dataSourceName)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v\n", err)
	}

	store := mydb.NewStore(conn)
	server := server.NewServer(store)

	// err = server.Start(serverAddress)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("Cannot start the server: %v\n", err)
	}
}
