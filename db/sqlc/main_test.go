package mydb_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/util"
	_ "github.com/lib/pq"
)

// const (
// 	driverName     = "postgres"
// 	dataSourceName = "postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable"
// )

var testQueries *mydb.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("Cannot load config: %v\n", err)
	}

	// testDB, err = sql.Open(driverName, dataSourceName)
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	testQueries = mydb.New(testDB)

	os.Exit(m.Run())
}
