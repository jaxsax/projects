package db

import (
	"context"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type ContentBlocklist struct {
	ID       uint64 `db:"id"`
	Strategy string `db:"strategy"`
	Content  string `db:"content"`
}

func (q *Queries) toTypesContentBlocklistStrategy(db *ContentBlocklist) (*types.ContentBLocklistStrategy, error) {
	return &types.ContentBLocklistStrategy{
		ID:       db.ID,
		Strategy: db.Strategy,
		Content:  db.Content,
	}, nil
}

func (q *Queries) ListBLocklistStrategies(ctx context.Context) ([]*types.ContentBLocklistStrategy, error) {
	stmt := "SELECT * FROM content_blocklist"
	rs, err := q.QueryxContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	strategies := make([]*types.ContentBLocklistStrategy, 0)
	for rs.Next() {
		var row ContentBlocklist
		if err := rs.StructScan(&row); err != nil {
			return nil, err
		}

		lt, err := q.toTypesContentBlocklistStrategy(&row)
		if err != nil {
			return nil, err
		}

		strategies = append(strategies, lt)
	}

	return strategies, nil
}
