package skippedlinks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Sqlite3 struct {
	db *sql.DB
}

var _ Repository = &Sqlite3{}

func NewSqliteRepository(db *sql.DB) *Sqlite3 {
	return &Sqlite3{
		db: db,
	}
}

func (sql *Sqlite3) Create(skippedLink SkippedLink) error {
	var dbLink models.SkippedLink
	dbLink.Link = skippedLink.Link
	dbLink.ErrorReason = skippedLink.ErrorReason

	err := dbLink.Insert(context.TODO(), sql.db, boil.Blacklist("id"))
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}
