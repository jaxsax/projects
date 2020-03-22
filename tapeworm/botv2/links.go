package botv2

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type link struct {
	ID        int64 `db:"id"`
	CreatedTS int64 `db:"created_ts"`
	CreatedBy int64 `db:"created_by"`
	Link      string
	Title     string
	ExtraData string `db:"extra_data"`
}

type Link struct {
	ID        int64
	CreatedTS int64
	CreatedBy int64
	Link      string
	Title     string
	ExtraData map[string]interface{}
}

type LinksDB struct {
	db *sqlx.DB
}

func NewLinksDB(db *sqlx.DB) *LinksDB {
	return &LinksDB{
		db: db,
	}
}

func (l *Link) toDBObject() link {
	var extraDataAsString = "{}"
	body, err := json.Marshal(l.ExtraData)
	if err == nil {
		extraDataAsString = string(body)
	}

	return link{
		ID:        l.ID,
		CreatedTS: l.CreatedTS,
		CreatedBy: l.CreatedBy,
		Link:      l.Link,
		Title:     l.Title,
		ExtraData: extraDataAsString,
	}
}

func (l *link) fromDBObject() Link {
	var extraData = map[string]interface{}{}
	_ = json.Unmarshal([]byte(l.ExtraData), &extraData)

	return Link{
		ID:        l.ID,
		CreatedTS: l.CreatedTS,
		CreatedBy: l.CreatedBy,
		Link:      l.Link,
		Title:     l.Title,
		ExtraData: extraData,
	}
}

func (l *LinksDB) List() ([]Link, error) {
	dbLinks := []link{}
	err := l.db.Select(&dbLinks, "SELECT * FROM links")
	if err != nil {
		return nil, err
	}

	links := make([]Link, len(dbLinks))
	for i, link := range dbLinks {
		links[i] = link.fromDBObject()
	}

	return links, nil
}

func (l *LinksDB) Create(links []Link) error {
	query := `INSERT INTO links(link, title, created_ts, created_by, extra_data)
				VALUES(:link, :title, :created_ts, :created_by, :extra_data)`

	for _, link := range links {
		_, err := l.db.NamedExec(query, link.toDBObject())
		if err != nil {
			return err
		}
	}
	return nil
}
