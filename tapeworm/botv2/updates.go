package botv2

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

type Update struct {
	ID   int64
	Data *tgbotapi.Update
}

type UpdateDB struct {
	*sqlx.DB
}

func NewUpdateDB(db *sqlx.DB) *UpdateDB {
	return &UpdateDB{db}
}

func (u *UpdateDB) Create(update Update) error {
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

	_, err = u.NamedExec(query, tempUpdate)
	if err != nil {
		return err
	}

	return nil
}
