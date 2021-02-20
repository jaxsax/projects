package sql

import (
	"encoding/json"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/links"
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

type linksRepository struct {
	db *sqlx.DB
}

// NewLinksRepository creates a new links.Repository backed by sql
func NewLinksRepository(db *sqlx.DB) links.Repository {
	return &linksRepository{
		db: db,
	}
}

func toDBObject(l links.Link) link {
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

func (l *link) fromDBObject() links.Link {
	var extraData = map[string]interface{}{}
	_ = json.Unmarshal([]byte(l.ExtraData), &extraData)

	return links.Link{
		ID:        l.ID,
		CreatedTS: l.CreatedTS,
		CreatedBy: l.CreatedBy,
		Link:      l.Link,
		Title:     l.Title,
		ExtraData: extraData,
	}
}

func (repo *linksRepository) CreateMany(links []links.Link) error {
	query := `INSERT INTO links(link, title, created_ts, created_by, extra_data)
				VALUES(:link, :title, :created_ts, :created_by, :extra_data)`

	for _, link := range links {
		_, err := repo.db.NamedExec(query, toDBObject(link))
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *linksRepository) List() ([]links.Link, error) {
	dbLinks := []link{}
	err := repo.db.Select(&dbLinks, "SELECT * FROM links ORDER BY created_ts desc")
	if err != nil {
		return []links.Link{}, fmt.Errorf("select: %w", err)
	}

	links := make([]links.Link, len(dbLinks))
	for i, link := range dbLinks {
		links[i] = link.fromDBObject()
	}

	return links, nil
}
