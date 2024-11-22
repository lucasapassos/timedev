package db

import (
	_ "embed"
	"log"
	"os"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDBConnection() *sql.DB {

	var database *sql.DB
	sqlpath := os.Getenv("SQLITE_PATH")

	var err error
	database, err = sql.Open("sqlite3", sqlpath)

	if err != nil {
		log.Fatal(err)
	}

	return database
}
