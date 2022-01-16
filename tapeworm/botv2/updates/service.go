package updates

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/utils"
	"github.com/jaxsax/projects/tapeworm/botv2/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Create(ctx context.Context, update Update) error {
	tx, err := utils.GetTx(ctx)
	if err != nil {
		return err
	}

	var db models.Update

	body, err := json.Marshal(update.Data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	db.Data = string(body)

	err = db.Insert(ctx, tx, boil.Blacklist("id"))
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}
