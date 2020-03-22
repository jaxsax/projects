package sql

import (
	"encoding/json"

	"github.com/jaxsax/projects/tapeworm/botv2/updates"
	"github.com/jmoiron/sqlx"
)

type updatesRepository struct {
	db *sqlx.DB
}

func NewUpdatesRepository(db *sqlx.DB) updates.Repository {
	return &updatesRepository{db}
}

func (repo *updatesRepository) Create(update updates.Update) error {
	query := `INSERT INTO updates(data)
				VALUES(:data)`

	body, err := json.Marshal(update.Data)
	if err != nil {
		return err
	}

	tempUpdate := struct {
		Data string
	}{
		string(body),
	}

	_, err = repo.db.NamedExec(query, tempUpdate)
	if err != nil {
		return err
	}

	return nil
}
