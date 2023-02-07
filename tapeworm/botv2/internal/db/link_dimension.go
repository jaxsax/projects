package db

import (
	"context"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type LinkDimension struct {
	ID     uint64 `db:"id"`
	LinkID uint64 `db:"link_id"`
	Kind   string `db:"kind"`
	Data   string `db:"data"`
}

func toDAOLinkDimension(dim *types.Dimension) (*LinkDimension, error) {
	return &LinkDimension{
		Kind: string(dim.Kind),
		Data: string(dim.Data),
	}, nil
}

func (q *Queries) CreateLinkDimension(ctx context.Context, linkID uint64, dimension *types.Dimension) error {
	ld, err := toDAOLinkDimension(dimension)
	if err != nil {
		return err
	}

	_, err = q.ExecContext(ctx, `
        INSERT INTO link_dimension (link_id, kind, data) VALUES(?, ?, ?)
		`, linkID, ld.Kind, ld.Data,
	)
	if err != nil {
		return err
	}

	return nil
}
