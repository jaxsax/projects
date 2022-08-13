package main

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jmoiron/sqlx"
)

type app struct {
	dbOptions      *db.Options
	migrateOptions *db.MigrateOptions
	logger         logr.Logger
}

func main() {
	a, err := initialize()
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("sqlite3", a.dbOptions.URI)
	if err != nil {
		a.logger.Error(err, "sql connect")
		return
	}

	if err := db.Ping(); err != nil {
		a.logger.Error(err, "sql ping")
		return
	}

	sqliteInstance, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		a.logger.Error(err, "build instance fromm db")
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"sqlite3",
		sqliteInstance,
	)
	if err != nil {
		a.logger.Error(err, "create database instance")
		return
	}

	if a.migrateOptions.Up && a.migrateOptions.Down {
		a.logger.Error(fmt.Errorf("cannot specify both up and down together"), "invalid options")
		return
	}

	if a.migrateOptions.Up {
		if err := m.Up(); err != nil {
			a.logger.Error(err, "migrate up failed")
			return
		}
	}

	if a.migrateOptions.Down {
		if err := m.Down(); err != nil {
			a.logger.Error(err, "migrate down failed")
			return
		}
	}

	a.logger.Info("migration successful", "dsn", a.dbOptions.URI)
}
