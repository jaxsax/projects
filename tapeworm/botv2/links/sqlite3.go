package links

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/jaxsax/projects/tapeworm/botv2/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Sqlite3 struct {
	db *sql.DB
}

var _ Repository = &Sqlite3{}

func NewSqliteRepository(db *sql.DB) *Sqlite3 {
	return &Sqlite3{
		db: db,
	}
}

func dbLinksToLink(m *models.Link) *Link {
	var data map[string]string

	_ = json.Unmarshal([]byte(m.ExtraData), &data)

	var domain *string
	if u, err := url.Parse(m.Link); err == nil {
		var host = u.Hostname()
		domain = &host
	}

	return &Link{
		ID:        m.ID,
		CreatedTS: m.CreatedTS,
		CreatedBy: m.CreatedBy,
		DeletedAt: m.DeletedAt,
		Link:      m.Link,
		Title:     m.Title,
		Domain:    domain,
		ExtraData: data,
	}
}

var successiveSpaces = regexp.MustCompile(`\s+`)

func linkToDBLink(link *Link) *models.Link {
	var extraData = "{}"
	marshalled, err := json.Marshal(link.ExtraData)
	if err == nil {
		extraData = string(marshalled)
	}

	// Cleanup titles
	var title = link.Title
	title = strings.TrimSpace(title)
	title = successiveSpaces.ReplaceAllLiteralString(title, " ")

	return &models.Link{
		ID:        link.ID,
		CreatedTS: link.CreatedTS,
		CreatedBy: link.CreatedBy,
		DeletedAt: link.DeletedAt,
		Link:      link.Link,
		Title:     title,
		ExtraData: extraData,
	}
}

func (sq *Sqlite3) CreateMany(links []Link) error {
	dbLinks := make([]*models.Link, 0, len(links))
	for _, link := range links {
		dbLinks = append(dbLinks, linkToDBLink(&link))
	}

	tx, err := sq.db.Begin()
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}

		_ = tx.Commit()
	}()

	for _, link := range dbLinks {
		err = link.Insert(context.TODO(), tx, boil.Blacklist("id"))
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	return nil
}

func (sq *Sqlite3) List() ([]Link, error) {
	objs, err := models.Links(
		models.LinkWhere.DeletedAt.EQ(0),
		qm.OrderBy("created_ts desc"),
	).All(context.TODO(), sq.db)
	if err != nil {
		return []Link{}, fmt.Errorf("all: %w", err)
	}

	links := make([]Link, 0, len(objs))
	for _, dbLink := range objs {
		links = append(links, *dbLinksToLink(dbLink))
	}

	return links, nil
}

func (sq *Sqlite3) ListMatchingIDs(IDs []int64) ([]Link, error) {

	objs, err := models.Links(
		models.LinkWhere.ID.IN(IDs),
		qm.OrderBy("created_ts desc"),
	).All(context.TODO(), sq.db)
	if err != nil {
		return []Link{}, fmt.Errorf("all: %w", err)
	}

	links := make([]Link, 0, len(objs))
	for _, dbLink := range objs {
		links = append(links, *dbLinksToLink(dbLink))
	}

	return links, nil
}
