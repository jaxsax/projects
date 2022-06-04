package db

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

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
	l.CreatedTS = uint64(link.CreatedAt)
	l.CreatedBy = link.CreatedByID

	extraDataBytes, err := json.Marshal(link.ExtraData)
	if err != nil {
		return nil, err
	}

	l.ExtraData = string(extraDataBytes)

	return &l, nil
}

func toTypesLink(link *Link) (*types.Link, error) {
	var l types.Link

	l.ID = link.ID
	l.Link = link.Link
	l.Title = link.Title
	l.CreatedAt = link.CreatedTS
	l.CreatedByID = link.CreatedBy

	if link.DeletedAt > 0 {
		deletedAt := time.Unix(int64(link.DeletedAt), 0)
		l.DeletedAt = &deletedAt
	}

	u, err := url.Parse(l.Link)
	if err != nil {
		return nil, err
	}

	l.Domain = u.Hostname()

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

func (q *Queries) ListLinks(ctx context.Context) ([]*types.Link, error) {
	rs, err := q.QueryxContext(ctx, `
		SELECT * FROM links WHERE deleted_at = 0
	`)
	if err != nil {
		return nil, err
	}

	links := make([]*types.Link, 0)
	for rs.Next() {
		var obj Link
		if err := rs.StructScan(&obj); err != nil {
			return nil, err
		}

		link, err := toTypesLink(&obj)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
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
