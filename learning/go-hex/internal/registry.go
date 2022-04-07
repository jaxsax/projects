package internal

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Registry struct {
	config *Config
	DB     *sqlx.DB
}

type Config struct {
	SQLPath string
}

func NewRegistry(cfg *Config) (*Registry, error) {
	db, err := sqlx.Open("sqlite3", cfg.SQLPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Registry{
		config: cfg,
		DB:     db,
	}, nil
}

func (r *Registry) Setup() error {
	driver, err := sqlite3.WithInstance(r.DB.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("get migrate instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"store.db",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	_, err = r.DB.Exec("INSERT INTO pets(name) values(?)", "Jonnay")
	if err != nil {
		log.Printf("setup data failure: %v", err)
	}

	return nil
}
