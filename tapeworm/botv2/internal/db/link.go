package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type Label struct {
	ID     uint64 `db:"id"`
	LinkID uint64 `db:"link_id"`
	Name   string `db:"label_name"`
}

type Link struct {
	ID           uint64 `db:"id"`
	Link         string `db:"link"`
	Title        string `db:"title"`
	CreatedTS    uint64 `db:"created_ts"`
	CreatedBy    uint64 `db:"created_by"`
	ExtraData    string `db:"extra_data"`
	Host         string `db:"host"`
	Path         string `db:"path"`
	DeletedAt    uint64 `db:"deleted_at"`
	DimCollected int    `db:"dim_collected"`
}

func toDAOLink(link *types.Link) (*Link, error) {
	var l Link

	l.ID = link.ID
	l.Link = link.Link
	l.Title = link.Title
	l.CreatedTS = uint64(link.CreatedAt)
	l.CreatedBy = link.CreatedByID
	l.Host = link.Domain
	l.Path = link.Path
	l.DimCollected = 0

	if link.DimCollected {
		l.DimCollected = 1
	}

	extraDataBytes, err := json.Marshal(link.ExtraData)
	if err != nil {
		return nil, err
	}

	l.ExtraData = string(extraDataBytes)

	return &l, nil
}

func (q *Queries) toTypesLink(link *Link) (*types.Link, error) {
	var l types.Link

	l.ID = link.ID
	l.Link = link.Link
	l.Title = link.Title
	l.CreatedAt = link.CreatedTS
	l.CreatedByID = link.CreatedBy
	l.Domain = link.Host
	l.Path = link.Path
	l.DimCollected = false

	if link.DeletedAt > 0 {
		deletedAt := time.Unix(int64(link.DeletedAt), 0)
		l.DeletedAt = &deletedAt
	}

	if link.DimCollected >= 1 {
		l.DimCollected = true
	}

	u, err := url.Parse(l.Link)
	if err != nil {
		return nil, err
	}

	l.Domain = u.Hostname()

	linkLabels, err := q.ListLinkLabels(context.TODO(), link.ID)
	if err != nil {
		return nil, err
	}

	l.Labels = linkLabels

	dimensions, err := q.ListLinkDimensions(context.TODO(), link.ID)
	if err != nil {
		return nil, err
	}

	l.Dimensions = dimensions

	return &l, nil
}

func (q *Queries) CreateLink(ctx context.Context, link *types.Link) error {
	l, err := toDAOLink(link)
	if err != nil {
		return err
	}

	_, err = q.ExecContext(ctx, `
		INSERT INTO links (link, title, created_ts, created_by, extra_data, host, path) VALUES (
			?, ?, ?, ?, ?, ?, ?
		)
	`, l.Link, l.Title, l.CreatedTS, l.CreatedBy, l.ExtraData, l.Host, l.Path)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queries) ListLinks(ctx context.Context) ([]*types.Link, error) {
	return q.ListLinksWithFilter(ctx, &types.LinkFilter{})
}

type andPair struct {
	fieldName string
	operator  string
	value     interface{}
}

func (q *Queries) buildLinkFilterStmt(selectStmt string, filter *types.LinkFilter) (string, []interface{}) {
	var andPairs []andPair

	andPairs = append(andPairs, andPair{
		fieldName: "deleted_at",
		operator:  "=",
		value:     0,
	})

	if filter.LinkWithoutScheme != "" {
		andPairs = append(andPairs, andPair{
			fieldName: "link",
			operator:  "LIKE",
			value:     fmt.Sprintf("%%://%s", filter.LinkWithoutScheme),
		})
	}

	if filter.Domain != "" {
		andPairs = append(andPairs, andPair{
			fieldName: "host",
			operator:  "=",
			value:     filter.Domain,
		})
	}

	if filter.TitleSearch != "" {
		andPairs = append(andPairs, andPair{
			fieldName: "lower(title)",
			operator:  "LIKE",
			value:     "%" + filter.TitleSearch + "%",
		})
	}

	if filter.DimCollected != nil {
		value := "0"
		if *filter.DimCollected {
			value = "1"
		}

		andPairs = append(andPairs, andPair{
			fieldName: "dim_collected",
			operator:  "=",
			value:     value,
		})
	}

	stmtParts := []string{
		selectStmt,
	}

	values := make([]interface{}, 0)
	if len(andPairs) > 0 {
		andStatements := make([]string, 0, len(andPairs))
		for _, p := range andPairs {
			andStatements = append(andStatements, fmt.Sprintf("%s %s ?", p.fieldName, p.operator))
			values = append(values, p.value)
		}

		stmtParts = append(stmtParts, "WHERE")
		stmtParts = append(stmtParts, strings.Join(andStatements, " AND "))
	}

	if filter.UniqueLink {
		stmtParts = append(stmtParts, "GROUP BY link")
	}

	stmtParts = append(stmtParts, "ORDER BY created_ts DESC")

	if filter.Limit > 0 {
		stmtParts = append(stmtParts, "LIMIT ?")
		values = append(values, strconv.Itoa(filter.Limit))
	}

	if filter.PageNumber > 1 && filter.Limit > 0 {
		stmtParts = append(stmtParts, "OFFSET ?")
		values = append(values, strconv.Itoa((filter.PageNumber-1)*filter.Limit))
	}

	stmt := strings.Join(stmtParts, " ")
	return stmt, values
}

func (q *Queries) toTypesLabel(l *Label) (*types.Label, error) {
	return &types.Label{
		Name: l.Name,
	}, nil
}

func (q *Queries) ListLinkLabels(ctx context.Context, linkID uint64) ([]*types.Label, error) {
	stmt := "SELECT * FROM link_label WHERE link_id = ?"
	rs, err := q.QueryxContext(ctx, stmt, linkID)
	if err != nil {
		return nil, err
	}

	labels := make([]*types.Label, 0)
	for rs.Next() {
		var label Label
		if err := rs.StructScan(&label); err != nil {
			return nil, err
		}

		lt, err := q.toTypesLabel(&label)
		if err != nil {
			return nil, err
		}

		labels = append(labels, lt)
	}

	return labels, nil
}

func (q *Queries) toTypesDimension(l *LinkDimension) (*types.Dimension, error) {
	return &types.Dimension{
		Kind: types.DimensionKind(l.Kind),
		Data: json.RawMessage(l.Data),
	}, nil
}

func (q *Queries) ListLinkDimensions(ctx context.Context, linkID uint64) ([]*types.Dimension, error) {
	stmt := "SELECT * FROM link_dimension WHERE link_id = ?"
	rs, err := q.QueryxContext(ctx, stmt, linkID)
	if err != nil {
		return nil, err
	}

	dimensions := make([]*types.Dimension, 0)
	for rs.Next() {
		var dimension LinkDimension
		if err := rs.StructScan(&dimension); err != nil {
			return nil, err
		}

		lt, err := q.toTypesDimension(&dimension)
		if err != nil {
			return nil, err
		}

		dimensions = append(dimensions, lt)
	}

	return dimensions, err
}

func (q *Queries) CountLinksWithFilter(ctx context.Context, filter *types.LinkFilter) (int, error) {
	stmt, values := q.buildLinkFilterStmt("SELECT COUNT(*) FROM links", filter)
	logging.FromContext(ctx).V(1).Info("query", "stmt", stmt, "values", values)

	var count int
	if err := q.GetContext(ctx, &count, stmt, values...); err != nil {
		return 0, err
	}

	return count, nil
}

func (q *Queries) ListLinksWithFilter(ctx context.Context, filter *types.LinkFilter) ([]*types.Link, error) {
	stmt, values := q.buildLinkFilterStmt("SELECT * FROM links", filter)
	logging.FromContext(ctx).V(1).Info("query", "stmt", stmt, "values", values)
	rs, err := q.QueryxContext(ctx, stmt, values...)
	if err != nil {
		return nil, err
	}

	links := make([]*types.Link, 0)
	for rs.Next() {
		var obj Link
		if err := rs.StructScan(&obj); err != nil {
			return nil, err
		}

		link, err := q.toTypesLink(&obj)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

func (q *Queries) GetLink(ctx context.Context, id uint64) (*types.Link, error) {
	link := new(Link)
	err := q.GetContext(ctx, link, "SELECT * FROM links WHERE id = ? and deleted_at = 0", id)
	if err != nil {
		return nil, err
	}

	typesLink, err := q.toTypesLink(link)
	if err != nil {
		return nil, err
	}

	return typesLink, nil
}

func (q *Queries) UpdateLink(ctx context.Context, link *types.Link) error {
	daoLink, err := toDAOLink(link)
	if err != nil {
		return err
	}

	_, err = q.ExecContext(ctx, `
			UPDATE links SET link = ?, title = ?, host = ?, path = ?
			WHERE id = ?
		`, daoLink.Link, daoLink.Title, daoLink.Host, daoLink.Path, daoLink.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queries) SetLinkDimCollected(ctx context.Context, linkID uint64, collected bool) error {
	collectedVal := 0
	if collected {
		collectedVal = 1
	}

	_, err := q.ExecContext(ctx, `
    UPDATE links SET dim_collected = ? WHERE id = ?
    `, collectedVal, linkID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateLinks(ctx context.Context, links []*types.Link) error {
	err := s.execTx(ctx, func(q *Queries) error {
		for _, link := range links {
			if err := q.UpdateLink(ctx, link); err != nil {
				return err
			}
		}

		return nil
	})

	return err
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

func (s *Store) UpdateLinkDimensions(ctx context.Context, link *types.Link, dimensions []*types.Dimension) error {
	return s.execTx(ctx, func(q *Queries) error {
		if err := q.SetLinkDimCollected(ctx, link.ID, true); err != nil {
			return fmt.Errorf("update link: %w", err)
		}

		for _, dim := range dimensions {
			if err := q.CreateLinkDimension(ctx, link.ID, dim); err != nil {
				return fmt.Errorf("create link dimension: %w", err)
			}
		}

		return nil
	})
}
