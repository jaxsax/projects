package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	URI string `long:"db_uri" description:"uri to connect to database" default:"./bot.db" env:"DB_URI"`
}

type Store struct {
	*Queries
	db *sqlx.DB

	linkIndex bleve.Index
}

func NewStore(
	db *sqlx.DB,
	titleIndex bleve.Index,
) *Store {
	return &Store{
		db: db,
		Queries: &Queries{
			db, db, db,
		},
		linkIndex: titleIndex,
	}
}

type SingleRow interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Queries struct {
	sqlx.ExecerContext
	sqlx.QueryerContext
	SingleRow
}

func (q *Queries) WithTx(tx *sqlx.Tx) *Queries {
	return &Queries{
		tx, tx, tx,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := &Queries{tx, tx, tx}
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

	pathName := "/tmp/links.bleve"
	indexingRule := bleve.NewIndexMapping()
	var linkIndex bleve.Index

	_, err = os.Stat(pathName)
	if os.IsNotExist(err) {
		index, err := bleve.New(pathName, indexingRule)
		if err != nil {
			return nil, err
		}

		linkIndex = index
	} else {
		index, err := bleve.Open(pathName)
		if err != nil {
			return nil, err
		}

		linkIndex = index
	}

	return NewStore(db, linkIndex), nil
}
