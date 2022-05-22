package db

import (
	"context"
	"encoding/json"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type Link struct {
	ID        uint64 `db:"id"`
	Link      string `db:"link"`
	Title     string `db:"title"`
	CreatedTS uint64 `db:"created_ts"`
	CreatedBy uint64 `db:"created_by"`
	ExtraData string `db:"extra_data"`
	DeletedAt uint64 `db:"deleted_at"`
}

func toDAOLink(link *types.Link) (*Link, error) {
	var l Link

	l.Link = link.Link
	l.Title = link.Title
	l.CreatedTS = uint64(link.CreatedAt.Unix())
	l.CreatedBy = link.CreatedByID

	extraDataBytes, err := json.Marshal(link.ExtraData)
	if err != nil {
		return nil, err
	}

	l.ExtraData = string(extraDataBytes)

	return &l, nil
}

func (q *Queries) CreateLink(ctx context.Context, link *types.Link) error {
	l, err := toDAOLink(link)
	if err != nil {
		return err
	}

	_, err = q.ExecContext(ctx, `
		INSERT INTO links (link, title, created_ts, created_by, extra_data) VALUES (
			?, ?, ?, ?, ?
		)
	`, l.Link, l.Title, l.CreatedTS, l.CreatedBy, l.ExtraData)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateLinks(ctx context.Context, links []*types.Link) error {
	err := s.execTx(ctx, func(q *Queries) error {
		for _, link := range links {
			if err := q.CreateLink(ctx, link); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
