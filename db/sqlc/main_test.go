package mydb

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	
	testDB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
