package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/abiyogaaron/simplebank-service/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load configuration file: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
