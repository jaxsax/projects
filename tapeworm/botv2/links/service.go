package links

import (
	"context"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/utils"
	"github.com/jaxsax/projects/tapeworm/botv2/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetLinks(ctx context.Context) ([]*Link, error) {
	tx, err := utils.GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tx: %w", err)
	}

	objs, err := models.Links(
		models.LinkWhere.DeletedAt.EQ(0),
		qm.OrderBy("created_ts desc"),
	).All(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	links := make([]*Link, 0, len(objs))
	for _, dbLink := range objs {
		links = append(links, dbLinksToLink(dbLink))
	}

	return links, nil
}
