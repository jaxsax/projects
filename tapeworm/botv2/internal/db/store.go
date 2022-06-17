package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	fdbbleve "github.com/jaxsax/projects/tapeworm/botv2/internal/fdb-bleve"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	URI               string `long:"db_uri" description:"uri to connect to database" default:"./bot.db" env:"DB_URI"`
	EnableBleveSearch bool   `long:"enable_bleve_search" env:"enable_bleve_search"`
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

	var linkIndex bleve.Index
	if opt.EnableBleveSearch {
		pathName := "/tmp/links.bleve"
		indexingRule := bleve.NewIndexMapping()

		_, err = os.Stat(pathName)
		fdbConfig := map[string]interface{}{
			"fdbAPIVersion": 710,
			"clusterFile":   "fdb.cluster",
		}
		if os.IsNotExist(err) {
			index, err := bleve.NewUsing(pathName, indexingRule, "upside_down", fdbbleve.Name, fdbConfig)
			if err != nil {
				return nil, fmt.Errorf("new using: %w", err)
			}

			linkIndex = index
		} else {
			index, err := bleve.OpenUsing(pathName, fdbConfig)
			if err != nil {
				return nil, fmt.Errorf("open using: %w", err)
			}

			linkIndex = index
		}
	}

	return NewStore(db, linkIndex), nil
}
