package main

import (
	"database/sql"
	"log"

	"github.com/aryyawijaya/simple-bank/api"
	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable"
	serverAddress  = "localhost:8080"
)

func main() {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v\n", err)
	}

	store := mydb.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalf("Cannot start the server: %v\n", err)
	}
}
