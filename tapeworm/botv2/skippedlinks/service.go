package skippedlinks

import (
	"context"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/utils"
	"github.com/jaxsax/projects/tapeworm/botv2/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Create(ctx context.Context, skippedLink SkippedLink) error {
	tx, err := utils.GetTx(ctx)
	if err != nil {
		return err
	}

	var dbLink models.SkippedLink
	dbLink.Link = skippedLink.Link
	dbLink.ErrorReason = skippedLink.ErrorReason

	err = dbLink.Insert(ctx, tx, boil.Blacklist("id"))
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}
