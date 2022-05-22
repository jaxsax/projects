package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	URI string `long:"db_uri" description:"uri to connect to database" default:"./bot.db" env:"DB_URI"`
}

type Store struct {
	*Queries
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db:      db,
		Queries: &Queries{db},
	}
}

type Queries struct {
	sqlx.ExecerContext
}

func (q *Queries) WithTx(tx *sqlx.Tx) *Queries {
	return &Queries{
		tx,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := &Queries{tx}
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func Setup(opt *Options) (*Store, error) {
	db, err := sqlx.Connect("sqlite3", opt.URI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return NewStore(db), nil
}