package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbsource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

/*
	it is the first phase which we made and if we run the test we got two different errors,
	first is:
		go: go.mod file not found in current directory or any parent directory; see 'go help modules'
	second is:
		Failed to connect to database: sql:unknown driver "postgres" (forgoten import?)
		- this is because database sql package just provides a generic interfaces around the sql database
			it needs to be used in conjuntion with a db driver in order to talk to a database engine.
			for this purpose we can use lib/pq, and install it with it's GitHub page,
			and run the installation command in terminal
			so when we get it's dependency and import it because we did not use it's functions go will remove it
			automatically when we save the file, and that is why we use and _ before the package name
*/

// we define it as a global parameter because we want to use it in all tests
// if we click on *Queries, this contains a DBTX which can be either a transaction or connection
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// before we made our sql db object as code below, but when we define our store we want to have access to it,
	//		therefore we should define it as a global variable in line 34
	// conn, err := sql.Open(dbDriver, dbsource)
	// if err != nil {
	// 	log.Fatal("Failed to connect to database: ", err)
	// }
	// testQueries = New(conn)
	// os.Exit(m.Run())
	var err error
	testDB, err = sql.Open(dbDriver, dbsource)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
