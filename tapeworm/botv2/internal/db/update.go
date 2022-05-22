package db

import (
	"context"
	"encoding/json"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type TelegramUpdate struct {
	ID   uint64 `db:"id"`
	Data string `db:"data"`
}

func (s *Queries) CreateTelegramUpdate(ctx context.Context, update *types.TelegramUpdate) error {
	updateBytes, err := json.Marshal(update.Data)
	if err != nil {
		return err
	}

	_, err = s.ExecContext(
		ctx,
		"INSERT INTO updates (data) VALUES(?)",
		string(updateBytes),
	)
	if err != nil {
		return err
	}

	return nil
}
