package updates

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (sql *Sqlite3) Create(update Update) error {
	var db models.Update

	body, err := json.Marshal(update.Data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	db.Data = string(body)

	err = db.Insert(context.TODO(), sql.db, boil.Blacklist("id"))
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}
