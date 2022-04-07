package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
		log.Printf("error opening database %v", err)
		os.Exit(1)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Printf("error create instance: %v", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"store.db",
		driver,
	)
	if err != nil {
		log.Printf("failed to create migrator: %v", err)
		os.Exit(1)
	}

	m.Up()
}
