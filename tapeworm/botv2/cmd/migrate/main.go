package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
)

var (
	flagParser     = flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)
	dbOptions      = &db.Options{}
	migrateOptions = &migrateOption{}
)

type migrateOption struct {
	Up   bool `long:"up" description:"migrate forward" env:"MIGRATE_UP"`
	Down bool `long:"down" description:"migrate backwards" env:"MIGRATE_DOWN"`
}

func main() {
	if _, err := flagParser.AddGroup("db", "", dbOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.AddGroup("migrate", "", migrateOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.Parse(); err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("sqlite3", dbOptions.URI)
	if err != nil {
		log.Fatalf("sql connect error=%v", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("sql ping error=%v", err)
		return
	}

	sqliteInstance, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatalf("sqlite.withinstance err=%v", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"sqlite3",
		sqliteInstance,
	)
	if err != nil {
		log.Fatalf("failed to create database instance err=%v", err)
		return
	}

	if migrateOptions.Up && migrateOptions.Down {
		log.Fatalf("What are you doing???")
		return
	}

	if migrateOptions.Up {
		if err := m.Up(); err != nil {
			log.Fatalf("migrate ip failed err=%v", err)
			return
		}
	}

	if migrateOptions.Down {
		if err := m.Down(); err != nil {
			log.Fatalf("migrate down failed err=%v", err)
			return
		}
	}

	log.Printf("Migration succcessful|urn=%s", dbOptions.URI)
}
